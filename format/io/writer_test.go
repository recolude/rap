package io_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/io"
	"github.com/recolude/rap/format/metadata"
	"github.com/recolude/rap/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_ErrorsWithNoValidEncoders(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := mocks.NewMockCaptureCollection(ctrl)
	stream.EXPECT().Signature().Return("test.data")

	rec := mocks.NewMockRecording(ctrl)
	rec.EXPECT().CaptureCollections().AnyTimes().Return([]format.CaptureCollection{stream})
	rec.EXPECT().Recordings().AnyTimes().Return([]format.Recording{})
	rec.EXPECT().BinaryReferences().AnyTimes().Return([]format.BinaryReference{})
	rec.EXPECT().Binaries().AnyTimes().Return([]format.Binary{})

	w := io.NewWriter(nil, false, nil)

	// ACT ====================================================================
	_, err := w.Write(rec)

	// ASSERT =================================================================
	assert.EqualError(t, err, "no encoder registered to handle stream: test.data")
}

func Test_PanicsWithNilRecording(t *testing.T) {
	// ARRANGE ================================================================
	w := io.NewWriter(nil, false, nil)

	// ACT/ASSERT =============================================================
	assert.PanicsWithError(t, "can not write nil recording", func() { w.Write(nil) })
}

func Test_ErrorsWithNilSubRecording(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rec := mocks.NewMockRecording(ctrl)
	rec.EXPECT().Recordings().AnyTimes().Return([]format.Recording{nil})
	rec.EXPECT().BinaryReferences().AnyTimes().Return([]format.BinaryReference{})
	rec.EXPECT().Binaries().AnyTimes().Return([]format.Binary{})

	out := bytes.Buffer{}

	w := io.NewRecoludeWriter(&out)

	// ACT ====================================================================
	written, err := w.Write(rec)

	// ASSERT =================================================================
	assert.Zero(t, written)
	assert.Len(t, out.Bytes(), 0)
	assert.EqualError(t, err, "can not serialize recording with nil sub-recordings")
}

func Test_ErrorsWithNilCaptureCollections(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rec := mocks.NewMockRecording(ctrl)
	rec.EXPECT().Recordings().AnyTimes().Return([]format.Recording{})
	rec.EXPECT().CaptureCollections().AnyTimes().Return([]format.CaptureCollection{nil})
	rec.EXPECT().BinaryReferences().AnyTimes().Return([]format.BinaryReference{})
	rec.EXPECT().Binaries().AnyTimes().Return([]format.Binary{})

	out := bytes.Buffer{}

	w := io.NewRecoludeWriter(&out)

	// ACT ====================================================================
	written, err := w.Write(rec)

	// ASSERT =================================================================
	assert.Zero(t, written)
	assert.Len(t, out.Bytes(), 0)
	assert.EqualError(t, err, "can not serialize recording with nil capture collections")
}

func Test_ErrorsWithNilBinaryReferences(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rec := mocks.NewMockRecording(ctrl)
	rec.EXPECT().Recordings().AnyTimes().Return([]format.Recording{})
	rec.EXPECT().CaptureCollections().AnyTimes().Return([]format.CaptureCollection{})
	rec.EXPECT().BinaryReferences().AnyTimes().Return([]format.BinaryReference{nil})
	rec.EXPECT().Metadata().AnyTimes().Return(metadata.EmptyBlock())
	rec.EXPECT().Binaries().AnyTimes().Return([]format.Binary{})

	out := bytes.Buffer{}

	w := io.NewRecoludeWriter(&out)

	// ACT ====================================================================
	written, err := w.Write(rec)

	// ASSERT =================================================================
	assert.Zero(t, written)
	assert.Len(t, out.Bytes(), 0)
	assert.EqualError(t, err, "can not serialize recording with nil binary references")
}

func Test_ErrorsWithNilBinaries(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rec := mocks.NewMockRecording(ctrl)
	rec.EXPECT().Recordings().AnyTimes().Return([]format.Recording{})
	rec.EXPECT().CaptureCollections().AnyTimes().Return([]format.CaptureCollection{})
	rec.EXPECT().BinaryReferences().AnyTimes().Return([]format.BinaryReference{})
	rec.EXPECT().Metadata().AnyTimes().Return(metadata.EmptyBlock())
	rec.EXPECT().Binaries().AnyTimes().Return([]format.Binary{nil})

	out := bytes.Buffer{}

	w := io.NewRecoludeWriter(&out)

	// ACT ====================================================================
	written, err := w.Write(rec)

	// ASSERT =================================================================
	assert.Zero(t, written)
	assert.Len(t, out.Bytes(), 0)
	assert.EqualError(t, err, "can not serialize recording with nil binaries")
}
