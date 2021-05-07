package position_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/encoding/position"
	positionStream "github.com/recolude/rap/pkg/streams/position"
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
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn})
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
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn, streamIn2})
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
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn, streamIn2})
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

func Test_Oct24_SingleEntry(t *testing.T) {
	// ARRANGE ================================================================
	captures := []positionStream.Capture{
		positionStream.NewCapture(1, 2, 3, 4),
		positionStream.NewCapture(2, 4, 5, 6),
		positionStream.NewCapture(3, 7, 8, 9),
	}

	streamIn := positionStream.NewStream("Pos", captures)

	encoder := position.NewEncoder(position.Oct24)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn})
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

				assert.InDelta(t, captures[i].Time(), positioniCapture.Time(), 0.01, "times are not equal: %.2f != %.2f", captures[i].Time(), positioniCapture.Time())
				assert.InDelta(
					t,
					captures[i].Position().X(),
					positioniCapture.Position().X(),
					0.01,
					"X components not equal: %.2f != %.2f", captures[i].Position().X(), positioniCapture.Position().X(),
				)
				assert.InDelta(
					t,
					captures[i].Position().Y(),
					positioniCapture.Position().Y(),
					0.01,
					"Y components not equal: %.2f != %.2f", captures[i].Position().Y(), positioniCapture.Position().Y(),
				)
				assert.InDelta(
					t,
					captures[i].Position().Z(),
					positioniCapture.Position().Z(),
					0.01,
					"Z components not equal: %.2f != %.2f", captures[i].Position().Z(), positioniCapture.Position().Z(),
				)
			}
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
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn, streamIn2})
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
				failureMessage := fmt.Sprintf("(%.2f, %.2f, %.2f) != (%.2f, %.2f, %.2f)", correct.X(), correct.Y(), correct.Z(), answer.X(), answer.Y(), answer.Z())

				assert.InDelta(t, captures[i].Time(), positioniCapture.Time(), 0.0003, "Mismatched Time")
				assert.InDelta(t, captures[i].Position().X(), positioniCapture.Position().X(), .2, failureMessage)
				assert.InDelta(t, captures[i].Position().Y(), positioniCapture.Position().Y(), .2, failureMessage)
				assert.InDelta(t, captures[i].Position().Z(), positioniCapture.Position().Z(), .2, failureMessage)
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
				assert.InDelta(t, captures2[i].Time(), positioniCapture.Time(), 0.0002, "Mismatched Time")
				assert.InDelta(t, captures2[i].Position().X(), positioniCapture.Position().X(), .2)
				assert.InDelta(t, captures2[i].Position().Y(), positioniCapture.Position().Y(), .2)
				assert.InDelta(t, captures2[i].Position().Z(), positioniCapture.Position().Z(), .2)
			}
		}
	}
}
