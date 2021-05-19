package event_test

import (
	"testing"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/event"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCollection(t *testing.T) {
	metadata := format.NewMetadataBlock(map[string]format.Property{
		"a": format.NewStringProperty("b"),
		"c": format.NewStringProperty("d"),
	})

	event1 := event.NewCapture(1, "Event 1", metadata)
	event2 := event.NewCapture(2, "Event 2", metadata)
	name := "Custom Events"

	// ACT ====================================================================
	collection := event.NewCollection(name, []event.Capture{event1, event2})
	captures := collection.Captures()

	// ASSERT =================================================================
	assert.Equal(t, name, collection.Name())
	assert.Len(t, captures, 2)
}
