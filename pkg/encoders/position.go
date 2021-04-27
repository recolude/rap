package encoders

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/recolude/rap/pkg/data"
)

type Position struct{}

func (p Position) encode(stream data.CaptureStream) ([]byte, error) {
	streamData := new(bytes.Buffer)

	// Write duration
	duration := streamDuration(stream)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(duration))
	_, err := streamData.Write(buf)
	if err != nil {
		return nil, err
	}

	return streamData.Bytes(), nil
}

func (p Position) Encode(streams []data.CaptureStream) ([]byte, [][]byte, error) {
	allStreamData := make([][]byte, len(streams))

	for i, stream := range streams {
		s, err := p.encode(stream)
		if err != nil {
			return nil, nil, err
		}
		allStreamData[i] = s
	}

	return nil, allStreamData, nil
}
