package position_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/EliCDavis/vector/vector3"
	"github.com/recolude/rap/format"
	positionCollection "github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/encoding/position"
	"github.com/stretchr/testify/assert"
)

func Test_Oct24_EmptyStream(t *testing.T) {
	// ARRANGE ================================================================
	captures := []positionCollection.Capture{}
	streamName := "Pos"
	streamIn := positionCollection.NewCollection(streamName, captures)

	encoder := position.NewEncoder(position.Oct24)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
	streamOut, decodeErr := encoder.Decode(streamName, header, streamsData[0], nil)

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 1)
	assert.NotNil(t, streamOut)
	assert.Len(t, streamOut.Captures(), 0)
}

func Test_Positions(t *testing.T) {
	continuousCaptures := make([]positionCollection.Capture, 1000)
	continuousTimes := make([]float64, len(continuousCaptures))
	curTime := -1000.0
	curPos := vector3.Zero[float64]()
	for i := 0; i < len(continuousCaptures); i++ {
		continuousCaptures[i] = positionCollection.NewCapture(
			curTime,
			curPos.X(),
			curPos.Y(),
			curPos.Z(),
		)
		continuousTimes[i] = curTime
		curPos = curPos.Add(vector3.New[float64](rand.Float64()*10, rand.Float64()*10, rand.Float64()*10))
		curTime += rand.Float64() * 10.0
	}

	tests := map[string]struct {
		captures []positionCollection.Capture
		time     []float64
	}{
		"nil positions": {captures: nil},
		"0-positions":   {captures: []positionCollection.Capture{}},
		"1-positions": {
			captures: []positionCollection.Capture{positionCollection.NewCapture(1.2, 1, 1, 1)},
			time:     []float64{1.2},
		},
		"2-positions": {
			captures: []positionCollection.Capture{positionCollection.NewCapture(1.2, 1, 1, 1), positionCollection.NewCapture(1.3, 4, 5, 6)},
			time:     []float64{1.2, 1.3},
		},
		"3-positions": {
			captures: []positionCollection.Capture{
				positionCollection.NewCapture(1.2, 1, 1, 1),
				positionCollection.NewCapture(1.3, 4, 5, 6),
				positionCollection.NewCapture(1.4, 4.1, 5.7, 6.0),
			},
			time: []float64{1.2, 1.3, 1.4},
		},
		"1000-continuous-positions": {captures: continuousCaptures, time: continuousTimes},
	}

	storageTechniques := []struct {
		displayName        string
		technique          position.StorageTechnique
		timeTollerance     float64
		positionTollerance float64
	}{
		{
			displayName:        "Raw64",
			technique:          position.Raw64,
			timeTollerance:     0,
			positionTollerance: 0,
		},
		{
			displayName:        "Raw32",
			technique:          position.Raw32,
			timeTollerance:     0.0005,
			positionTollerance: 0.0003,
		},
		{
			displayName:        "Oct24",
			technique:          position.Oct24,
			timeTollerance:     0.01,
			positionTollerance: 0.04,
		},
		{
			displayName:        "Oct48",
			technique:          position.Oct48,
			timeTollerance:     0.001,
			positionTollerance: 0.0004,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				streamIn := positionCollection.NewCollection(technique.displayName, tc.captures)

				encoder := position.NewEncoder(technique.technique)

				// ACT ====================================================================
				header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
				streamOut, decodeErr := encoder.Decode(technique.displayName, header, streamsData[0], tc.time)

				// ASSERT =================================================================
				assert.NoError(t, encodeErr)
				assert.NoError(t, decodeErr)
				assert.Len(t, header, 0)
				assert.Len(t, streamsData, 1)
				if assert.NotNil(t, streamOut) {
					assert.Equal(t, streamIn.Name(), streamOut.Name())
					if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {
						for i, c := range streamOut.Captures() {
							positioniCapture, ok := c.(positionCollection.Capture)
							if assert.True(t, ok) == false {
								break
							}

							assert.InDelta(t, tc.captures[i].Time(), positioniCapture.Time(), technique.timeTollerance, "times are not equal: %.2f != %.2f", tc.captures[i].Time(), positioniCapture.Time())
							if assert.InDelta(
								t,
								tc.captures[i].Position().X(),
								positioniCapture.Position().X(),
								technique.positionTollerance,
								"X components not equal: %.2f != %.2f", tc.captures[i].Position().X(), positioniCapture.Position().X(),
							) == false {
								break
							}
							if assert.InDelta(
								t,
								tc.captures[i].Position().Y(),
								positioniCapture.Position().Y(),
								technique.positionTollerance,
								"Y components not equal: %.2f != %.2f", tc.captures[i].Position().Y(), positioniCapture.Position().Y(),
							) == false {
								break
							}
							if assert.InDelta(
								t,
								tc.captures[i].Position().Z(),
								positioniCapture.Position().Z(),
								technique.positionTollerance,
								"Z components not equal: %.2f != %.2f", tc.captures[i].Position().Z(), positioniCapture.Position().Z(),
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

func Test_Oct24_MultipleStreams(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]positionCollection.Capture, 1000)
	captureTimes := make([]float64, len(captures))
	curTime := 1.0
	for i := 0; i < len(captures); i++ {
		captures[i] = positionCollection.NewCapture(
			curTime,
			rand.Float64()*100,
			rand.Float64()*100,
			rand.Float64()*100,
		)
		captureTimes[i] = curTime
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionCollection.NewCollection("Pos", captures)

	captures2 := make([]positionCollection.Capture, 3000)
	capture2Times := make([]float64, len(captures2))
	curTime2 := 1.0
	for i := 0; i < len(captures2); i++ {
		captures2[i] = positionCollection.NewCapture(
			curTime,
			rand.Float64()*100,
			rand.Float64()*100,
			rand.Float64()*100,
		)
		capture2Times[i] = curTime
		curTime2 += rand.Float64() * 10.0
	}
	streamIn2 := positionCollection.NewCollection("Pos2", captures2)

	encoder := position.NewEncoder(position.Oct24)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn, streamIn2})
	streamOut, decodeErr := encoder.Decode("Pos", header, streamsData[0], captureTimes)
	streamOut2, decodeErr2 := encoder.Decode("Pos2", header, streamsData[1], capture2Times)

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.NoError(t, decodeErr2)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 2)
	if assert.NotNil(t, streamOut) {
		assert.Equal(t, streamIn.Name(), streamOut.Name())
		if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {
			for i, c := range streamOut.Captures() {
				positioniCapture, ok := c.(positionCollection.Capture)
				if assert.True(t, ok) == false {
					break
				}

				correct := captures[i].Position()
				answer := positioniCapture.Position()
				failureMessage := fmt.Sprintf("[%d]: (%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", i, correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.Equal(t, captures[i].Time(), positioniCapture.Time(), "Mismatched Time")
				assert.InDelta(t, captures[i].Position().X(), positioniCapture.Position().X(), .7, failureMessage)
				assert.InDelta(t, captures[i].Position().Y(), positioniCapture.Position().Y(), .7, failureMessage)
				assert.InDelta(t, captures[i].Position().Z(), positioniCapture.Position().Z(), .7, failureMessage)
			}
		}
	}

	if assert.NotNil(t, streamOut2) {
		assert.Equal(t, streamIn2.Name(), streamOut2.Name())
		if assert.Len(t, streamOut2.Captures(), len(streamIn2.Captures())) {
			for i, c := range streamOut2.Captures() {
				positioniCapture, ok := c.(positionCollection.Capture)
				if assert.True(t, ok) == false {
					break
				}

				correct := captures2[i].Position()
				answer := positioniCapture.Position()
				failureMessage := fmt.Sprintf("[%d]: (%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", i, correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.Equal(t, captures2[i].Time(), positioniCapture.Time(), "Mismatched Time")
				assert.InDelta(t, captures2[i].Position().X(), positioniCapture.Position().X(), .8, failureMessage)
				assert.InDelta(t, captures2[i].Position().Y(), positioniCapture.Position().Y(), .8, failureMessage)
				assert.InDelta(t, captures2[i].Position().Z(), positioniCapture.Position().Z(), .8, failureMessage)
			}
		}
	}
}

func Test_Oct24_Continuous(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]positionCollection.Capture, 1000)
	capturetimes := make([]float64, len(captures))
	streamName := "Pos"
	curTime := 1.0
	curPos := vector3.Zero[float64]()
	for i := 0; i < len(captures); i++ {
		captures[i] = positionCollection.NewCapture(
			curTime,
			curPos.X(),
			curPos.Y(),
			curPos.Z(),
		)
		capturetimes[i] = curTime
		curPos = curPos.Add(vector3.New[float64](rand.Float64()*10, rand.Float64()*10, rand.Float64()*10))
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionCollection.NewCollection(streamName, captures)

	encoder := position.NewEncoder(position.Oct24)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
	streamOut, decodeErr := encoder.Decode(streamName, header, streamsData[0], capturetimes)

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 1)
	if assert.NotNil(t, streamOut) {
		assert.Equal(t, streamIn.Name(), streamOut.Name())
		if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {
			for i, c := range streamOut.Captures() {
				positioniCapture, ok := c.(positionCollection.Capture)
				if assert.True(t, ok) == false {
					break
				}

				correct := captures[i].Position()
				answer := positioniCapture.Position()
				failureMessage := fmt.Sprintf("[%d]: (%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", i, correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.Equal(t, captures[i].Time(), positioniCapture.Time(), "Mismatched Time")
				assert.InDelta(t, captures[i].Position().X(), positioniCapture.Position().X(), .04, failureMessage)
				assert.InDelta(t, captures[i].Position().Y(), positioniCapture.Position().Y(), .04, failureMessage)
				assert.InDelta(t, captures[i].Position().Z(), positioniCapture.Position().Z(), .04, failureMessage)
			}
		}
	}
}
