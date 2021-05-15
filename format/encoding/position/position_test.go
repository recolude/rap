package position_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/EliCDavis/vector"
	"github.com/recolude/rap/format"
	positionStream "github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/encoding/position"
	"github.com/stretchr/testify/assert"
)

func Test_Position_Raw64(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]positionStream.Capture, 1000)
	curTime := 0.0
	for i := 0; i < len(captures); i++ {
		captures[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionStream.NewStream("Pos", captures)
	encoder := position.NewEncoder(position.Raw64)

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
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}
				assert.Equal(t, captures[i].Time(), positioniCapture.Time())
				assert.Equal(t, captures[i].Position().X(), positioniCapture.Position().X())
				assert.Equal(t, captures[i].Position().Y(), positioniCapture.Position().Y())
				assert.Equal(t, captures[i].Position().Z(), positioniCapture.Position().Z())
			}
		}
	}
}

func Test_Position_MultipleStreams(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]positionStream.Capture, 1000)
	curTime := 0.0
	for i := 0; i < len(captures); i++ {
		captures[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionStream.NewStream("Pos", captures)

	captures2 := make([]positionStream.Capture, 3000)
	curTime2 := 0.0
	for i := 0; i < len(captures2); i++ {
		captures2[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime2 += rand.Float64() * 10.0
	}
	streamIn2 := positionStream.NewStream("Pos2", captures2)

	encoder := position.NewEncoder(position.Raw64)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn, streamIn2})
	streamOut, decodeErr := encoder.Decode(header, streamsData[0])
	streamOut2, decodeErr2 := encoder.Decode(header, streamsData[1])

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
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}
				assert.Equal(t, captures[i].Time(), positioniCapture.Time())
				assert.Equal(t, captures[i].Position().X(), positioniCapture.Position().X())
				assert.Equal(t, captures[i].Position().Y(), positioniCapture.Position().Y())
				assert.Equal(t, captures[i].Position().Z(), positioniCapture.Position().Z())
			}
		}
	}

	if assert.NotNil(t, streamOut2) {
		assert.Equal(t, streamIn2.Name(), streamOut2.Name())
		if assert.Len(t, streamOut2.Captures(), len(streamIn2.Captures())) {
			for i, c := range streamOut2.Captures() {
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}
				assert.Equal(t, captures2[i].Time(), positioniCapture.Time())
				assert.Equal(t, captures2[i].Position().X(), positioniCapture.Position().X())
				assert.Equal(t, captures2[i].Position().Y(), positioniCapture.Position().Y())
				assert.Equal(t, captures2[i].Position().Z(), positioniCapture.Position().Z())
			}
		}
	}
}

func Test_Raw32_MultipleStreams(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]positionStream.Capture, 1000)
	curTime := 1.0
	for i := 0; i < len(captures); i++ {
		captures[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionStream.NewStream("Pos", captures)

	captures2 := make([]positionStream.Capture, 3000)
	curTime2 := 1.0
	for i := 0; i < len(captures2); i++ {
		captures2[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime2 += rand.Float64() * 10.0
	}
	streamIn2 := positionStream.NewStream("Pos2", captures2)

	encoder := position.NewEncoder(position.Raw32)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn, streamIn2})
	streamOut, decodeErr := encoder.Decode(header, streamsData[0])
	streamOut2, decodeErr2 := encoder.Decode(header, streamsData[1])

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
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}

				assert.InEpsilon(t, captures[i].Time(), positioniCapture.Time(), 0.000001)
				assert.InEpsilon(t, captures[i].Position().X(), positioniCapture.Position().X(), 0.000001)
				assert.InEpsilon(t, captures[i].Position().Y(), positioniCapture.Position().Y(), 0.000001)
				assert.InEpsilon(t, captures[i].Position().Z(), positioniCapture.Position().Z(), 0.000001)
			}
		}
	}

	if assert.NotNil(t, streamOut2) {
		assert.Equal(t, streamIn2.Name(), streamOut2.Name())
		if assert.Len(t, streamOut2.Captures(), len(streamIn2.Captures())) {
			for i, c := range streamOut2.Captures() {
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}
				assert.InEpsilon(t, captures2[i].Time(), positioniCapture.Time(), 0.000001)
				assert.InEpsilon(t, captures2[i].Position().X(), positioniCapture.Position().X(), 0.000001)
				assert.InEpsilon(t, captures2[i].Position().Y(), positioniCapture.Position().Y(), 0.000001)
				assert.InEpsilon(t, captures2[i].Position().Z(), positioniCapture.Position().Z(), 0.000001)
			}
		}
	}
}

func Test_Oct24_EmptyStream(t *testing.T) {
	// ARRANGE ================================================================
	captures := []positionStream.Capture{}

	streamIn := positionStream.NewStream("Pos", captures)

	encoder := position.NewEncoder(position.Oct24)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
	streamOut, decodeErr := encoder.Decode(header, streamsData[0])

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 1)
	assert.NotNil(t, streamOut)
	assert.Len(t, streamOut.Captures(), 0)
}

