package euler

import (
	"fmt"

	"github.com/EliCDavis/vector/vector3"
)

type Capture struct {
	time  float64
	euler vector3.Float64
}

func NewEulerZXYCapture(time, eulerX, eulerY, eulerZ float64) Capture {
	return Capture{
		time:  time,
		euler: vector3.New(eulerX, eulerY, eulerZ),
	}
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) EulerZXY() vector3.Float64 {
	return c.euler
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] Rotation - %.2f, %.2f, %.2f", c.time, c.euler.X(), c.euler.Y(), c.euler.Z())
}
