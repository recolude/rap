package event_test

import (
	"fmt"
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
		},
	}

	storageTechniques := []struct {
		displayName    string
		technique      event.StorageTechnique
		timeTollerance float64
	}{
		{
			displayName:    "Raw64",
			technique:      event.Raw64,
			timeTollerance: 0,
		},
		{
			displayName:    "Raw32",
			technique:      event.Raw32,
			timeTollerance: 0.0005,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				streamIn := eventStream.NewCollection(tc.streamName, tc.captures)

				encoder := event.NewEncoder(technique.technique)

				// ACT ====================================================================
				header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
				streamOut, decodeErr := encoder.Decode(header, streamsData[0])

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

							assert.InDelta(t, tc.captures[i].Time(), eventCapture.Time(), technique.timeTollerance, "times are not equal: %.2f != %.2f", tc.captures[i].Time(), eventCapture.Time())
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
}
