package enum

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/enum"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type Encoder struct{}

func NewEncoder() Encoder {
	return Encoder{}
}

func (p Encoder) Accepts(stream format.CaptureCollection) bool {
	return stream.Signature() == "recolude.enum"
}

func (p Encoder) Signature() string {
	return "recolude.enum"
}

func (p Encoder) Version() uint {
	return 0
}

func (p Encoder) Encode(streams []format.CaptureCollection) ([]byte, [][]byte, error) {
	enumMemberMapping := make(map[string]int)

	streamDataBuffers := make([]bytes.Buffer, len(streams))

	for bufferIndex, stream := range streams {
		enmstr := stream.(enum.Collection)

		// Build mapping from enum member to index in header
		indexMapping := make([]uint, len(enmstr.EnumMembers()))
		for i, member := range enmstr.EnumMembers() {
			if val, ok := enumMemberMapping[member]; ok {
				indexMapping[i] = uint(val)
			} else {
				indexMapping[i] = uint(len(enumMemberMapping))
				enumMemberMapping[member] = len(enumMemberMapping)
			}
		}

		// Write Enum Members indexes
		streamDataBuffers[bufferIndex].Write(rapbinary.UvarintArrayToBytes(indexMapping))

		for _, c := range stream.Captures() {
			enumCapture, ok := c.(enum.Capture)
			if !ok {
				return nil, nil, errors.New("capture is not of type enum")
			}

			valueBuf := make([]byte, binary.MaxVarintLen64)
			read := binary.PutUvarint(valueBuf, uint64(enumCapture.Value()))
			streamDataBuffers[bufferIndex].Write(valueBuf[:read])
		}
	}

	// Build header
	headerBuffer := bytes.Buffer{}
	headerMembers := make([]string, len(enumMemberMapping))
	for key, val := range enumMemberMapping {
		headerMembers[val] = key
	}
	headerBuffer.Write(rapbinary.StringArrayToBytes(headerMembers))

	streamData := make([][]byte, len(streams))
	for i, buffer := range streamDataBuffers {
		streamData[i] = buffer.Bytes()
	}

	return headerBuffer.Bytes(), streamData, nil
}

func (p Encoder) Decode(name string, header []byte, streamData []byte, times []float64) (format.CaptureCollection, error) {
	allEnumMembers, _, err := rapbinary.ReadStringArray(bytes.NewReader(header))
	if err != nil {
		return nil, err
	}

	reader := rapbinary.NewErrReader(bytes.NewBuffer(streamData))

	enumMemberIndexes, _, _ := rapbinary.ReadUvarIntArray(reader)

	enumMembers := make([]string, len(enumMemberIndexes))
	for i, indeces := range enumMemberIndexes {
		enumMembers[i] = allEnumMembers[indeces]
	}

	captures := make([]enum.Capture, len(times))
	for i := 0; i < len(times); i++ {
		value, _ := binary.ReadUvarint(reader)
		captures[i] = enum.NewCapture(times[i], int(value))
	}

	return enum.NewCollection(name, enumMembers, captures), reader.Error()
}
