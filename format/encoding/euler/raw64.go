package euler

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/recolude/rap/format/collection/euler"
)

func encodeRaw64(captures []euler.Capture) []byte {
	streamData := new(bytes.Buffer)

	buf := make([]byte, binary.MaxVarintLen64)

	for _, capture := range captures {
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.EulerZXY().X()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.EulerZXY().Y()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.EulerZXY().Z()))
		streamData.Write(buf)
	}
	return streamData.Bytes()
}

func decodeRaw64(streamData *bytes.Reader, times []float64) ([]euler.Capture, error) {
	// streamData := bytes.NewReader(captureData)

	captures := make([]euler.Capture, len(times))
	buf := make([]byte, binary.MaxVarintLen64)
	for i := 0; i < len(times); i++ {
		_, err := streamData.Read(buf)
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

		captures[i] = euler.NewEulerZXYCapture(times[i], x, y, z)
	}

	return captures, nil
}
