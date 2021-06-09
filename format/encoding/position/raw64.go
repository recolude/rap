package position

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/recolude/rap/format/collection/position"
)

func encodeRaw64(captures []position.Capture) []byte {
	streamData := new(bytes.Buffer)

	buf := make([]byte, binary.MaxVarintLen64)
	size := binary.PutUvarint(buf, uint64(len(captures)))
	streamData.Write(buf[:size])

	for _, capture := range captures {
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Position().X()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Position().Y()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Position().Z()))
		streamData.Write(buf)
	}
	return streamData.Bytes()
}

func decodeRaw64(streamData *bytes.Reader, times []float64) ([]position.Capture, error) {
	numCaptures, err := binary.ReadUvarint(streamData)
	if err != nil {
		return nil, err
	}

	captures := make([]position.Capture, numCaptures)
	buf := make([]byte, binary.MaxVarintLen64)
	for i := 0; i < int(numCaptures); i++ {
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

		captures[i] = position.NewCapture(times[i], x, y, z)
	}

	return captures, nil
}
