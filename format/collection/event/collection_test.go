package event_test

import (
	"testing"

	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_Collection(t *testing.T) {
	block := metadata.NewBlock(map[string]metadata.Property{
		"a": metadata.NewStringProperty("b"),
		"c": metadata.NewStringProperty("d"),
	})

	event1 := event.NewCapture(1, "Event 1", block)
	event2 := event.NewCapture(2, "Event 2", block)
	event3 := event.NewCapture(3, "Event 3", block)
	event4 := event.NewCapture(4, "Event 4", block)
	name := "Custom Events"

	// ACT ====================================================================
	collection := event.NewCollection(name, []event.Capture{event1, event2, event3, event4})
	captures := collection.Captures()
	slicedCaptures := collection.Slice(1.5, 3.5)

	// ASSERT =================================================================
	assert.Equal(t, name, collection.Name())
	assert.Equal(t, "recolude.event", collection.Signature())
	assert.Len(t, captures, 4)
	// assert.Equal(t, "Event 1", captures[0].)

	assert.Len(t, slicedCaptures.Captures(), 2)
}
