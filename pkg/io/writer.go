package io

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	rapbinary "github.com/recolude/rap/internal/io/binary"
	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/encoding"
)

type encoderStreamMapping struct {
	encoder encoding.Encoder
	streams []data.CaptureStream
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

func (w Writer) evaluateStreams(recording data.Recording) ([]encoderStreamMapping, error) {
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
				streamsSatisfied[streamIndex] = true
			}
		}

		if len(mapping.streams) > 0 {
			mappings = append(mappings, mapping)
		}
	}

	for i, stream := range recording.CaptureStreams() {
		if streamsSatisfied[i] == false {
			return nil, fmt.Errorf("no encoder registered to handle stream: %s", stream.Signature())
		}
	}

	return mappings, nil
}

func writeMetadata(out io.Writer, metadata map[string]string) (int, error) {
	metadataBlock := make([]string, len(metadata)*2)
	i := 0
	for key, val := range metadata {
		metadataBlock[i] = key
		metadataBlock[i+1] = val
		i += 2
	}
	return out.Write(rapbinary.StringArrayToBytes(metadataBlock))
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

func (w Writer) Write(recording data.Recording) (int, error) {
	if recording == nil {
		panic(errors.New("can not write nil recording"))
	}

	encoderMappings, err := w.evaluateStreams(recording)
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

	// Write name
	written, err = w.out.Write(rapbinary.StringToBytes(recording.Name()))
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

	// Write metadata
	written, err = writeMetadata(w.out, recording.Metadata())
	totalBytesWritten += written
	if err != nil {
		return totalBytesWritten, err
	}

	for _, val := range encoderMappings {
		header, streamsEncoded, err := val.encoder.Encode(val.streams)
		if err != nil {
			return totalBytesWritten, err
		}

		// Write header
		written, err = w.out.Write(rapbinary.BytesArrayToBytes(header))
		totalBytesWritten += written
		if err != nil {
			return totalBytesWritten, err
		}

		// Write number of streams
		numStreams := make([]byte, 4)
		read := binary.PutUvarint(numStreams, uint64(len(streamsEncoded)))
		written, err = w.out.Write(numStreams[:read])
		totalBytesWritten += written
		if err != nil {
			return totalBytesWritten, err
		}

		// Write all streams
		for _, encodedStream := range streamsEncoded {
			written, err = w.out.Write(rapbinary.BytesArrayToBytes(encodedStream))
			totalBytesWritten += written
			if err != nil {
				return totalBytesWritten, err
			}
		}
	}

	return totalBytesWritten, nil
}
