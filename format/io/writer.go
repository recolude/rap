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
	"github.com/recolude/rap/format/metadata"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

// https://dave.cheney.net/2019/01/27/eliminate-error-handling-by-eliminating-errors
type errWriter struct {
	io.Writer
	err error
	n   int
}

func (e *errWriter) Write(buf []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	var n int
	n, e.err = e.Writer.Write(buf)
	e.n += n
	return n, e.err
}

func (e *errWriter) TotalWritten() int {
	return e.n
}

func (e *errWriter) Close() error {
	return e.err
}

type encoderCollectionMapping struct {
	encoder         encoding.Encoder
	collections     []format.CaptureCollection
	collectionOrder []int
}

type Writer struct {
	encoders             []encoding.Encoder
	timeStorageTechnique TimeStorageTechnique
	compress             bool
	out                  io.Writer
}

// NewRecoludeWriter builds a new recording writer with default recolude
// encoders.
func NewRecoludeWriter(out io.Writer) Writer {
	return Writer{
		encoders: []encoding.Encoder{
			event.NewEncoder(event.Raw32),
			position.NewEncoder(position.Oct48),
			euler.NewEncoder(euler.Raw32),
			enum.NewEncoder(),
		},
		compress:             true,
		timeStorageTechnique: BST16,
		out:                  out,
	}
}

