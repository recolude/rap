package format_test

import (
	"testing"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	rec := format.NewRecording(
		"some-id",
		"dum name",
		[]format.CaptureCollection{
			position.NewCollection("t", []position.Capture{
				position.NewCapture(1, 2, 3, 4),
				position.NewCapture(3, 4, 5, 6),
				position.NewCapture(4, 5, 6, 7),
				position.NewCapture(10, 11, 12, 13),
				position.NewCapture(11, 12, 13, 14),
			}),
		},
		[]format.Recording{
			format.NewRecording("123", "child", []format.CaptureCollection{
				position.NewCollection("empty", nil),
				position.NewCollection("Position", []position.Capture{
					position.NewCapture(1, 2, 3, 4),
					position.NewCapture(4, 5, 6, 7),
					position.NewCapture(3, 4, 5, 6),
					position.NewCapture(10, 11, 12, 13),
					position.NewCapture(11, 12, 13, 14),
				}),
			}, nil, metadata.EmptyBlock(), nil, nil),
		},
		metadata.EmptyBlock(),
		nil,
		nil,
	)

	err := format.Validate(
		rec,
		format.RequireChronologicalCapture(true),
	)

	assert.EqualError(t, err, "[123] child: Position capture collection violates chronological event validator")
}
