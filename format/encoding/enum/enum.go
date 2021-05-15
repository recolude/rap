package enum

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/enum"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type StorageTechnique int

const (
	// Raw64 encodes time with 64 bit precision
	Raw64 StorageTechnique = iota

	// Raw32 encodes time with 32 bit precision
	Raw32
)

type Encoder struct {
	technique StorageTechnique
}

func NewEncoder(technique StorageTechnique) Encoder {
	return Encoder{technique}
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
		enmstr := stream.(enum.Stream)
		streamDataBuffers[bufferIndex].Write(rapbinary.StringArrayToBytes(enmstr.EnumMembers()))

		// Write technique
		streamDataBuffers[bufferIndex].WriteByte(byte(p.technique))

		// Write Num Captures
		numCaptures := make([]byte, 4)
		read := binary.PutUvarint(numCaptures, uint64(len(stream.Captures())))
		streamDataBuffers[bufferIndex].Write(numCaptures[:read])

		for _, c := range stream.Captures() {
			enumCapture, ok := c.(enum.Capture)
			if !ok {
				return nil, nil, errors.New("capture is not of type enum")
			}

			switch p.technique {
			case Raw64:
				binary.Write(&streamDataBuffers[bufferIndex], binary.LittleEndian, enumCapture.Time())
			case Raw32:
				binary.Write(&streamDataBuffers[bufferIndex], binary.LittleEndian, float32(enumCapture.Time()))
			}

			valueBuf := make([]byte, 4)
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

func (p Encoder) Decode(header []byte, streamData []byte) (format.CaptureCollection, error) {
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

	// Read Storage Technique
	typeByte, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	encodingTechnique := StorageTechnique(typeByte)

	// Read Num Captures
	numCaptures, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}

	captures := make([]enum.Capture, numCaptures)
	for i := 0; i < int(numCaptures); i++ {
		var time float64

		switch encodingTechnique {
		case Raw64:
			binary.Read(buf, binary.LittleEndian, &time)

		case Raw32:
			var time32 float32
			binary.Read(buf, binary.LittleEndian, &time32)
			time = float64(time32)
		}

		value, err := binary.ReadUvarint(buf)
		if err != nil {
			return nil, err
		}
		captures[i] = enum.NewCapture(time, int(value))
	}

	return enum.NewStream(streamName, enumMembers, captures), nil
}
