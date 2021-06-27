package event_test

import (
	"testing"

	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCapture(t *testing.T) {
	time := 123.0
	block := metadata.NewBlock(map[string]metadata.Property{
		"a": metadata.NewStringProperty("b"),
		"c": metadata.NewStringProperty("d"),
	})
	capture := event.NewCapture(time, "My Name", block)

	// ASSERT =================================================================
	assert.Equal(t, time, capture.Time())
	assert.Equal(t, "[123.00] My Name", capture.String())
	assert.Len(t, capture.Metadata().Mapping(), 2)
	assert.Equal(t, metadata.NewStringProperty("b"), capture.Metadata().Mapping()["a"])
	assert.Equal(t, metadata.NewStringProperty("d"), capture.Metadata().Mapping()["c"])
}
