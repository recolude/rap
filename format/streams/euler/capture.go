package euler

import (
	"fmt"

	"github.com/EliCDavis/vector"
)

type Capture struct {
	time  float64
	euler vector.Vector3
}

func NewEulerZXYCapture(time, eulerX, eulerY, eulerZ float64) Capture {
	return Capture{
		time:  time,
		euler: vector.NewVector3(eulerX, eulerY, eulerZ),
	}
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) EulerZXY() vector.Vector3 {
	return c.euler
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] Rotation - %.2f, %.2f, %.2f", c.time, c.euler.X(), c.euler.Y(), c.euler.Z())
}
