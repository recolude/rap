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
	streamDataBuffers := make([]bytes.Buffer, len(streams))
	for bufferIndex, stream := range streams {
		// Write Stream Name
		streamDataBuffers[bufferIndex].Write(rapbinary.StringToBytes(stream.Name()))

		// Write Enum Members
		enmstr := stream.(enum.Collection)
		streamDataBuffers[bufferIndex].Write(rapbinary.StringArrayToBytes(enmstr.EnumMembers()))

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

	streamData := make([][]byte, len(streams))
	for i, buffer := range streamDataBuffers {
		streamData[i] = buffer.Bytes()
	}

	return nil, streamData, nil
}

func (p Encoder) Decode(header []byte, streamData []byte, times []float64) (format.CaptureCollection, error) {
	buf := bytes.NewBuffer(streamData)

	// Read Name
	streamName, _, err := rapbinary.ReadString(buf)
	if err != nil {
		return nil, err
	}

	enumMembers, _, err := rapbinary.ReadStringArray(buf)
	if err != nil {
		return nil, err
	}

	captures := make([]enum.Capture, len(times))
	for i := 0; i < len(times); i++ {
		value, err := binary.ReadUvarint(buf)
		if err != nil {
			return nil, err
		}
		captures[i] = enum.NewCapture(times[i], int(value))
	}

	return enum.NewCollection(streamName, enumMembers, captures), nil
}
