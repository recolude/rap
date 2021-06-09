package float_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/recolude/rap/format"
	floatCollection "github.com/recolude/rap/format/collection/float"
	"github.com/recolude/rap/format/encoding/float"
	"github.com/stretchr/testify/assert"
)

func Test_Float(t *testing.T) {
	continuousCaptures := make([]floatCollection.Capture, 2000)
	curTime := -1000.0
	curPos := -1000.0
	for i := 0; i < len(continuousCaptures); i++ {
		continuousCaptures[i] = floatCollection.NewCapture(
			curTime,
			curPos,
		)
		curPos += rand.Float64() * 10.0
		curTime += rand.Float64() * 10.0
	}

	tests := map[string]struct {
		streamName string
		captures   []floatCollection.Capture
	}{
		"nil floats": {streamName: "", captures: nil},
		"0-float":    {streamName: "empty stream", captures: []floatCollection.Capture{}},
		"1-float": {
			streamName: "ahhhh",
			captures: []floatCollection.Capture{
				floatCollection.NewCapture(1.2, 1),
			},
		},
		"2-floats": {
			streamName: "222",
			captures: []floatCollection.Capture{
				floatCollection.NewCapture(-1.2, 0),
				floatCollection.NewCapture(0, 1.3),
			},
		},
		"3-floats": {
			streamName: "ahhhh",
			captures: []floatCollection.Capture{
				floatCollection.NewCapture(-1.2, 1),
				floatCollection.NewCapture(1.2, -1),
				floatCollection.NewCapture(1.3, 1.3),
			},
		},
		// "2000-continuous-floats": {captures: continuousCaptures},
	}

	storageTechniques := []struct {
		displayName     string
		technique       float.StorageTechnique
		timeTollerance  float64
		valueTollerance float64
	}{
		{
			displayName:    "Raw64",
			technique:      float.Raw64,
			timeTollerance: 0,
		},
		{
			displayName:    "Raw32",
			technique:      float.Raw32,
			timeTollerance: 0.0005,
		},
		// {
		// 	displayName:    "BST16",
		// 	technique:      float.BST16,
		// 	timeTollerance: 0.0005,
		// },
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				collectionIn := floatCollection.NewCollection(tc.streamName, tc.captures)

				encoder := float.NewEncoder(technique.technique)

				// ACT ====================================================================
				header, collectionData, encodeErr := encoder.Encode([]format.CaptureCollection{collectionIn})
				collectionOut, decodeErr := encoder.Decode(header, collectionData[0], nil)

				// ASSERT =================================================================
				assert.NoError(t, encodeErr)
				assert.NoError(t, decodeErr)
				assert.NotNil(t, collectionOut)
				assert.Len(t, collectionData, 1)
				assert.Equal(t, tc.streamName, collectionOut.Name())
				if assert.NotNil(t, collectionOut) {
					assert.Equal(t, collectionIn.Name(), collectionOut.Name())
					if assert.Len(t, collectionOut.Captures(), len(collectionIn.Captures())) {
						for i, c := range collectionOut.Captures() {
							capture, ok := c.(floatCollection.Capture)
							if assert.True(t, ok) == false {
								break
							}

							assert.InDelta(t, tc.captures[i].Time(), capture.Time(), technique.timeTollerance, "times are not equal: %.2f != %.2f", tc.captures[i].Time(), capture.Time())
							assert.InDelta(t, tc.captures[i].Value(), capture.Value(), technique.timeTollerance, "values are not equal: %.2f != %.2f", tc.captures[i].Value(), capture.Value())

						}
					}
				}
			})
		}
	}
}
