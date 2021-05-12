package euler

import (
	"bytes"
	"encoding/binary"

	"github.com/recolude/rap/format/streams/euler"
)

func encodeRaw32(captures []euler.Capture) []byte {
	streamData := new(bytes.Buffer)

	buf := make([]byte, 8)
	size := binary.PutUvarint(buf, uint64(len(captures)))
	streamData.Write(buf[:size])

	for _, capture := range captures {
		binary.Write(streamData, binary.LittleEndian, float32(capture.Time()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.EulerZXY().X()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.EulerZXY().Y()))
		binary.Write(streamData, binary.LittleEndian, float32(capture.EulerZXY().Z()))
	}
	return streamData.Bytes()
}

func decodeRaw32(streamData *bytes.Reader) ([]euler.Capture, error) {
	// streamData := bytes.NewReader(captureData)

	numCaptures, err := binary.ReadUvarint(streamData)
	if err != nil {
		return nil, err
	}

	captures := make([]euler.Capture, numCaptures)
	for i := 0; i < int(numCaptures); i++ {
		var time float32
		var x float32
		var y float32
		var z float32

		binary.Read(streamData, binary.LittleEndian, &time)
		binary.Read(streamData, binary.LittleEndian, &x)
		binary.Read(streamData, binary.LittleEndian, &y)
		binary.Read(streamData, binary.LittleEndian, &z)
		captures[i] = euler.NewEulerZXYCapture(float64(time), float64(x), float64(y), float64(z))
	}

	return captures, nil
}
