package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_Summarize(t *testing.T) {
	// ARRANGE ================================================================
	recIn := format.NewRecording(
		"",
		"Test Recording",
		[]format.CaptureCollection{
			position.NewCollection(
				"Position",
				[]position.Capture{
					position.NewCapture(1, 1, 2, 3),
					position.NewCapture(2, 4, 5, 6),
					position.NewCapture(4, 7, 8, 9),
					position.NewCapture(7, 10, 11, 12),
				},
			),
			position.NewCollection(
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
				[]format.CaptureCollection{
					position.NewCollection(
						"Child Position",
						[]position.Capture{
							position.NewCapture(1, 1, 2, 3),
							position.NewCapture(2, 4, 5, 6),
							position.NewCapture(4, 7, 8, 9),
							position.NewCapture(7, 10, 11, 12),
						},
					),
					position.NewCollection(
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
				metadata.NewBlock(
					map[string]metadata.Property{
						"a":  metadata.NewStringProperty("bee"),
						"ce": metadata.NewStringProperty("dee"),
					},
				),
				nil,
				nil,
			),
		},
		metadata.NewBlock(
			map[string]metadata.Property{
				"a":  metadata.NewStringProperty("bee"),
				"ce": metadata.NewStringProperty("dee"),
			},
		),
		nil,
		nil,
	)

	answerBuilder := strings.Builder{}
	answerBuilder.WriteString("Name:                    Test Recording\n")
	answerBuilder.WriteString("Sub Recordings:          1\n")
	answerBuilder.WriteString("Total Position Captures: 16\n")
	answerBuilder.WriteString("Total Euler Captures:    0\n")
	answerBuilder.WriteString("Total Event Captures:    0\n")
	answerBuilder.WriteString("Total Enum Captures:     0\n")
	answerBuilder.WriteString("Total Other Captures:    0\n")

	out := bytes.Buffer{}

	// ACT ====================================================================
	printSummary(&out, recIn)

	// ASSERT =================================================================
	assert.Equal(t, answerBuilder.String(), out.String())
}
