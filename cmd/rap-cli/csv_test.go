package main

import (
	"bytes"
	"testing"

	"github.com/recolude/rap/format/streams/position"
	"github.com/stretchr/testify/assert"
)

func Test_CSV_Simple(t *testing.T) {
	// ARRANGE ================================================================
	csv := `id, name, time, x, y, z
0, bob, 1, 2, 3, 4
0, bob, 2, 5, 6, 7
0, bob, 3, 8, 9, 10
`

	// ACT ====================================================================
	recording, err := RecordingFromCSV(bytes.NewReader([]byte(csv)))

	// ASSERT =================================================================
	assert.NoError(t, err)
	assert.NotNil(t, recording)
	if assert.Len(t, recording.Recordings(), 1) == false {
		return
	}

	if assert.Len(t, recording.Recordings()[0].CaptureStreams(), 1) == false {
		return
	}

	stream := recording.Recordings()[0].CaptureStreams()[0]
	if assert.NotNil(t, stream) {
		assert.Equal(t, "Position", stream.Name())
		assert.Equal(t, "recolude.position", stream.Signature())
		if assert.Len(t, stream.Captures(), 3) {
			assert.Equal(t, position.NewCapture(1, 2, 3, 4), stream.Captures()[0])
			assert.Equal(t, position.NewCapture(2, 5, 6, 7), stream.Captures()[1])
			assert.Equal(t, position.NewCapture(3, 8, 9, 10), stream.Captures()[2])
		}
	}
}
