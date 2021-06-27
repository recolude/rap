package event_test

import (
	"testing"

	"github.com/recolude/rap/format"
	eventStream "github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/encoding/event"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_SingleEvent(t *testing.T) {
	tests := map[string]struct {
		streamName string
		captures   []eventStream.Capture
		times      []float64
	}{
		"nil events": {streamName: "", captures: nil},
		"0-events":   {streamName: "empty stream", captures: []eventStream.Capture{}},
		"1-events": {
			streamName: "ahhhh",
			captures: []eventStream.Capture{
				eventStream.NewCapture(
					1.2,
					"Damage",
					metadata.NewBlock(map[string]metadata.Property{
						"dealer": metadata.NewStringProperty("player"),
						"damage": metadata.NewIntProperty(7),
					}),
				),
			},
			times: []float64{1.2},
		},
		"2-events": {
			streamName: "ahhhh",
			captures: []eventStream.Capture{
				eventStream.NewCapture(
					1.2,
					"",
					metadata.NewBlock(map[string]metadata.Property{
						"dealer": metadata.NewStringProperty("arer"),
						"damage": metadata.NewIntProperty(7),
					}),
				),
				eventStream.NewCapture(
					1.6,
					"ttt",
					metadata.NewBlock(map[string]metadata.Property{
						"dealer": metadata.NewStringProperty("watcher"),
						"damage": metadata.NewIntProperty(40),
					}),
				),
			},
			times: []float64{1.2, 1.6},
		},
	}

	encoder := event.NewEncoder()
	assert.Equal(t, "recolude.event", encoder.Signature())
	assert.Equal(t, uint(0), encoder.Version())

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			streamIn := eventStream.NewCollection(tc.streamName, tc.captures)
			assert.True(t, encoder.Accepts(streamIn))

			// ACT ====================================================================
			header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
			streamOut, decodeErr := encoder.Decode(tc.streamName, header, streamsData[0], tc.times)

			// ASSERT =================================================================
			assert.NoError(t, encodeErr)
			assert.NoError(t, decodeErr)
			assert.NotNil(t, streamOut)
			assert.Len(t, streamsData, 1)
			assert.Equal(t, tc.streamName, streamOut.Name())
			if assert.NotNil(t, streamOut) {
				assert.Equal(t, streamIn.Name(), streamOut.Name())
				if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {
					for i, c := range streamOut.Captures() {
						eventCapture, ok := c.(eventStream.Capture)
						if assert.True(t, ok) == false {
							break
						}

						assert.Equal(t, tc.captures[i].Time(), eventCapture.Time(), "times are not equal: %.2f != %.2f", tc.captures[i].Time(), eventCapture.Time())
						assert.Equal(t, tc.captures[i].Name(), eventCapture.Name())

						assert.Len(t, eventCapture.Metadata().Mapping(), len(tc.captures[i].Metadata().Mapping()))
						for key, val := range tc.captures[i].Metadata().Mapping() {
							assert.Equal(t, val, eventCapture.Metadata().Mapping()[key])
						}
					}
				}
			}
		})
	}
}
