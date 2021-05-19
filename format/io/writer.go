package io

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	"github.com/recolude/rap/format/encoding/enum"
	"github.com/recolude/rap/format/encoding/euler"
	"github.com/recolude/rap/format/encoding/event"
	"github.com/recolude/rap/format/encoding/position"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type encoderCollectionMapping struct {
	encoder         encoding.Encoder
	collections     []format.CaptureCollection
	collectionOrder []int
}

type Writer struct {
	encoders []encoding.Encoder
	out      io.Writer
}

// NewRecoludeWriter builds a new recording writer with default recolude
// encoders.
func NewRecoludeWriter(out io.Writer) Writer {
	return Writer{
		encoders: []encoding.Encoder{
			event.NewEncoder(event.Raw32),
			position.NewEncoder(position.Oct48),
			euler.NewEncoder(euler.Raw32),
			enum.NewEncoder(enum.Raw32),
		},
		out: out,
	}
}

// NewWriter builds a new writer using the encoders provided.
func NewWriter(encoders []encoding.Encoder, out io.Writer) Writer {
	return Writer{
		encoders: encoders,
		out:      out,
	}
}

func calcNumStreams(recording format.Recording) int {
	total := 0
	for _, rec := range recording.Recordings() {
		total += calcNumStreams(rec)
	}
	return len(recording.CaptureCollections()) + total
}

// accumulateMetdataKeys builds out a mapping of metadata keys to some unique
// index. Used to ensure the key is only ever written once to file.
func accumulateMetdataKeys(recording format.Recording, keyMappingToIndex map[string]int) {
	keyCount := len(keyMappingToIndex)

	for key, _ := range recording.Metadata().Mapping() {
		if _, ok := keyMappingToIndex[key]; !ok {
			keyMappingToIndex[key] = keyCount
			keyCount++
		}
	}

	for _, ref := range recording.BinaryReferences() {
		for key, _ := range ref.Metadata().Mapping() {
			if _, ok := keyMappingToIndex[key]; !ok {
				keyMappingToIndex[key] = keyCount
				keyCount++
			}
		}
	}

	for _, rec := range recording.Recordings() {
		accumulateMetdataKeys(rec, keyMappingToIndex)
	}
}

func (w Writer) evaluateCollections(recording format.Recording, offset int) ([]encoderCollectionMapping, int, error) {
	mappings := make([]encoderCollectionMapping, 0)
	streamsSatisfied := make([]bool, len(recording.CaptureCollections()))
	for i := range recording.CaptureCollections() {
		streamsSatisfied[i] = false
	}

	for i, encoder := range w.encoders {

		mapping := encoderCollectionMapping{encoder: w.encoders[i]}
		for streamIndex, stream := range recording.CaptureCollections() {
			if streamsSatisfied[streamIndex] == false && encoder.Accepts(stream) {
				mapping.collections = append(mapping.collections, stream)
				mapping.collectionOrder = append(mapping.collectionOrder, streamIndex+offset)
				streamsSatisfied[streamIndex] = true
			}
		}

		if len(mapping.collections) > 0 {
			mappings = append(mappings, mapping)
		}
	}

	for i, stream := range recording.CaptureCollections() {
		if streamsSatisfied[i] == false {
			return nil, 0, fmt.Errorf("no encoder registered to handle stream: %s", stream.Signature())
		}
	}

	curOffset := offset + len(recording.CaptureCollections())

	for _, childRecording := range recording.Recordings() {
		childMappings, newOffset, err := w.evaluateCollections(childRecording, curOffset)
		curOffset = newOffset
		if err != nil {
			return nil, 0, err
		}
		for _, childMap := range childMappings {
			found := false
			for i, ourMap := range mappings {
				if ourMap.encoder.Signature() == childMap.encoder.Signature() {
					mappings[i].collections = append(ourMap.collections, childMap.collections...)
					mappings[i].collectionOrder = append(ourMap.collectionOrder, childMap.collectionOrder...)
					found = true
				}
			}
			if found == false {
				mappings = append(mappings, childMap)
			}
		}
	}

	return mappings, curOffset, nil
}

func writeMetadata(out io.Writer, keyMappingToIndex map[string]int, metadata format.Metadata) (int, error) {
	metadataIndices := make([]uint, len(metadata.Mapping()))

	metadataValuesBuffer := bytes.Buffer{}
	// metadataValues := make([][]byte, len(metadata.Mapping()))
	i := 0
	for key, val := range metadata.Mapping() {
		metadataIndices[i] = uint(keyMappingToIndex[key])
		metadataValuesBuffer.WriteByte(val.Code())
		metadataValuesBuffer.Write(val.Data())
		// metadataValues[i] = append([]byte{val.Code()}, val.Data()...)
		i++
	}

	totalBytes := 0

	written, err := out.Write(rapbinary.UintArrayToBytes(metadataIndices))
	totalBytes += written
	if err != nil {
		return written, err
	}

	written, err = out.Write(rapbinary.BytesArrayToBytes(metadataValuesBuffer.Bytes()))
	totalBytes += written

	return totalBytes, err
}

func writeEncoders(out io.Writer, encoders []encoderCollectionMapping) (int, error) {
	encoderSignatures := make([]string, len(encoders))
	encoderVersions := make([]uint, len(encoders))
	i := 0
	for _, encoderMapping := range encoders {
		encoderSignatures[i] = encoderMapping.encoder.Signature()
		encoderVersions[i] = encoderMapping.encoder.Version()
		i++
	}

	written, err := out.Write(rapbinary.StringArrayToBytes(encoderSignatures))
	if err != nil {
		return written, err
	}

	writtenVersions, err := out.Write(rapbinary.UintArrayToBytes(encoderVersions))

	return writtenVersions + written, err
}

