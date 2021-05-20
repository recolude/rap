package event_test

import (
	"testing"

	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCapture(t *testing.T) {
	time := 123.0
	metadata := metadata.NewBlock(map[string]metadata.Property{
		"a": metadata.NewStringProperty("b"),
		"c": metadata.NewStringProperty("d"),
	})
	capture := event.NewCapture(time, "My Name", metadata)

	// ASSERT =================================================================
	assert.Equal(t, time, capture.Time())
	assert.Equal(t, "[123.00] My Name", capture.String())
}
