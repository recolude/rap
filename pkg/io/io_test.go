package io_test

import (
	"bytes"
	"testing"

	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/encoding"
	"github.com/recolude/rap/pkg/io"
	"github.com/recolude/rap/pkg/streams/position"
	"github.com/stretchr/testify/assert"
)

func Test_HandlesOneRecordingOneStream(t *testing.T) {
	// ARRANGE ================================================================
	fileData := new(bytes.Buffer)

	encoders := []encoding.Encoder{
		encoding.NewPositionEncoder(encoding.Raw64),
	}

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	recIn := data.NewRecording(
		"Test Recording",
		[]data.CaptureStream{
			position.NewStream(
				"Position",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
		},
		nil,
		map[string]string{
			"a":  "bee",
			"ce": "dee",
		},
		nil,
	)

	// ACT ====================================================================
	n, errWrite := w.Write(recIn)
	recOut, nOut, errRead := r.Read()

	// ASSERT =================================================================
	assert.NoError(t, errWrite)
	assert.NoError(t, errRead)
	assert.Equal(t, n, nOut)
	if assert.NotNil(t, recOut) {
		assert.Equal(t, recIn.Name(), recOut.Name())

		if assert.Equal(t, len(recIn.Metadata()), len(recOut.Metadata())) {
			for key, element := range recIn.Metadata() {
				assert.Equal(t, element, recOut.Metadata()[key])
			}
		}

		assert.Equal(t, len(recIn.CaptureStreams()), len(recOut.CaptureStreams()))
	}
}
