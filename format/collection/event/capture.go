package event

import (
	"fmt"
)

type Capture struct {
	time     float64
	name     string
	metadata map[string]string
}

func NewCapture(time float64, name string, metadata map[string]string) Capture {
	return Capture{
		time:     time,
		name:     name,
		metadata: metadata,
	}
}

func (c Capture) Name() string {
	return c.name
}

func (c Capture) Metadata() map[string]string {
	return c.metadata
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] %s", c.time, c.name)
}
