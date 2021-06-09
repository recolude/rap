package euler

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/recolude/rap/format/collection/euler"
	binaryutil "github.com/recolude/rap/internal/io/binary"
)

func wrapEulerAngle(angle float64) float64 {
	return angle - (360 * math.Floor(angle/360))
}

func encodeRaw16(captures []euler.Capture) []byte {
	streamData := new(bytes.Buffer)

	if len(captures) == 0 {
		return streamData.Bytes()
	}

	if len(captures) == 1 {
		binary.Write(streamData, binary.LittleEndian, float32(wrapEulerAngle(captures[0].EulerZXY().X())))
		binary.Write(streamData, binary.LittleEndian, float32(wrapEulerAngle(captures[0].EulerZXY().Y())))
		binary.Write(streamData, binary.LittleEndian, float32(wrapEulerAngle(captures[0].EulerZXY().Z())))
		return streamData.Bytes()
	}

	buffer2Byes := make([]byte, 2)
	for _, capture := range captures {
		binaryutil.UnsignedFloatBSTToBytes(wrapEulerAngle(capture.EulerZXY().X()), 0, 360, buffer2Byes)
		streamData.Write(buffer2Byes)

		binaryutil.UnsignedFloatBSTToBytes(wrapEulerAngle(capture.EulerZXY().Y()), 0, 360, buffer2Byes)
		streamData.Write(buffer2Byes)

		binaryutil.UnsignedFloatBSTToBytes(wrapEulerAngle(capture.EulerZXY().Z()), 0, 360, buffer2Byes)
		streamData.Write(buffer2Byes)
	}
	return streamData.Bytes()
}

func decodeRaw16(streamData *bytes.Reader, times []float64) ([]euler.Capture, error) {
	// streamData := bytes.NewReader(captureData)

	if len(times) == 1 {
		var posX float32
		var posY float32
		var posZ float32
		err := binary.Read(streamData, binary.LittleEndian, &posX)
		err = binary.Read(streamData, binary.LittleEndian, &posY)
		err = binary.Read(streamData, binary.LittleEndian, &posZ)
		return []euler.Capture{euler.NewEulerZXYCapture(times[0], float64(posX), float64(posY), float64(posZ))}, err
	}

	captures := make([]euler.Capture, len(times))
	buffer := make([]byte, 2)

	for i := 0; i < len(times); i++ {
		streamData.Read(buffer)
		x := binaryutil.BytesToUnisngedFloatBST(0, 360, buffer)
		streamData.Read(buffer)
		y := binaryutil.BytesToUnisngedFloatBST(0, 360, buffer)
		streamData.Read(buffer)
		z := binaryutil.BytesToUnisngedFloatBST(0, 360, buffer)
		captures[i] = euler.NewEulerZXYCapture(times[i], float64(x), float64(y), float64(z))
	}

	return captures, nil
}
