package io_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	enumEncoding "github.com/recolude/rap/format/encoding/enum"
	eulerEncoding "github.com/recolude/rap/format/encoding/euler"
	eventEncoding "github.com/recolude/rap/format/encoding/event"
	positionEncoding "github.com/recolude/rap/format/encoding/position"
	"github.com/recolude/rap/format/io"
	"github.com/recolude/rap/format/streams/enum"
	"github.com/recolude/rap/format/streams/euler"
	"github.com/recolude/rap/format/streams/event"
	"github.com/recolude/rap/format/streams/position"
	"github.com/stretchr/testify/assert"
)

func assertRecordingsMatch(t *testing.T, recExpected, recActual format.Recording) bool {
	if recExpected != nil && recActual == nil {
		t.Error("Expected recording to not be nil, but it was")
		return false
	}

	if assert.Equal(t, len(recExpected.Binaries()), len(recActual.Binaries())) == false {
		return false
	}

	if assert.Equal(t, len(recExpected.Recordings()), len(recActual.Recordings()), "Mismatch child recordings") == false {
		return false
	}

	for i, childRec := range recActual.Recordings() {
		if assertRecordingsMatch(t, recExpected.Recordings()[i], childRec) == false {
			return false
		}
	}

	if assert.NotNil(t, recActual) == false {
		return false
	}

	if assert.Equal(t, recExpected.Name(), recActual.Name()) == false {
		return false
	}

	if assert.Equal(t, len(recExpected.Metadata()), len(recActual.Metadata())) == false {
		return false
	}

	for key, element := range recExpected.Metadata() {
		assert.Equal(t, element, recActual.Metadata()[key])
	}

	if assert.Equal(t, len(recExpected.CaptureStreams()), len(recActual.CaptureStreams())) == false {
		return false
	}

	for streamIndex, stream := range recExpected.CaptureStreams() {
		assert.Equal(t, stream.Name(), recActual.CaptureStreams()[streamIndex].Name())

		for i, correctCapture := range recExpected.CaptureStreams()[streamIndex].Captures() {
			assert.Equal(t, correctCapture.Time(), recActual.CaptureStreams()[streamIndex].Captures()[i].Time())
		}

	}

	return true
}