func Test_Positions(t *testing.T) {
	continuousCaptures := make([]positionStream.Capture, 1000)
	curTime := 1.0
	curPos := vector.Vector3Zero()
	for i := 0; i < len(continuousCaptures); i++ {
		continuousCaptures[i] = positionStream.NewCapture(
			curTime,
			curPos.X(),
			curPos.Y(),
			curPos.Z(),
		)
		curPos = curPos.Add(vector.NewVector3(rand.Float64()*10, rand.Float64()*10, rand.Float64()*10))
		curTime += rand.Float64() * 10.0
	}

	tests := map[string]struct {
		captures []positionStream.Capture
	}{
		"nil positions": {captures: nil},
		"0-positions":   {captures: []positionStream.Capture{}},
		"1-positions":   {captures: []positionStream.Capture{positionStream.NewCapture(1.2, 1, 1, 1)}},
		"2-positions":   {captures: []positionStream.Capture{positionStream.NewCapture(1.2, 1, 1, 1), positionStream.NewCapture(1.3, 4, 5, 6)}},
		"3-positions": {
			captures: []positionStream.Capture{
				positionStream.NewCapture(1.2, 1, 1, 1),
				positionStream.NewCapture(1.3, 4, 5, 6),
				positionStream.NewCapture(1.4, 4.1, 5.7, 6.0),
			},
		},
		"1000-continuous-positions": {captures: continuousCaptures},
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
			timeTollerance:     0.01,
			positionTollerance: 0.0004,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				streamIn := positionStream.NewStream("Pos", tc.captures)

				encoder := position.NewEncoder(technique.technique)

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
							positioniCapture, ok := c.(positionStream.Capture)
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
	captures := make([]positionStream.Capture, 1000)
	curTime := 1.0
	for i := 0; i < len(captures); i++ {
		captures[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*100,
			rand.Float64()*100,
			rand.Float64()*100,
		)
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionStream.NewStream("Pos", captures)

	captures2 := make([]positionStream.Capture, 3000)
	curTime2 := 1.0
	for i := 0; i < len(captures2); i++ {
		captures2[i] = positionStream.NewCapture(
			curTime,
			rand.Float64()*100,
			rand.Float64()*100,
			rand.Float64()*100,
		)
		curTime2 += rand.Float64() * 10.0
	}
	streamIn2 := positionStream.NewStream("Pos2", captures2)

	encoder := position.NewEncoder(position.Oct24)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn, streamIn2})
	streamOut, decodeErr := encoder.Decode(header, streamsData[0])
	streamOut2, decodeErr2 := encoder.Decode(header, streamsData[1])

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
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}

				correct := captures[i].Position()
				answer := positioniCapture.Position()
				failureMessage := fmt.Sprintf("[%d]: (%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", i, correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.InDelta(t, captures[i].Time(), positioniCapture.Time(), 0.003, "Mismatched Time")
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
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}

				correct := captures2[i].Position()
				answer := positioniCapture.Position()
				failureMessage := fmt.Sprintf("[%d]: (%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", i, correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.InDelta(t, captures2[i].Time(), positioniCapture.Time(), 0.003, "Mismatched Time")
				assert.InDelta(t, captures2[i].Position().X(), positioniCapture.Position().X(), .8, failureMessage)
				assert.InDelta(t, captures2[i].Position().Y(), positioniCapture.Position().Y(), .8, failureMessage)
				assert.InDelta(t, captures2[i].Position().Z(), positioniCapture.Position().Z(), .8, failureMessage)
			}
		}
	}
}

func Test_Oct24_Continuous(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]positionStream.Capture, 1000)
	curTime := 1.0
	curPos := vector.Vector3Zero()
	for i := 0; i < len(captures); i++ {
		captures[i] = positionStream.NewCapture(
			curTime,
			curPos.X(),
			curPos.Y(),
			curPos.Z(),
		)
		curPos = curPos.Add(vector.NewVector3(rand.Float64()*10, rand.Float64()*10, rand.Float64()*10))
		curTime += rand.Float64() * 10.0
	}
	streamIn := positionStream.NewStream("Pos", captures)

	encoder := position.NewEncoder(position.Oct24)

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
				positioniCapture, ok := c.(positionStream.Capture)
				if assert.True(t, ok) == false {
					break
				}

				correct := captures[i].Position()
				answer := positioniCapture.Position()
				failureMessage := fmt.Sprintf("[%d]: (%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", i, correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.InDelta(t, captures[i].Time(), positioniCapture.Time(), 0.0003, "Mismatched Time")
				assert.InDelta(t, captures[i].Position().X(), positioniCapture.Position().X(), .04, failureMessage)
				assert.InDelta(t, captures[i].Position().Y(), positioniCapture.Position().Y(), .04, failureMessage)
				assert.InDelta(t, captures[i].Position().Z(), positioniCapture.Position().Z(), .04, failureMessage)
			}
		}
	}
}
