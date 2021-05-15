package enum

import (
	"fmt"
)

type Capture struct {
	time  float64
	value int
}

func NewCapture(time float64, value int) Capture {
	return Capture{
		time:  time,
		value: value,
	}
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) Value() int {
	return c.value
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] Enum - %d", c.time, c.value)
}
