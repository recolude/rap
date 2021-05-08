package event_test

import (
	"testing"

	"github.com/recolude/rap/format/streams/event"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCapture(t *testing.T) {
	time := 123.0

	capture := event.NewCapture(time, "My Name", map[string]string{"a": "b", "c": "d"})

	// ASSERT =================================================================
	assert.Equal(t, time, capture.Time())
	assert.Equal(t, "[123.00] My Name", capture.String())
}
