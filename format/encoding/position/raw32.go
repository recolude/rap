package position

import (
	"bytes"
	"encoding/binary"

	"github.com/recolude/rap/format/collection/position"
)

func encodeRaw32(captures []position.Capture) []byte {
	streamData := new(bytes.Buffer)

	for _, capture := range captures {
		binary.Write(streamData, binary.LittleEndian, float32(capture.Position().X()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.Position().Y()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.Position().Z()))
	}

	return streamData.Bytes()
}

func decodeRaw32(streamData *bytes.Reader, times []float64) ([]position.Capture, error) {
	captures := make([]position.Capture, len(times))
	for i := 0; i < len(times); i++ {
		var x float32
		var y float32
		var z float32

		binary.Read(streamData, binary.LittleEndian, &x)
		binary.Read(streamData, binary.LittleEndian, &y)
		binary.Read(streamData, binary.LittleEndian, &z)
		captures[i] = position.NewCapture(times[i], float64(x), float64(y), float64(z))
	}

	return captures, nil
}
