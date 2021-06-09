package euler_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/recolude/rap/format"
	eulerCollection "github.com/recolude/rap/format/collection/euler"
	"github.com/recolude/rap/format/encoding/euler"
	"github.com/stretchr/testify/assert"
)

func Test_Euler(t *testing.T) {
	continuousCaptures := make([]eulerCollection.Capture, 1000)
	continuousTimes := make([]float64, len(continuousCaptures))
	curTime := 1.0
	for i := 0; i < len(continuousCaptures); i++ {
		continuousTimes[i] = curTime
		continuousCaptures[i] = eulerCollection.NewEulerZXYCapture(
			curTime,
			(rand.Float64() * 360),
			(rand.Float64() * 360),
			(rand.Float64() * 360),
		)
		curTime += rand.Float64() * 10.0
	}

	tests := map[string]struct {
		captures []eulerCollection.Capture
		times    []float64
	}{
		"nil rotations": {captures: nil},
		"0-rotations":   {captures: []eulerCollection.Capture{}},
		"1-rotations": {
			captures: []eulerCollection.Capture{eulerCollection.NewEulerZXYCapture(1.2, 1, 1, 1)},
			times:    []float64{1.2},
		},
		"2-rotations": {
			captures: []eulerCollection.Capture{eulerCollection.NewEulerZXYCapture(1.2, 1, 1, 1), eulerCollection.NewEulerZXYCapture(1.3, 4, 5, 6)},
			times:    []float64{1.2, 1.3},
		},
		"3-rotations": {
			captures: []eulerCollection.Capture{
				eulerCollection.NewEulerZXYCapture(1.2, 1, 1, 1),
				eulerCollection.NewEulerZXYCapture(1.3, 4, 5, 6),
				eulerCollection.NewEulerZXYCapture(1.4, 4.1, 5.7, 6.0),
			},
			times: []float64{1.2, 1.3, 1.4},
		},
		"1000-rotations": {captures: continuousCaptures, times: continuousTimes},
	}

	storageTechniques := []struct {
		displayName        string
		technique          euler.StorageTechnique
		positionTollerance float64
	}{
		{
			displayName:        "Raw64",
			technique:          euler.Raw64,
			positionTollerance: 0,
		},
		{
			displayName:        "Raw32",
			technique:          euler.Raw32,
			positionTollerance: 0.0003,
		},
		{
			displayName:        "Raw16",
			technique:          euler.Raw16,
			positionTollerance: 0.04,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				streamIn := eulerCollection.NewCollection("Pos", tc.captures)

				encoder := euler.NewEncoder(technique.technique)

				// ACT ====================================================================
				header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
				streamOut, decodeErr := encoder.Decode(header, streamsData[0], tc.times)

				// ASSERT =================================================================
				assert.NoError(t, encodeErr)
				assert.NoError(t, decodeErr)
				assert.Len(t, header, 0)
				assert.Len(t, streamsData, 1)
				if assert.NotNil(t, streamOut) {
					assert.Equal(t, streamIn.Name(), streamOut.Name())
					if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {
						for i, c := range streamOut.Captures() {
							positioniCapture, ok := c.(eulerCollection.Capture)
							if assert.True(t, ok) == false {
								break
							}

							assert.Equal(t, tc.captures[i].Time(), positioniCapture.Time(), "times are not equal: %.2f != %.2f", tc.captures[i].Time(), positioniCapture.Time())
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

func Test_EulerRaw16Wrapping(t *testing.T) {
	// ARRANGE ================================================================
	testCaptures := []struct {
		in  eulerCollection.Capture
		out eulerCollection.Capture
	}{
		{
			in:  eulerCollection.NewEulerZXYCapture(0, 10+360, 20, 30),
			out: eulerCollection.NewEulerZXYCapture(0, 10, 20, 30),
		},
		{
			in:  eulerCollection.NewEulerZXYCapture(1, 10, 20, 30+360),
			out: eulerCollection.NewEulerZXYCapture(1, 10, 20, 30),
		},
		{
			in:  eulerCollection.NewEulerZXYCapture(2, 10, 20+360, 30),
			out: eulerCollection.NewEulerZXYCapture(2, 10, 20, 30),
		},
		{
			in:  eulerCollection.NewEulerZXYCapture(3, 10+360, 20+720, 30+1080),
			out: eulerCollection.NewEulerZXYCapture(3, 10, 20, 30),
		},
		{
			in:  eulerCollection.NewEulerZXYCapture(4, -10, -20, -30),
			out: eulerCollection.NewEulerZXYCapture(4, 350, 340, 330),
		},
		{
			in:  eulerCollection.NewEulerZXYCapture(5, -10-360, -720, -360),
			out: eulerCollection.NewEulerZXYCapture(5, 350, 0, 0),
		},
	}

	capturesIn := make([]eulerCollection.Capture, len(testCaptures))
	correctAnswers := make([]eulerCollection.Capture, len(testCaptures))
	for i, capture := range testCaptures {
		capturesIn[i] = capture.in
		correctAnswers[i] = capture.out
	}

	streamIn := eulerCollection.NewCollection("Rot", capturesIn)
	encoder := euler.NewEncoder(euler.Raw16)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
	streamOut, decodeErr := encoder.Decode(header, streamsData[0], []float64{0, 1, 2, 3, 4, 5})

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 1)
	assert.Len(t, streamOut.Captures(), len(streamIn.Captures()))

	for i, attempt := range streamOut.Captures() {

		rotationCapture, ok := attempt.(eulerCollection.Capture)
		if assert.True(t, ok) == false {
			break
		}

		assert.InDelta(t, correctAnswers[i].Time(), rotationCapture.Time(), 0.001, "times are not equal: %.2f != %.2f", correctAnswers[i].Time(), rotationCapture.Time())
		if assert.InDelta(
			t,
			correctAnswers[i].EulerZXY().X(),
			rotationCapture.EulerZXY().X(),
			0.005,
			"[%d] X components not equal: %.2f != %.2f", i, correctAnswers[i].EulerZXY().X(), rotationCapture.EulerZXY().X(),
		) == false {
			break
		}
		if assert.InDelta(
			t,
			correctAnswers[i].EulerZXY().Y(),
			rotationCapture.EulerZXY().Y(),
			0.005,
			"Y components not equal: %.2f != %.2f", correctAnswers[i].EulerZXY().Y(), rotationCapture.EulerZXY().Y(),
		) == false {
			break
		}
		if assert.InDelta(
			t,
			correctAnswers[i].EulerZXY().Z(),
			rotationCapture.EulerZXY().Z(),
			0.005,
			"Z components not equal: %.2f != %.2f", correctAnswers[i].EulerZXY().Z(), rotationCapture.EulerZXY().Z(),
		) == false {
			break
		}
	}
}
