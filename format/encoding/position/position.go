package position

import (
	"bytes"
	"fmt"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/streams/position"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type StorageTechnique int

const (
	// Raw64 encodes all values at fullest precision, costing 256 bits per
	// capture
	Raw64 StorageTechnique = iota

	// Raw32 encodes all values at 32bit precision, costing 128 bits per
	// capture
	Raw32

	// Oct48 stores all values in a oct tree of depth 16, costing 64 bits per
	// capture (time is stored in 16 bits)
	Oct48

	// Oct24 stores all values in a oct tree of depth 8, costing 40 bits per
	// capture (time is stored in 16 bits)
	Oct24
)

type Encoder struct {
	technique StorageTechnique
}

func NewEncoder(technique StorageTechnique) Encoder {
	return Encoder{technique: technique}
}

func (p Encoder) encode(stream format.CaptureStream) ([]byte, error) {
	streamData := new(bytes.Buffer)

	streamData.Write(rapbinary.StringToBytes(stream.Name()))

	castedCaptureData := make([]position.Capture, len(stream.Captures()))
	for i, c := range stream.Captures() {
		castedCaptureData[i] = c.(position.Capture)
	}

	streamData.WriteByte(byte(p.technique))

	switch p.technique {
	case Raw64:
		streamData.Write(encodeRaw64(castedCaptureData))
		break
	case Raw32:
		streamData.Write(encodeRaw32(castedCaptureData))
		break

	case Oct24:
		d, err := encodeOct24(castedCaptureData)
		if err != nil {
			return nil, err
		}
		streamData.Write(d)
		break
	}

	return streamData.Bytes(), nil
}

func (p Encoder) Encode(streams []format.CaptureStream) ([]byte, [][]byte, error) {
	allStreamData := make([][]byte, len(streams))

	for i, stream := range streams {
		s, err := p.encode(stream)
		if err != nil {
			return nil, nil, err
		}
		allStreamData[i] = s
	}

	return nil, allStreamData, nil
}

func decode(data []byte) (format.CaptureStream, error) {
	reader := bytes.NewReader(data)

	name, _, err := rapbinary.ReadString(reader)
	if err != nil {
		return nil, err
	}

	typeByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	encodingTechnique := StorageTechnique(typeByte)

	switch encodingTechnique {
	case Raw64:
		captures, err := decodeRaw64(reader)
		if err != nil {
			return nil, err
		}
		return position.NewStream(name, captures), nil

	case Raw32:
		captures, err := decodeRaw32(reader)
		if err != nil {
			return nil, err
		}
		return position.NewStream(name, captures), nil

	case Oct24:
		captures, err := decodeOct24(reader)
		if err != nil {
			return nil, err
		}
		return position.NewStream(name, captures), nil
	}

	return nil, fmt.Errorf("Unknown positional encoding technique: %d", int(encodingTechnique))
}

func (p Encoder) Decode(header []byte, streamData []byte) (format.CaptureStream, error) {
	return decode(streamData)
}

func (p Encoder) Accepts(stream format.CaptureStream) bool {
	return stream.Signature() == "recolude.position"
}

func (p Encoder) Signature() string {
	return "recolude.position"
}

func (p Encoder) Version() uint {
	return 0
}
