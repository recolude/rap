package encoding_test

import (
	"math/rand"
	"testing"

	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/encoding"
	"github.com/recolude/rap/pkg/streams/position"
	"github.com/stretchr/testify/assert"
)

func Test_Position_Raw64(t *testing.T) {
	// ARRANGE ================================================================
	captures := make([]position.Capture, 1000)
	curTime := 0.0
	for i := 0; i < len(captures); i++ {
		captures[i] = position.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime += rand.Float64() * 10.0
	}
	streamIn := position.NewStream("Pos", captures)
	encoder := encoding.NewPositionEncoder(encoding.Raw64)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn})
	streamsOut, decodeErr := encoder.Decode(header, streamsData)

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 1)
	assert.Len(t, streamsOut, 1)
	if assert.NotNil(t, streamsOut[0]) {
		assert.Equal(t, streamIn.Name(), streamsOut[0].Name())
		if assert.Len(t, streamsOut[0].Captures(), len(streamIn.Captures())) {
			for i, c := range streamsOut[0].Captures() {
				positioniCapture, ok := c.(position.Capture)
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
	captures := make([]position.Capture, 1000)
	curTime := 0.0
	for i := 0; i < len(captures); i++ {
		captures[i] = position.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime += rand.Float64() * 10.0
	}
	streamIn := position.NewStream("Pos", captures)

	captures2 := make([]position.Capture, 3000)
	curTime2 := 0.0
	for i := 0; i < len(captures2); i++ {
		captures2[i] = position.NewCapture(
			curTime,
			rand.Float64()*10000,
			rand.Float64()*10000,
			rand.Float64()*10000,
		)
		curTime2 += rand.Float64() * 10.0
	}
	streamIn2 := position.NewStream("Pos2", captures2)

	encoder := encoding.NewPositionEncoder(encoding.Raw64)

	// ACT ====================================================================
	header, streamsData, encodeErr := encoder.Encode([]data.CaptureStream{streamIn, streamIn2})
	streamsOut, decodeErr := encoder.Decode(header, streamsData)

	// ASSERT =================================================================
	assert.NoError(t, encodeErr)
	assert.NoError(t, decodeErr)
	assert.Len(t, header, 0)
	assert.Len(t, streamsData, 2)
	assert.Len(t, streamsOut, 2)
	if assert.NotNil(t, streamsOut[0]) {
		assert.Equal(t, streamIn.Name(), streamsOut[0].Name())
		if assert.Len(t, streamsOut[0].Captures(), len(streamIn.Captures())) {
			for i, c := range streamsOut[0].Captures() {
				positioniCapture, ok := c.(position.Capture)
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

	if assert.NotNil(t, streamsOut[1]) {
		assert.Equal(t, streamIn2.Name(), streamsOut[1].Name())
		if assert.Len(t, streamsOut[1].Captures(), len(streamIn2.Captures())) {
			for i, c := range streamsOut[1].Captures() {
				positioniCapture, ok := c.(position.Capture)
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
