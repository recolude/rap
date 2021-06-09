package float

import "fmt"

type Capture struct {
	time float64
	x    float64
}

func NewCapture(time, x float64) Capture {
	return Capture{
		time: time,
		x:    x,
	}
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] - %.2f", c.time, c.x)
}

func (c Capture) Value() float64 {
	return c.x
}
