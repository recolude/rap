package float

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/float"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type StorageTechnique int

const (
	// Raw encodes with maximum precision
	Raw64 StorageTechnique = iota

	// Raw32 truncates all number values to 32bit
	Raw32

	// BST16 encodes within 16 bits.
	BST16
)

type Encoder struct {
	technique StorageTechnique
}

func NewEncoder(technique StorageTechnique) Encoder {
	return Encoder{technique}
}

func (p Encoder) Accepts(stream format.CaptureCollection) bool {
	return stream.Signature() == "recolude.float"
}

func (p Encoder) Signature() string {
	return "recolude.float"
}

func (p Encoder) Version() uint {
	return 0
}

func encode64(out io.Writer, captures []format.Capture) error {
	for _, c := range captures {
		floatCapture, ok := c.(float.Capture)
		if !ok {
			return errors.New("capture is not of type float")
		}
		binary.Write(out, binary.LittleEndian, floatCapture.Value())
	}
	return nil
}

func encode32(out io.Writer, captures []format.Capture) error {
	for _, c := range captures {
		floatCapture, ok := c.(float.Capture)
		if !ok {
			return errors.New("capture is not of type float")
		}
		binary.Write(out, binary.LittleEndian, float32(floatCapture.Value()))
	}
	return nil
}

func encodeBST16(out io.Writer, captures []format.Capture) error {
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	for _, c := range captures {
		floatCapture, ok := c.(float.Capture)
		if !ok {
			return errors.New("capture is not of type float")
		}

		if floatCapture.Value() > maxVal {
			maxVal = floatCapture.Value()
		}
		if floatCapture.Value() < minVal {
			minVal = floatCapture.Value()
		}
	}
	binary.Write(out, binary.LittleEndian, float32(minVal))
	binary.Write(out, binary.LittleEndian, float32(maxVal))

	valBuffer := make([]byte, 2)
	for _, c := range captures {
		capture := c.(float.Capture)
		rapbinary.UnsignedFloatBSTToBytes(capture.Value(), minVal, maxVal-minVal, valBuffer)
		_, err := out.Write(valBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Encoder) Encode(streams []format.CaptureCollection) ([]byte, [][]byte, error) {
	streamDataBuffers := make([]bytes.Buffer, len(streams))
	for bufferIndex, stream := range streams {

		// Write technique
		streamDataBuffers[bufferIndex].WriteByte(byte(p.technique))

		switch p.technique {
		case Raw64:
			err := encode64(&streamDataBuffers[bufferIndex], stream.Captures())
			if err != nil {
				return nil, nil, err
			}
			break

		case Raw32:
			err := encode32(&streamDataBuffers[bufferIndex], stream.Captures())
			if err != nil {
				return nil, nil, err
			}
			break

		case BST16:
			err := encodeBST16(&streamDataBuffers[bufferIndex], stream.Captures())
			if err != nil {
				return nil, nil, err
			}
			break
		}
	}

	streamData := make([][]byte, len(streams))
	for i, buffer := range streamDataBuffers {
		streamData[i] = buffer.Bytes()
	}

	return nil, streamData, nil
}

func decodeBST16(in io.Reader, times []float64) ([]float.Capture, error) {
	var min float32
	var max float32

	err := binary.Read(in, binary.LittleEndian, &min)
	if err != nil {
		return nil, err
	}

	err = binary.Read(in, binary.LittleEndian, &max)
	if err != nil {
		return nil, err
	}

	captures := make([]float.Capture, len(times))

	buffer := make([]byte, 2)
	for i, time := range times {
		in.Read(buffer)
		value := rapbinary.BytesToUnisngedFloatBST(float64(min), float64(max-min), buffer)
		captures[i] = float.NewCapture(time, value)
	}

	return captures, nil
}

func decodeRaw64(in io.Reader, times []float64) ([]float.Capture, error) {
	captures := make([]float.Capture, len(times))
	var value float64
	for i, time := range times {
		err := binary.Read(in, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		captures[i] = float.NewCapture(time, value)
	}
	return captures, nil
}

func decodeRaw32(in io.Reader, times []float64) ([]float.Capture, error) {
	captures := make([]float.Capture, len(times))
	var value32 float32
	for i, time := range times {
		err := binary.Read(in, binary.LittleEndian, &value32)
		if err != nil {
			return nil, err
		}
		captures[i] = float.NewCapture(time, float64(value32))
	}
	return captures, nil
}

func (p Encoder) Decode(name string, header []byte, streamData []byte, times []float64) (format.CaptureCollection, error) {
	buf := bytes.NewBuffer(streamData)

	// Read Storage Technique
	typeByte, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	encodingTechnique := StorageTechnique(typeByte)

	var captures []float.Capture
	switch encodingTechnique {
	case Raw64:
		captures, err = decodeRaw64(buf, times)
		if err != nil {
			return nil, err
		}

	case Raw32:
		captures, err = decodeRaw32(buf, times)
		if err != nil {
			return nil, err
		}

	case BST16:
		captures, err = decodeBST16(buf, times)
		if err != nil {
			return nil, err
		}
	}

	return float.NewCollection(name, captures), nil
}
