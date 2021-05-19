package event_test

import (
	"testing"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/event"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCapture(t *testing.T) {
	time := 123.0
	metadata := format.NewMetadataBlock(map[string]format.Property{
		"a": format.NewStringProperty("b"),
		"c": format.NewStringProperty("d"),
	})
	capture := event.NewCapture(time, "My Name", metadata)

	// ASSERT =================================================================
	assert.Equal(t, time, capture.Time())
	assert.Equal(t, "[123.00] My Name", capture.String())
}
