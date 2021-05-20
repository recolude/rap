package event_test

import (
	"testing"

	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCollection(t *testing.T) {
	block := metadata.NewBlock(map[string]metadata.Property{
		"a": metadata.NewStringProperty("b"),
		"c": metadata.NewStringProperty("d"),
	})

	event1 := event.NewCapture(1, "Event 1", block)
	event2 := event.NewCapture(2, "Event 2", block)
	name := "Custom Events"

	// ACT ====================================================================
	collection := event.NewCollection(name, []event.Capture{event1, event2})
	captures := collection.Captures()

	// ASSERT =================================================================
	assert.Equal(t, name, collection.Name())
	assert.Len(t, captures, 2)
}
