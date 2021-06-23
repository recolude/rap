package format_test

import (
	"testing"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	rec := format.NewRecording("some-id", "dum name", []format.CaptureCollection{
		position.NewCollection("t", []position.Capture{
			position.NewCapture(1, 2, 3, 4),
			position.NewCapture(3, 4, 5, 6),
			position.NewCapture(4, 5, 6, 7),
			position.NewCapture(10, 11, 12, 13),
			position.NewCapture(11, 12, 13, 14),
		}),
	}, nil, metadata.EmptyBlock(), nil, nil)

	recSliced := format.Slice(
		rec,
		format.BeginningOfSlice(3),
		format.EndOfSlice(10),
		format.KeepBinariesInSlice(false),
	)

	if assert.NotNil(t, recSliced) == false {
		return
	}

	assert.Len(t, recSliced.CaptureCollections(), 1)
	assert.Len(t, recSliced.CaptureCollections()[0].Captures(), 2)
	assert.Equal(t, 3.0, recSliced.CaptureCollections()[0].Captures()[0].Time())
	assert.Equal(t, 4.0, recSliced.CaptureCollections()[0].Captures()[1].Time())
}
