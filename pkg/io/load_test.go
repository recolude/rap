package io

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var v1DirectoryTestData string = "../../test/data/io/v1"

func Test_Load_PanicsOnNilReader(t *testing.T) {
	assert.Panics(t, func() {
		Load(nil)
	})
}

func Test_Load_ErrorsOnEmptyBuffer(t *testing.T) {
	// ARRANGE ================================================================
	emptyBuffer := bytes.Buffer{}

	// ACT ====================================================================
	rec, bytesRead, err := Load(&emptyBuffer)

	// ASSERT =================================================================
	assert.Nil(t, rec)
	assert.Equal(t, 0, bytesRead)
	assert.EqualError(t, err, io.EOF.Error())
}

func Test_Load_ErrorsOnUnrecognizedFileVersion(t *testing.T) {
	// ARRANGE ================================================================
	buf := bytes.Buffer{}
	buf.Write([]byte{3})

	// ACT ====================================================================
	rec, bytesRead, err := Load(&buf)

	// ASSERT =================================================================
	assert.Nil(t, rec)
	assert.Equal(t, 1, bytesRead)
	assert.EqualError(t, err, "Unrecognized file version: 3")
}

func TestLoad(t *testing.T) {
	f, err := os.Open(filepath.Join(v1DirectoryTestData, "Demo 38subj v1.rap"))
	if assert.NoError(t, err) == false {
		return
	}

	allBytes, err := ioutil.ReadAll(f)
	if assert.NoError(t, err) == false {
		return
	}

	// ACT ====================================================================
	rec, bytesRead, err := Load(bytes.NewReader(allBytes))

	// ASSERT =================================================================
	if assert.NoError(t, err) == false {
		return
	}
	if assert.NotNil(t, rec) == false {
		return
	}

	assert.Equal(t, len(allBytes), bytesRead)

	assert.Equal(t, "Demo", rec.Name())
	assert.Len(t, rec.CaptureStreams(), 1)
	assert.Equal(t, "Custom Event", rec.CaptureStreams()[0].Name())
	assert.Len(t, rec.Recordings(), 38)

	for _, subj := range rec.Recordings() {
		assert.Len(t, subj.CaptureStreams(), 4)
		assert.Equal(t, "Position", subj.CaptureStreams()[0].Name())
		assert.Equal(t, "Rotation", subj.CaptureStreams()[1].Name())
		assert.Equal(t, "Custom Event", subj.CaptureStreams()[2].Name())
		assert.Equal(t, "Life Cycle", subj.CaptureStreams()[3].Name())
	}
}
