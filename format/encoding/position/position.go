package position

import (
	"bytes"
	"fmt"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
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

func (p Encoder) encode(stream format.CaptureCollection) ([]byte, error) {
	streamData := new(bytes.Buffer)

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

	case Oct48:
		d, err := encodeOct48(castedCaptureData)
		if err != nil {
			return nil, err
		}
		streamData.Write(d)
		break
	}

	return streamData.Bytes(), nil
}

func (p Encoder) Encode(streams []format.CaptureCollection) ([]byte, [][]byte, error) {
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

func decode(streamName string, data []byte, times []float64) (format.CaptureCollection, error) {
	reader := bytes.NewReader(data)

	typeByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	encodingTechnique := StorageTechnique(typeByte)

	switch encodingTechnique {
	case Raw64:
		captures, err := decodeRaw64(reader, times)
		if err != nil {
			return nil, err
		}
		return position.NewCollection(streamName, captures), nil

	case Raw32:
		captures, err := decodeRaw32(reader, times)
		if err != nil {
			return nil, err
		}
		return position.NewCollection(streamName, captures), nil

	case Oct24:
		captures, err := decodeOct24(reader, times)
		if err != nil {
			return nil, err
		}
		return position.NewCollection(streamName, captures), nil

	case Oct48:
		captures, err := decodeOct48(reader, times)
		if err != nil {
			return nil, err
		}
		return position.NewCollection(streamName, captures), nil
	}

	return nil, fmt.Errorf("Unknown positional encoding technique: %d", int(encodingTechnique))
}

func (p Encoder) Decode(streamName string, header []byte, streamData []byte, times []float64) (format.CaptureCollection, error) {
	return decode(streamName, streamData, times)
}

func (p Encoder) Accepts(stream format.CaptureCollection) bool {
	return stream.Signature() == "recolude.position"
}

func (p Encoder) Signature() string {
	return "recolude.position"
}

func (p Encoder) Version() uint {
	return 0
}
