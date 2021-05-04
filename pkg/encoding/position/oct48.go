package position

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/recolude/rap/pkg/streams/position"
)

func encodeOct48(captures []position.Capture) ([]byte, error) {
	streamData := new(bytes.Buffer)

	err := binary.Write(streamData, binary.LittleEndian, captures[0].Time())
	if err != nil {
		return nil, err
	}

	startingTime := math.Inf(1)
	endingTime := math.Inf(-1)

	for _, capture := range captures {
		if capture.Time() < startingTime {
			startingTime = capture.Time()
		}
		if capture.Time() > endingTime {
			endingTime = capture.Time()
		}
	}

	duration := endingTime - startingTime

	err = binary.Write(streamData, binary.LittleEndian, duration)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 8)
	size := binary.PutUvarint(buf, uint64(len(captures)))
	streamData.Write(buf[:size])

	for _, capture := range captures {
		binary.Write(streamData, binary.LittleEndian, float32(capture.Time()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.Position().X()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.Position().Y()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.Position().Z()))
	}

	return streamData.Bytes(), nil
}
