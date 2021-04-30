package position

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	rapbinary "github.com/recolude/rap/internal/io/binary"
	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/streams/position"
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
)

type Encoder struct {
	technique StorageTechnique
}

func NewEncoder(technique StorageTechnique) Encoder {
	return Encoder{technique: technique}
}

func encodeRaw64(captures []position.Capture) []byte {
	streamData := new(bytes.Buffer)

	buf := make([]byte, 8)
	size := binary.PutUvarint(buf, uint64(len(captures)))
	streamData.Write(buf[:size])

	for _, capture := range captures {
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Time()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Position().X()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Position().Y()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Position().Z()))
		streamData.Write(buf)
	}
	return streamData.Bytes()
}

func (p Encoder) encode(stream data.CaptureStream) ([]byte, error) {
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
	}

	// buf := make([]byte, 8)

	// // Write Start
	// binary.LittleEndian.PutUint64(buf, math.Float64bits(stream.Captures()[0].Time()))
	// _, err := streamData.Write(buf)
	// if err != nil {
	// 	return nil, err
	// }

	// // Write duration
	// duration := streamDuration(stream)
	// binary.LittleEndian.PutUint64(buf, math.Float64bits(duration))
	// _, err = streamData.Write(buf)
	// if err != nil {
	// 	return nil, err
	// }

	return streamData.Bytes(), nil
}

func (p Encoder) Encode(streams []data.CaptureStream) ([]byte, [][]byte, error) {
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

func decodeRaw64(streamData *bytes.Reader) ([]position.Capture, error) {
	// streamData := bytes.NewReader(captureData)

	numCaptures, err := binary.ReadUvarint(streamData)
	if err != nil {
		return nil, err
	}

	captures := make([]position.Capture, numCaptures)
	buf := make([]byte, 8)
	for i := 0; i < int(numCaptures); i++ {
		_, err = streamData.Read(buf)
		if err != nil {
			return nil, err
		}
		time := math.Float64frombits(binary.LittleEndian.Uint64(buf))

		_, err = streamData.Read(buf)
		if err != nil {
			return nil, err
		}
		x := math.Float64frombits(binary.LittleEndian.Uint64(buf))

		_, err = streamData.Read(buf)
		if err != nil {
			return nil, err
		}
		y := math.Float64frombits(binary.LittleEndian.Uint64(buf))

		_, err = streamData.Read(buf)
		if err != nil {
			return nil, err
		}
		z := math.Float64frombits(binary.LittleEndian.Uint64(buf))

		captures[i] = position.NewCapture(time, x, y, z)
	}

	return captures, nil
}

func decode(data []byte) (data.CaptureStream, error) {
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
	}

	return nil, fmt.Errorf("Unknown positional encoding technique: %d", int(encodingTechnique))
}

func (p Encoder) Decode(header []byte, streamData []byte) (data.CaptureStream, error) {
	return decode(streamData)
}

func (p Encoder) Accepts(stream data.CaptureStream) bool {
	return stream.Signature() == "recolude.position"
}

func (p Encoder) Signature() string {
	return "recolude.position"
}

func (p Encoder) Version() uint {
	return 0
}
