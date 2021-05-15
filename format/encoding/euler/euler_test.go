package euler_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/recolude/rap/format"
	eulerStream "github.com/recolude/rap/format/collection/euler"
	"github.com/recolude/rap/format/encoding/euler"
	"github.com/stretchr/testify/assert"
)

func Test_Euler(t *testing.T) {
	continuousCaptures := make([]eulerStream.Capture, 1000)
	curTime := 1.0
	for i := 0; i < len(continuousCaptures); i++ {
		continuousCaptures[i] = eulerStream.NewEulerZXYCapture(
			curTime,
			(-rand.Float64()*360)+360,
			(-rand.Float64()*360)+360,
			(-rand.Float64()*360)+360,
		)
		curTime += rand.Float64() * 10.0
	}

	tests := map[string]struct {
		captures []eulerStream.Capture
	}{
		"nil rotations": {captures: nil},
		"0-rotations":   {captures: []eulerStream.Capture{}},
		"1-rotations":   {captures: []eulerStream.Capture{eulerStream.NewEulerZXYCapture(1.2, 1, 1, 1)}},
		"2-rotations":   {captures: []eulerStream.Capture{eulerStream.NewEulerZXYCapture(1.2, 1, 1, 1), eulerStream.NewEulerZXYCapture(1.3, 4, 5, 6)}},
		"3-rotations": {
			captures: []eulerStream.Capture{
				eulerStream.NewEulerZXYCapture(1.2, 1, 1, 1),
				eulerStream.NewEulerZXYCapture(1.3, 4, 5, 6),
				eulerStream.NewEulerZXYCapture(1.4, 4.1, 5.7, 6.0),
			},
		},
		"1000-rotations": {captures: continuousCaptures},
	}

	storageTechniques := []struct {
		displayName        string
		technique          euler.StorageTechnique
		timeTollerance     float64
		positionTollerance float64
	}{
		{
			displayName:        "Raw64",
			technique:          euler.Raw64,
			timeTollerance:     0,
			positionTollerance: 0,
		},
		{
			displayName:        "Raw32",
			technique:          euler.Raw32,
			timeTollerance:     0.0005,
			positionTollerance: 0.0003,
		},
		{
			displayName:        "Raw16",
			technique:          euler.Raw16,
			timeTollerance:     0.01,
			positionTollerance: 0.04,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				streamIn := eulerStream.NewStream("Pos", tc.captures)

				encoder := euler.NewEncoder(technique.technique)

				// ACT ====================================================================
				header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
				streamOut, decodeErr := encoder.Decode(header, streamsData[0])

				// ASSERT =================================================================
				assert.NoError(t, encodeErr)
				assert.NoError(t, decodeErr)
				assert.Len(t, header, 0)
				assert.Len(t, streamsData, 1)
				if assert.NotNil(t, streamOut) {
					assert.Equal(t, streamIn.Name(), streamOut.Name())
					if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {
						for i, c := range streamOut.Captures() {
							positioniCapture, ok := c.(eulerStream.Capture)
							if assert.True(t, ok) == false {
								break
							}

							assert.InDelta(t, tc.captures[i].Time(), positioniCapture.Time(), technique.timeTollerance, "times are not equal: %.2f != %.2f", tc.captures[i].Time(), positioniCapture.Time())
							if assert.InDelta(
								t,
								tc.captures[i].EulerZXY().X(),
								positioniCapture.EulerZXY().X(),
								technique.positionTollerance,
								"[%d] X components not equal: %.2f != %.2f", i, tc.captures[i].EulerZXY().X(), positioniCapture.EulerZXY().X(),
							) == false {
								break
							}
							if assert.InDelta(
								t,
								tc.captures[i].EulerZXY().Y(),
								positioniCapture.EulerZXY().Y(),
								technique.positionTollerance,
								"Y components not equal: %.2f != %.2f", tc.captures[i].EulerZXY().Y(), positioniCapture.EulerZXY().Y(),
							) == false {
								break
							}
							if assert.InDelta(
								t,
								tc.captures[i].EulerZXY().Z(),
								positioniCapture.EulerZXY().Z(),
								technique.positionTollerance,
								"Z components not equal: %.2f != %.2f", tc.captures[i].EulerZXY().Z(), positioniCapture.EulerZXY().Z(),
							) == false {
								break
							}
						}
					}
				}
			})
		}
	}
}
