package event_test

import (
	"testing"

	"github.com/recolude/rap/format/collection/event"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCollection(t *testing.T) {
	event1 := event.NewCapture(1, "Event 1", map[string]string{"a": "b", "c": "d"})
	event2 := event.NewCapture(2, "Event 2", map[string]string{"a": "b", "c": "d"})
	name := "Custom Events"

	// ACT ====================================================================
	collection := event.NewCollection(name, []event.Capture{event1, event2})
	captures := collection.Captures()

	// ASSERT =================================================================
	assert.Equal(t, name, collection.Name())
	assert.Len(t, captures, 2)
}
