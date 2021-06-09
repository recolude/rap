package euler

import (
	"bytes"
	"encoding/binary"

	"github.com/recolude/rap/format/collection/euler"
)

func encodeRaw32(captures []euler.Capture) []byte {
	streamData := new(bytes.Buffer)
	for _, capture := range captures {
		binary.Write(streamData, binary.LittleEndian, float32(capture.EulerZXY().X()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.EulerZXY().Y()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.EulerZXY().Z()))
	}
	return streamData.Bytes()
}

func decodeRaw32(streamData *bytes.Reader, times []float64) ([]euler.Capture, error) {
	var x float32
	var y float32
	var z float32

	captures := make([]euler.Capture, len(times))
	for i := 0; i < len(times); i++ {
		binary.Read(streamData, binary.LittleEndian, &x)
		binary.Read(streamData, binary.LittleEndian, &y)
		binary.Read(streamData, binary.LittleEndian, &z)
		captures[i] = euler.NewEulerZXYCapture(times[i], float64(x), float64(y), float64(z))
	}

	return captures, nil
}
