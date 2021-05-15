package position

import (
	"fmt"

	"github.com/EliCDavis/vector"
)

type Capture struct {
	time     float64
	position vector.Vector3
}

func NewCapture(time, x, y, z float64) Capture {
	return Capture{
		time:     time,
		position: vector.NewVector3(x, y, z),
	}
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] - %.2f, %.2f, %.2f", c.time, c.position.X(), c.position.Y(), c.position.Z())
}

func (c Capture) Position() vector.Vector3 {
	return c.position
}