// NewWriter builds a new writer using the encoders provided.
func NewWriter(encoders []encoding.Encoder, compress bool, out io.Writer, timeStorageTechnique TimeStorageTechnique) Writer {
	return Writer{
		encoders:             encoders,
		out:                  out,
		compress:             compress,
		timeStorageTechnique: timeStorageTechnique,
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
		for key := range ref.Metadata().Mapping() {
			if _, ok := keyMappingToIndex[key]; !ok {
				keyMappingToIndex[key] = keyCount
				keyCount++
			}
		}
	}

	for _, bin := range recording.Binaries() {
		for key := range bin.Metadata().Mapping() {
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
	for i, collection := range recording.CaptureCollections() {
		if collection == nil {
			return nil, 0, errors.New("can not serialize recording with nil capture collections")
		}
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

func writeMetadata(out io.Writer, keyMappingToIndex map[string]int, block metadata.Block) (int, error) {
	metadataIndices := make([]uint, len(block.Mapping()))

	metadataValuesBuffer := bytes.Buffer{}
	i := 0
	for key, val := range block.Mapping() {
		metadataIndices[i] = uint(keyMappingToIndex[key])

		metadata.WriteProprty(&metadataValuesBuffer, val)
		i++
	}

	totalBytes := 0

	written, err := out.Write(rapbinary.UvarintArrayToBytes(metadataIndices))
	totalBytes += written
	if err != nil {
		return totalBytes, err
	}

	written, err = out.Write(metadataValuesBuffer.Bytes())
	totalBytes += written
	if err != nil {
		return totalBytes, err
	}

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

	totalWritten, err := out.Write(rapbinary.StringArrayToBytes(encoderSignatures))
	if err != nil {
		return totalWritten, err
	}

	varByte := make([]byte, binary.MaxVarintLen64)
	for _, version := range encoderVersions {
		read := binary.PutUvarint(varByte, uint64(version))
		written, err := out.Write(varByte[:read])
		totalWritten += written
		if err != nil {
			return totalWritten, err
		}
	}

	return totalWritten, err
}

func recurseRecordingToBytes(out io.Writer, recording format.Recording, keyMappingToIndex map[string]int, encodingBlocks [][]byte, streamIndexToEncoderUsedIndex []int, offset int, tech TimeStorageTechnique) (int, int, error) {
	ew := &errWriter{Writer: out}

	// Write id
	ew.Write(rapbinary.StringToBytes(recording.ID()))

	// Write name
	ew.Write(rapbinary.StringToBytes(recording.Name()))

	// Write metadata
	writeMetadata(ew, keyMappingToIndex, recording.Metadata())

	// Write number of streams
	numStreams := make([]byte, binary.MaxVarintLen64)
	read := binary.PutUvarint(numStreams, uint64(len(recording.CaptureCollections())))
	ew.Write(numStreams[:read])

	// Write all streams
	for streamIndex := range recording.CaptureCollections() {
		// Write index of the encoder used to encode stream
		numStreams := make([]byte, binary.MaxVarintLen64)
		read := binary.PutUvarint(numStreams, uint64(streamIndexToEncoderUsedIndex[offset+streamIndex]))
		ew.Write(numStreams[:read])

		encodeTime(tech, ew, recording.CaptureCollections()[streamIndex].Captures())

		// Write stream data
		ew.Write(rapbinary.BytesArrayToBytes(encodingBlocks[offset+streamIndex]))
	}

	// Write number of references
	numReferences := make([]byte, binary.MaxVarintLen64)
	read = binary.PutUvarint(numReferences, uint64(len(recording.BinaryReferences())))
	ew.Write(numReferences[:read])

	// Write binary references
	for _, ref := range recording.BinaryReferences() {
		ew.Write(rapbinary.StringToBytes(ref.Name()))
		ew.Write(rapbinary.StringToBytes(ref.URI()))

		refSize := make([]byte, binary.MaxVarintLen64)
		read = binary.PutUvarint(refSize, ref.Size())
		ew.Write(refSize[:read])

		writeMetadata(ew, keyMappingToIndex, ref.Metadata())
	}

	// Write number of binaries
	numBinaries := make([]byte, binary.MaxVarintLen64)
	read = binary.PutUvarint(numBinaries, uint64(len(recording.Binaries())))
	ew.Write(numBinaries[:read])

	// Write binaries
	for _, bin := range recording.Binaries() {
		ew.Write(rapbinary.StringToBytes(bin.Name()))

		refSize := make([]byte, binary.MaxVarintLen64)
		read = binary.PutUvarint(refSize, bin.Size())
		ew.Write(refSize[:read])

		writeMetadata(ew, keyMappingToIndex, bin.Metadata())

		actualbinaryWritten, _ := io.Copy(ew, bin.Data())
		if actualbinaryWritten != int64(bin.Size()) {
			panic("Binary data written was larger than size in signature")
		}
	}

	// Write number of recordings
	numRecordings := make([]byte, binary.MaxVarintLen64)
	read = binary.PutUvarint(numRecordings, uint64(len(recording.Recordings())))
	ew.Write(numRecordings[:read])

	// Write all child recordings
	newOffset := offset + len(recording.CaptureCollections())
	for _, rec := range recording.Recordings() {
		_, updatedOffset, err := recurseRecordingToBytes(ew, rec, keyMappingToIndex, encodingBlocks, streamIndexToEncoderUsedIndex, newOffset, tech)
		if err != nil {
			return ew.TotalWritten(), -1, err
		}
		newOffset = updatedOffset
	}

	return ew.TotalWritten(), newOffset, ew.err
}

func checkForNilInterfaces(recording format.Recording) error {
	if recording == nil {
		return errors.New("can not serialize recording with nil sub-recordings")
	}

	for _, bin := range recording.BinaryReferences() {
		if bin == nil {
			return errors.New("can not serialize recording with nil binary references")
		}
	}

	for _, bin := range recording.Binaries() {
		if bin == nil {
			return errors.New("can not serialize recording with nil binaries")
		}
	}

	for _, rec := range recording.Recordings() {
		err := checkForNilInterfaces(rec)
		if err != nil {
			return err
		}
	}

	return nil
}

// Write will take the recording provided and write it to the underlying stream
// the writer was built with.
func (w Writer) Write(recording format.Recording) (int, error) {
	if recording == nil {
		panic(errors.New("can not write nil recording"))
	}

	err := checkForNilInterfaces(recording)
	if err != nil {
		return 0, err
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

	// Build compression writer
	var compressWriter io.WriteCloser
	if w.compress {
		written, err = w.out.Write([]byte{1})
		totalBytesWritten += written
		if err != nil {
			return totalBytesWritten, err
		}
		compressWriter, err = flate.NewWriter(w.out, 9 /*Best Compression*/)
		if err != nil {
			return totalBytesWritten, err
		}
	} else {
		written, err = w.out.Write([]byte{0})
		totalBytesWritten += written
		if err != nil {
			return totalBytesWritten, err
		}
		compressWriter = &errWriter{Writer: w.out}
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
		written, err = compressWriter.Write(rapbinary.BytesArrayToBytes(header))
		totalBytesWritten += written
		if err != nil {
			return totalBytesWritten, err
		}
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
	written, _, err = recurseRecordingToBytes(compressWriter, recording, keyMappingToIndex, encodingBlocks, streamIndexToEncoderUsedIndex, 0, w.timeStorageTechnique)
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