func Test_HandlesOneRecordingOneStream(t *testing.T) {
	// ARRANGE ================================================================
	fileData := new(bytes.Buffer)

	encoders := []encoding.Encoder{
		positionEncoding.NewEncoder(positionEncoding.Raw64),
	}

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	recIn := format.NewRecording(
		"",
		"Test Recording",
		[]format.CaptureStream{
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
	assertRecordingsMatch(t, recIn, recOut)
}

func Test_HandlesOneRecordingTwoStream(t *testing.T) {
	// ARRANGE ================================================================
	fileData := new(bytes.Buffer)

	encoders := []encoding.Encoder{
		positionEncoding.NewEncoder(positionEncoding.Raw64),
	}

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	recIn := format.NewRecording(
		"",
		"Test Recording",
		[]format.CaptureStream{
			position.NewStream(
				"Position",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
			position.NewStream(
				"Position2",
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
	assertRecordingsMatch(t, recIn, recOut)
}

func Test_HandlesNestedRecordings(t *testing.T) {
	// ARRANGE ================================================================
	fileData := new(bytes.Buffer)

	encoders := []encoding.Encoder{
		positionEncoding.NewEncoder(positionEncoding.Raw64),
	}

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	recIn := format.NewRecording(
		"",
		"Test Recording",
		[]format.CaptureStream{
			position.NewStream(
				"Position",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
			position.NewStream(
				"Position2",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
		},
		[]format.Recording{
			format.NewRecording(
				"",
				"Child Recording",
				[]format.CaptureStream{
					position.NewStream(
						"Child Position",
						[]position.Capture{
							position.NewCapture(1, 1, 2, 3),
							position.NewCapture(2, 4, 5, 6),
							position.NewCapture(4, 7, 8, 9),
							position.NewCapture(7, 10, 11, 12),
						},
					),
					position.NewStream(
						"Child Position2",
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
			),
		},
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
	assertRecordingsMatch(t, recIn, recOut)
}

func Test_EncodersWithHeaders(t *testing.T) {
	// ARRANGE ================================================================
	fileData := new(bytes.Buffer)

	encoders := []encoding.Encoder{
		positionEncoding.NewEncoder(positionEncoding.Raw64),
		eulerEncoding.NewEncoder(eulerEncoding.Raw64),
		eventEncoding.NewEncoder(eventEncoding.Raw64),
		enumEncoding.NewEncoder(enumEncoding.Raw64),
	}

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	recIn := format.NewRecording(
		"",
		"Test Recording",
		[]format.CaptureStream{
			position.NewStream(
				"Position",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
				},
			),
		},
		[]format.Recording{
			format.NewRecording(
				"",
				"Child",
				[]format.CaptureStream{
					event.NewStream("ahhh", []event.Capture{
						event.NewCapture(1, "att", map[string]string{"1": "2"}),
					}),
				},
				nil,
				nil,
				nil,
			),
		},
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
	assertRecordingsMatch(t, recIn, recOut)
}

func Test_HandlesMultipleEncoders(t *testing.T) {
	// ARRANGE ================================================================
	fileData := new(bytes.Buffer)

	encoders := []encoding.Encoder{
		positionEncoding.NewEncoder(positionEncoding.Raw64),
		eulerEncoding.NewEncoder(eulerEncoding.Raw64),
		eventEncoding.NewEncoder(eventEncoding.Raw64),
		enumEncoding.NewEncoder(enumEncoding.Raw64),
	}

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	recIn := format.NewRecording(
		"",
		"Test Recording",
		[]format.CaptureStream{
			position.NewStream(
				"Position",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
			position.NewStream(
				"Position2",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
		},
		[]format.Recording{
			format.NewRecording(
				"",
				"Child Recording",
				[]format.CaptureStream{
					event.NewStream("ahhh", []event.Capture{
						event.NewCapture(1, "att", map[string]string{"1": "2"}),
					}),
					position.NewStream(
						"Child Position",
						[]position.Capture{
							position.NewCapture(1, 1, 2, 3),
							position.NewCapture(2, 4, 5, 6),
							position.NewCapture(4, 7, 8, 9),
							position.NewCapture(7, 10, 11, 12),
						},
					),
					position.NewStream(
						"Child Position2",
						[]position.Capture{
							position.NewCapture(1, 1, 2, 3),
							position.NewCapture(2, 4, 5, 6),
							position.NewCapture(4, 7, 8, 9),
							position.NewCapture(7, 10, 11, 12),
						},
					),
					euler.NewStream(
						"Rot",
						[]euler.Capture{
							euler.NewEulerZXYCapture(1, 1, 2, 3),
							euler.NewEulerZXYCapture(2, 4, 5, 6),
							euler.NewEulerZXYCapture(4, 7, 8, 9),
							euler.NewEulerZXYCapture(7, 10, 11, 12),
						},
					),
					enum.NewStream(
						"cmon",
						[]string{"A", "n"},
						[]enum.Capture{
							enum.NewCapture(1, 1),
						},
					),
				},
				nil,
				map[string]string{
					"a":  "bee",
					"ce": "dee",
				},
				nil,
			),
			format.NewRecording(
				"",
				"Child 2 Recording",
				[]format.CaptureStream{
					position.NewStream(
						"Child Position",
						[]position.Capture{
							position.NewCapture(1, 1, 2, 3),
							position.NewCapture(2, 4, 5, 6),
							position.NewCapture(4, 7, 8, 9),
							position.NewCapture(7, 10, 11, 12),
						},
					),
					position.NewStream(
						"Child Position2",
						[]position.Capture{
							position.NewCapture(1, 1, 2, 3),
							position.NewCapture(2, 4, 5, 6),
							position.NewCapture(4, 7, 8, 9),
							position.NewCapture(7, 10, 11, 12),
						},
					),
					euler.NewStream(
						"Rot",
						[]euler.Capture{
							euler.NewEulerZXYCapture(1, 1, 2, 3),
							euler.NewEulerZXYCapture(2, 4, 5, 6),
							euler.NewEulerZXYCapture(4, 7, 8, 9),
							euler.NewEulerZXYCapture(7, 10, 11, 12),
						},
					),
				},
				nil,
				map[string]string{
					"a":  "bee",
					"ce": "dee",
				},
				nil,
			),
		},
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
	assertRecordingsMatch(t, recIn, recOut)
}

func Test_Uprade(t *testing.T) {
	f, err := os.Open(filepath.Join(v1DirectoryTestData, "Demo 38subj v1.rap"))
	if assert.NoError(t, err) == false {
		return
	}

	allBytes, err := ioutil.ReadAll(f)
	if assert.NoError(t, err) == false {
		return
	}

	encoders := []encoding.Encoder{
		positionEncoding.NewEncoder(positionEncoding.Raw64),
		eulerEncoding.NewEncoder(eulerEncoding.Raw64),
		eventEncoding.NewEncoder(eventEncoding.Raw64),
		enumEncoding.NewEncoder(enumEncoding.Raw64),
	}
	fileData := new(bytes.Buffer)

	w := io.NewWriter(encoders, fileData)
	r := io.NewReader(encoders, fileData)

	// ACT ====================================================================
	rec, _, err := io.Load(bytes.NewReader(allBytes))
	n, errWrite := w.Write(rec)
	recOut, nOut, errRead := r.Read()

	// ASSERT =================================================================
	assert.NoError(t, errWrite)
	assert.NoError(t, errRead)
	assert.Equal(t, n, nOut)
	assertRecordingsMatch(t, rec, recOut)
}
