package euler

import (
	"bytes"
	"fmt"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/euler"
)

type StorageTechnique int

const (
	// Raw64 encodes all values at fullest precision, costing 192 bits per
	// capture
	Raw64 StorageTechnique = iota

	// Raw32 encodes all values at 32bit precision, costing 96 bits per
	// capture
	Raw32

	// Raw16 stores all values at 16 bit precision, costiting 48 bits per
	// capture
	Raw16
)

type Encoder struct {
	technique StorageTechnique
}

func NewEncoder(technique StorageTechnique) Encoder {
	return Encoder{technique: technique}
}

func (p Encoder) encode(stream format.CaptureCollection) ([]byte, error) {
	streamData := new(bytes.Buffer)

	castedCaptureData := make([]euler.Capture, len(stream.Captures()))
	for i, c := range stream.Captures() {
		castedCaptureData[i] = c.(euler.Capture)
	}

	streamData.WriteByte(byte(p.technique))

	switch p.technique {
	case Raw64:
		streamData.Write(encodeRaw64(castedCaptureData))
		break
	case Raw32:
		streamData.Write(encodeRaw32(castedCaptureData))
		break
	case Raw16:
		streamData.Write(encodeRaw16(castedCaptureData))
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

func decode(name string, data []byte, times []float64) (format.CaptureCollection, error) {
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
		return euler.NewCollection(name, captures), nil

	case Raw32:
		captures, err := decodeRaw32(reader, times)
		if err != nil {
			return nil, err
		}
		return euler.NewCollection(name, captures), nil

	case Raw16:
		captures, err := decodeRaw16(reader, times)
		if err != nil {
			return nil, err
		}
		return euler.NewCollection(name, captures), nil
	}

	return nil, fmt.Errorf("Unknown euler encoding technique: %d", int(encodingTechnique))
}

func (p Encoder) Decode(name string, header []byte, streamData []byte, times []float64) (format.CaptureCollection, error) {
	return decode(name, streamData, times)
}

func (p Encoder) Accepts(stream format.CaptureCollection) bool {
	return stream.Signature() == "recolude.euler"
}

func (p Encoder) Signature() string {
	return "recolude.euler"
}

func (p Encoder) Version() uint {
	return 0
}
