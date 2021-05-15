package main

import (
	"bytes"
	"testing"

	"github.com/recolude/rap/format/collection/position"
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

	if assert.Len(t, recording.Recordings()[0].CaptureCollections(), 1) == false {
		return
	}

	collection := recording.Recordings()[0].CaptureCollections()[0]
	if assert.NotNil(t, collection) {
		assert.Equal(t, "Position", collection.Name())
		assert.Equal(t, "recolude.position", collection.Signature())
		if assert.Len(t, collection.Captures(), 3) {
			assert.Equal(t, position.NewCapture(1, 2, 3, 4), collection.Captures()[0])
			assert.Equal(t, position.NewCapture(2, 5, 6, 7), collection.Captures()[1])
			assert.Equal(t, position.NewCapture(3, 8, 9, 10), collection.Captures()[2])
		}
	}
}
