package euler

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/recolude/rap/format/collection/euler"
	binaryutil "github.com/recolude/rap/internal/io/binary"
)

func encodeRaw16(captures []euler.Capture) []byte {
	streamData := new(bytes.Buffer)

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

	// Num captures
	buf := make([]byte, 8)
	size := binary.PutUvarint(buf, uint64(len(captures)))
	streamData.Write(buf[:size])

	if len(captures) == 0 {
		return streamData.Bytes()
	}

	binary.Write(streamData, binary.LittleEndian, float32(startingTime))

	if len(captures) == 1 {
		binary.Write(streamData, binary.LittleEndian, float32(captures[0].EulerZXY().X()))
		binary.Write(streamData, binary.LittleEndian, float32(captures[0].EulerZXY().Y()))
		binary.Write(streamData, binary.LittleEndian, float32(captures[0].EulerZXY().Z()))
		return streamData.Bytes()
	}

	binary.Write(streamData, binary.LittleEndian, float32(maxTimeDifference))

	totalledQuantizedDuration := startingTime
	buffer2Byes := make([]byte, 2)
	for _, capture := range captures {
		// Write Time
		duration := capture.Time() - totalledQuantizedDuration
		binaryutil.UnsignedFloatBSTToBytes(duration, 0, maxTimeDifference, buffer2Byes)
		streamData.Write(buffer2Byes)

		// Read back quantized time to fix drifting
		totalledQuantizedDuration += binaryutil.BytesToUnisngedFloatBST(0, maxTimeDifference, buffer2Byes)

		binaryutil.UnsignedFloatBSTToBytes(capture.EulerZXY().X(), 0, 360, buffer2Byes)
		streamData.Write(buffer2Byes)

		binaryutil.UnsignedFloatBSTToBytes(capture.EulerZXY().Y(), 0, 360, buffer2Byes)
		streamData.Write(buffer2Byes)

		binaryutil.UnsignedFloatBSTToBytes(capture.EulerZXY().Z(), 0, 360, buffer2Byes)
		streamData.Write(buffer2Byes)
	}
	return streamData.Bytes()
}

func decodeRaw16(streamData *bytes.Reader) ([]euler.Capture, error) {
	// streamData := bytes.NewReader(captureData)

	numCaptures, err := binary.ReadUvarint(streamData)
	if err != nil {
		return nil, err
	}

	if numCaptures == 0 {
		return make([]euler.Capture, 0), nil
	}

	var startTime float32
	err = binary.Read(streamData, binary.LittleEndian, &startTime)
	if err != nil {
		return nil, err
	}

	if numCaptures == 1 {
		var posX float32
		var posY float32
		var posZ float32
		err = binary.Read(streamData, binary.LittleEndian, &posX)
		err = binary.Read(streamData, binary.LittleEndian, &posY)
		err = binary.Read(streamData, binary.LittleEndian, &posZ)
		return []euler.Capture{euler.NewEulerZXYCapture(float64(startTime), float64(posX), float64(posY), float64(posZ))}, err
	}

	var maxTimeDifference float32
	err = binary.Read(streamData, binary.LittleEndian, &maxTimeDifference)
	if err != nil {
		return nil, err
	}

	captures := make([]euler.Capture, numCaptures)
	buffer := make([]byte, 2)
	currentTime := float64(startTime)

	for i := 0; i < int(numCaptures); i++ {
		streamData.Read(buffer)
		time := binaryutil.BytesToUnisngedFloatBST(0, float64(maxTimeDifference), buffer)
		currentTime += time

		streamData.Read(buffer)
		x := binaryutil.BytesToUnisngedFloatBST(0, 360, buffer)
		streamData.Read(buffer)
		y := binaryutil.BytesToUnisngedFloatBST(0, 360, buffer)
		streamData.Read(buffer)
		z := binaryutil.BytesToUnisngedFloatBST(0, 360, buffer)
		captures[i] = euler.NewEulerZXYCapture(currentTime, float64(x), float64(y), float64(z))
	}

	return captures, nil
}