func recurseRecordingToBytes(recording format.Recording, keyMappingToIndex map[string]int, encodingBlocks [][]byte, streamIndexToEncoderUsedIndex []int, offset int) (int, []byte) {
	out := bytes.Buffer{}

	// Write id
	out.Write(rapbinary.StringToBytes(recording.ID()))

	// Write name
	out.Write(rapbinary.StringToBytes(recording.Name()))

	// Write metadata
	writeMetadata(&out, keyMappingToIndex, recording.Metadata())

	// Write number of streams
	numStreams := make([]byte, 4)
	read := binary.PutUvarint(numStreams, uint64(len(recording.CaptureCollections())))
	out.Write(numStreams[:read])

	// Write all streams
	for streamIndex := range recording.CaptureCollections() {
		// Write index of the encoder used to encode stream
		numStreams := make([]byte, 4)
		read := binary.PutUvarint(numStreams, uint64(streamIndexToEncoderUsedIndex[offset+streamIndex]))
		out.Write(numStreams[:read])

		// Write stream data
		out.Write(rapbinary.BytesArrayToBytes(encodingBlocks[offset+streamIndex]))
	}

	// Write number of references
	numReferences := make([]byte, 4)
	read = binary.PutUvarint(numReferences, uint64(len(recording.BinaryReferences())))
	out.Write(numReferences[:read])

	// Write binary references
	for _, ref := range recording.BinaryReferences() {
		out.Write(rapbinary.StringToBytes(ref.Name()))
		out.Write(rapbinary.StringToBytes(ref.URI()))

		refSize := make([]byte, 4)
		read = binary.PutUvarint(refSize, ref.Size())
		out.Write(refSize[:read])

		writeMetadata(&out, keyMappingToIndex, ref.Metadata())
	}
	// out.Write(rapbinary.StringArrayToBytes(binaryRefNames))

	// Write number of recordings
	numRecordings := make([]byte, 4)
	read = binary.PutUvarint(numRecordings, uint64(len(recording.Recordings())))
	out.Write(numRecordings[:read])

	// Write all child recordings
	newOffset := offset + len(recording.CaptureCollections())
	for _, rec := range recording.Recordings() {
		updatedOffset, recordingData := recurseRecordingToBytes(rec, keyMappingToIndex, encodingBlocks, streamIndexToEncoderUsedIndex, newOffset)
		newOffset = updatedOffset
		out.Write(rapbinary.BytesArrayToBytes(recordingData))
	}

	return newOffset, out.Bytes()
}

// Write will take the recording provided and write it to the underlying stream
// the writer was built with.
func (w Writer) Write(recording format.Recording) (int, error) {
	if recording == nil {
		panic(errors.New("can not write nil recording"))
	}

	encoderMappings, _, err := w.evaluateCollections(recording, 0)
	if err != nil {
		return 0, err
	}

	totalBytesWritten := 0

	// Write version number
	written, err := w.out.Write([]byte{2})
	totalBytesWritten += written
	if err != nil {
		return totalBytesWritten, err
	}

	// Write encoders used
	written, err = writeEncoders(w.out, encoderMappings)
	totalBytesWritten += written
	if err != nil {
		return totalBytesWritten, err
	}

	numStreams := calcNumStreams(recording)
	encodingBlocks := make([][]byte, numStreams)
	streamIndexToEncoderUsedIndex := make([]int, numStreams)

	for encoderIndex, val := range encoderMappings {
		header, streamsEncoded, err := val.encoder.Encode(val.collections)
		if err != nil {
			return totalBytesWritten, err
		}

		for i, order := range val.collectionOrder {
			encodingBlocks[order] = streamsEncoded[i]
			streamIndexToEncoderUsedIndex[order] = encoderIndex
		}

		// Write header
		written, err = w.out.Write(rapbinary.BytesArrayToBytes(header))
		totalBytesWritten += written
		if err != nil {
			return totalBytesWritten, err
		}
	}

	// Build compression buffer
	// compressBuffer := bytes.Buffer{}
	compressWriter, err := flate.NewWriter(w.out, 9 /*Best Compression*/)
	if err != nil {
		return totalBytesWritten, err
	}

	// Write metadata keys
	keyMappingToIndex := make(map[string]int)
	accumulateMetdataKeys(recording, keyMappingToIndex)
	allKeys := make([]string, len(keyMappingToIndex))
	for key, index := range keyMappingToIndex {
		allKeys[index] = key
	}
	written, err = compressWriter.Write(rapbinary.StringArrayToBytes(allKeys))
	totalBytesWritten += written
	if err != nil {
		return totalBytesWritten, err
	}

	// Write out all recordings
	_, allRecData := recurseRecordingToBytes(recording, keyMappingToIndex, encodingBlocks, streamIndexToEncoderUsedIndex, 0)
	written, err = compressWriter.Write(rapbinary.BytesArrayToBytes(allRecData))
	totalBytesWritten += written
	if err != nil {
		return totalBytesWritten, err
	}

	err = compressWriter.Close()
	if err != nil {
		return totalBytesWritten, err
	}

	return totalBytesWritten, err
}
