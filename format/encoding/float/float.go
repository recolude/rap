package float

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

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
		binary.Write(out, binary.LittleEndian, floatCapture.Time())
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
		binary.Write(out, binary.LittleEndian, float32(floatCapture.Time()))
		binary.Write(out, binary.LittleEndian, float32(floatCapture.Value()))
	}
	return nil
}

// func encodeBST16(out io.Writer, captures []format.Capture) error {
// 	minTime := math.Inf(1)
// 	maxTime := math.Inf(-1)
// 	minVal := math.Inf(1)
// 	maxVal := math.Inf(-1)
// 	for _, c := range captures {
// 		floatCapture, ok := c.(float.Capture)
// 		if !ok {
// 			return errors.New("capture is not of type float")
// 		}
// 		if floatCapture.Time() > maxTime {
// 			maxTime = floatCapture.Time()
// 		}
// 		if floatCapture.Time() < minTime {
// 			minTime = floatCapture.Time()
// 		}

// 		if floatCapture.Value() > maxVal {
// 			maxVal = floatCapture.Value()
// 		}
// 		if floatCapture.Value() < minVal {
// 			minVal = floatCapture.Value()
// 		}

// 		binary.Write(out, binary.LittleEndian, float32(floatCapture.Time()))
// 		binary.Write(out, binary.LittleEndian, float32(floatCapture.Value()))
// 	}

// 	totalledQuantizedDuration := minTime
// 	for _, c := range captures {
// 		capture := c.(float.Capture)

// 		// Write Time
// 		duration := capture.Time() - totalledQuantizedDuration
// 		rapbinary.UnsignedFloatBSTToBytes(duration, 0, maxTimeDifference, timeBuffer)
// 		_, err := collectionData.Write(timeBuffer)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return nil
// }

func (p Encoder) Encode(streams []format.CaptureCollection) ([]byte, [][]byte, error) {
	streamDataBuffers := make([]bytes.Buffer, len(streams))
	for bufferIndex, stream := range streams {
		// Write Stream Name
		streamDataBuffers[bufferIndex].Write(rapbinary.StringToBytes(stream.Name()))

		// Write technique
		streamDataBuffers[bufferIndex].WriteByte(byte(p.technique))

		// Write Num Captures
		numCaptures := make([]byte, binary.MaxVarintLen64)
		read := binary.PutUvarint(numCaptures, uint64(len(stream.Captures())))
		streamDataBuffers[bufferIndex].Write(numCaptures[:read])

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

	captures := make([]float.Capture, numCaptures)
	for i := 0; i < int(numCaptures); i++ {
		var time float64
		var value float64

		switch encodingTechnique {
		case Raw64:
			binary.Read(buf, binary.LittleEndian, &time)
			binary.Read(buf, binary.LittleEndian, &value)

		case Raw32:
			var time32 float32
			binary.Read(buf, binary.LittleEndian, &time32)
			time = float64(time32)
			var value32 float32
			binary.Read(buf, binary.LittleEndian, &value32)
			value = float64(value32)
		}

		captures[i] = float.NewCapture(time, value)
	}

	return float.NewCollection(streamName, captures), nil
}
