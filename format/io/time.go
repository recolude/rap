package io

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/recolude/rap/format"
	rapbinary "github.com/recolude/rap/internal/io/binary"
)

type TimeStorageTechnique int

const (
	// Raw64 encodes time with 64 bit precision
	Raw64 TimeStorageTechnique = iota

	// Raw32 encodes time with 32 bit precision
	Raw32

	BST16
)

func encodeTime64(out io.Writer, captures []format.Capture) error {
	for _, c := range captures {
		err := binary.Write(out, binary.LittleEndian, c.Time())
		if err != nil {
			return err
		}
	}
	return nil
}

func encodeTime32(out io.Writer, captures []format.Capture) error {
	for _, c := range captures {
		err := binary.Write(out, binary.LittleEndian, float32(c.Time()))
		if err != nil {
			return err
		}
	}
	return nil
}

func encodeTimeBST16(out io.Writer, captures []format.Capture) error {
	if len(captures) == 0 {
		return nil
	}

	startingTime := math.Inf(1)
	endingTime := math.Inf(-1)
	maxTimeDifference := math.Inf(-1)

	for i, capture := range captures {
		if capture.Time() < startingTime {
			startingTime = capture.Time()
		}
		if capture.Time() > endingTime {
			endingTime = capture.Time()
		}

		if i > 0 {
			timeDifference := capture.Time() - captures[i-1].Time()
			if timeDifference > maxTimeDifference {
				maxTimeDifference = timeDifference
			}
		}

	}

	binary.Write(out, binary.LittleEndian, float32(startingTime))

	if len(captures) == 1 {
		return nil
	}
	binary.Write(out, binary.LittleEndian, float32(maxTimeDifference))

	totalledQuantizedDuration := startingTime
	buffer2Byes := make([]byte, 2)
	for i := 1; i < len(captures); i++ {
		// Write Time
		duration := captures[i].Time() - totalledQuantizedDuration
		rapbinary.UnsignedFloatBSTToBytes(duration, 0, maxTimeDifference, buffer2Byes)
		out.Write(buffer2Byes)

		// Read back quantized time to fix drifting
		totalledQuantizedDuration += rapbinary.BytesToUnisngedFloatBST(0, maxTimeDifference, buffer2Byes)
	}
	return nil
}

func encodeTime(technique TimeStorageTechnique, out io.Writer, captures []format.Capture) (int, error) {
	dataBuffer := bytes.Buffer{}

	// Write technique
	dataBuffer.WriteByte(byte(technique))

	// Write Num Captures
	numCaptures := make([]byte, binary.MaxVarintLen64)
	read := binary.PutUvarint(numCaptures, uint64(len(captures)))
	dataBuffer.Write(numCaptures[:read])

	switch technique {
	case Raw64:
		err := encodeTime64(&dataBuffer, captures)
		if err != nil {
			return 0, err
		}
		break

	case Raw32:
		err := encodeTime32(&dataBuffer, captures)
		if err != nil {
			return 0, err
		}
		break

	case BST16:
		err := encodeTimeBST16(&dataBuffer, captures)
		if err != nil {
			return 0, err
		}
		break
	}
	return out.Write(dataBuffer.Bytes())
}

func decodeTime64(in io.Reader, numCaptures int) ([]float64, error) {
	times := make([]float64, numCaptures)
	for i := 0; i < int(numCaptures); i++ {
		var time float64

		binary.Read(in, binary.LittleEndian, &time)

		times[i] = time
	}
	return times, nil
}

func decodeTime32(in io.Reader, numCaptures int) ([]float64, error) {
	times := make([]float64, numCaptures)
	for i := 0; i < int(numCaptures); i++ {

		var time32 float32
		binary.Read(in, binary.LittleEndian, &time32)

		times[i] = float64(time32)
	}

	return times, nil
}

func decodeTimeBST16(in io.Reader, numCaptures int) ([]float64, error) {
	if numCaptures == 0 {
		return make([]float64, 0), nil
	}

	var startTime float32
	err := binary.Read(in, binary.LittleEndian, &startTime)
	if err != nil {
		return nil, err
	}

	if numCaptures == 1 {
		return []float64{float64(startTime)}, nil
	}

	var maxTimeDifference float32
	err = binary.Read(in, binary.LittleEndian, &maxTimeDifference)
	if err != nil {
		return nil, err
	}

	captures := make([]float64, numCaptures)
	captures[0] = float64(startTime)
	buffer := make([]byte, 2)
	currentTime := float64(startTime)

	for i := 1; i < int(numCaptures); i++ {
		in.Read(buffer)
		time := rapbinary.BytesToUnisngedFloatBST(0, float64(maxTimeDifference), buffer)
		currentTime += time

		captures[i] = currentTime
	}

	return captures, nil
}

func decodeTime(in io.Reader) ([]float64, error) {
	typeByte := []byte{0}

	// Read Storage Technique
	_, err := in.Read(typeByte)
	if err != nil {
		return nil, err
	}
	encodingTechnique := TimeStorageTechnique(typeByte[0])

	// Read Num Captures
	numCaptures, _, err := rapbinary.ReadUvarint(in)
	if err != nil {
		return nil, err
	}

	switch encodingTechnique {
	case Raw64:
		return decodeTime64(in, int(numCaptures))

	case Raw32:
		return decodeTime32(in, int(numCaptures))

	case BST16:
		return decodeTimeBST16(in, int(numCaptures))
	}

	return nil, fmt.Errorf("unrecognized time encoding: %d", encodingTechnique)
}
