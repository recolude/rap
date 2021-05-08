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
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type encoderStreamMapping struct {
	encoder     encoding.Encoder
	streams     []format.CaptureStream
	streamOrder []int
}

type Writer struct {
	encoders []encoding.Encoder
	out      io.Writer
}

func NewWriter(encoders []encoding.Encoder, out io.Writer) Writer {
	return Writer{
		encoders: encoders,
		out:      out,
	}
}

func allRecordingsWithin(recording format.Recording) []format.Recording {
	recs := make([]format.Recording, 1)
	recs[0] = recording
	for _, rec := range recording.Recordings() {
		recs = append(recs, allRecordingsWithin(rec)...)
	}
	return recs
}

func calcNumStreams(recording format.Recording) int {
	total := 0
	for _, rec := range recording.Recordings() {
		total += calcNumStreams(rec)
	}
	return len(recording.CaptureStreams()) + total
}

// accumulateMetdataKeys builds out a mapping of metadata keys to some unique
// index. Used to ensure the key is only ever written once to file.
func accumulateMetdataKeys(recording format.Recording, keyMappingToIndex map[string]int) {
	keyCount := len(keyMappingToIndex)

	for key, _ := range recording.Metadata() {
		if _, ok := keyMappingToIndex[key]; !ok {
			keyMappingToIndex[key] = keyCount
			keyCount++
		}
	}

	for _, rec := range recording.Recordings() {
		accumulateMetdataKeys(rec, keyMappingToIndex)
	}
}

func (w Writer) evaluateStreams(recording format.Recording, offset int) ([]encoderStreamMapping, int, error) {
	mappings := make([]encoderStreamMapping, 0)
	streamsSatisfied := make([]bool, len(recording.CaptureStreams()))
	for i := range recording.CaptureStreams() {
		streamsSatisfied[i] = false
	}

	for i, encoder := range w.encoders {

		mapping := encoderStreamMapping{encoder: w.encoders[i]}
		for streamIndex, stream := range recording.CaptureStreams() {
			if streamsSatisfied[streamIndex] == false && encoder.Accepts(stream) {
				mapping.streams = append(mapping.streams, stream)
				mapping.streamOrder = append(mapping.streamOrder, streamIndex+offset)
				streamsSatisfied[streamIndex] = true
			}
		}

		if len(mapping.streams) > 0 {
			mappings = append(mappings, mapping)
		}
	}

	for i, stream := range recording.CaptureStreams() {
		if streamsSatisfied[i] == false {
			return nil, 0, fmt.Errorf("no encoder registered to handle stream: %s", stream.Signature())
		}
	}

	curOffset := offset + len(recording.CaptureStreams())

	for _, childRecording := range recording.Recordings() {
		childMappings, newOffset, err := w.evaluateStreams(childRecording, curOffset)
		curOffset = newOffset
		if err != nil {
			return nil, 0, err
		}
		for _, childMap := range childMappings {
			for _, ourMap := range mappings {
				if ourMap.encoder.Signature() == childMap.encoder.Signature() {
					ourMap.streams = append(ourMap.streams, childMap.streams...)
				}
			}
		}
	}

	return mappings, curOffset, nil
}

func writeMetadata(out io.Writer, keyMappingToIndex map[string]int, metadata map[string]string) (int, error) {
	metadataIndices := make([]uint, len(metadata))
	metadataValues := make([]string, len(metadata))
	i := 0
	for key, val := range metadata {
		metadataIndices[i] = uint(keyMappingToIndex[key])
		metadataValues[i] = val
		i++
	}

	totalBytes := 0

	written, err := out.Write(rapbinary.UintArrayToBytes(metadataIndices))
	totalBytes += written
	if err != nil {
		return written, err
	}

	written, err = out.Write(rapbinary.StringArrayToBytes(metadataValues))
	totalBytes += written

	return totalBytes, err
}

func writeEncoders(out io.Writer, encoders []encoderStreamMapping) (int, error) {
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

	// Write name
	out.Write(rapbinary.StringToBytes(recording.Name()))

	// Write metadata
	writeMetadata(&out, keyMappingToIndex, recording.Metadata())

	// Write number of streams
	numStreams := make([]byte, 4)
	read := binary.PutUvarint(numStreams, uint64(len(recording.CaptureStreams())))
	out.Write(numStreams[:read])

	// Write all streams
	for streamIndex := range recording.CaptureStreams() {
		// Write index of the encoder used to encode stream
		numStreams := make([]byte, 4)
		read := binary.PutUvarint(numStreams, uint64(streamIndexToEncoderUsedIndex[streamIndex]))
		out.Write(numStreams[:read])

		// Write stream data
		out.Write(rapbinary.BytesArrayToBytes(encodingBlocks[offset+streamIndex]))
	}

	// Write number of recordings
	numRecordings := make([]byte, 4)
	read = binary.PutUvarint(numRecordings, uint64(len(recording.Recordings())))
	out.Write(numRecordings[:read])

	// Write all child recordings
	newOffset := offset + len(recording.CaptureStreams())
	for _, rec := range recording.Recordings() {
		updatedOffset, recordingData := recurseRecordingToBytes(rec, keyMappingToIndex, encodingBlocks, streamIndexToEncoderUsedIndex, newOffset)
		newOffset = updatedOffset
		out.Write(rapbinary.BytesArrayToBytes(recordingData))
	}

	return newOffset, out.Bytes()
}

func (w Writer) Write(recording format.Recording) (int, error) {
	if recording == nil {
		panic(errors.New("can not write nil recording"))
	}

	encoderMappings, _, err := w.evaluateStreams(recording, 0)
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
		header, streamsEncoded, err := val.encoder.Encode(val.streams)
		if err != nil {
			return totalBytesWritten, err
		}

		for i, order := range val.streamOrder {
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
