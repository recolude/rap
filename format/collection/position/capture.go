package position

import (
	"fmt"

	"github.com/EliCDavis/vector/vector3"
)

type Capture struct {
	time     float64
	position vector3.Float64
}

func NewCapture(time, x, y, z float64) Capture {
	return Capture{
		time:     time,
		position: vector3.New[float64](x, y, z),
	}
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] - %.2f, %.2f, %.2f", c.time, c.position.X(), c.position.Y(), c.position.Z())
}

func (c Capture) Position() vector3.Float64 {
	return c.position
}
