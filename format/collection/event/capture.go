package event

import (
	"fmt"

	"github.com/recolude/rap/format/metadata"
)

type Capture struct {
	time  float64
	name  string
	block metadata.Block
}

func NewCapture(time float64, name string, block metadata.Block) Capture {
	return Capture{
		time:  time,
		name:  name,
		block: block,
	}
}

func (c Capture) Name() string {
	return c.name
}

func (c Capture) Metadata() metadata.Block {
	return c.block
}

func (c Capture) Time() float64 {
	return c.time
}

func (c Capture) String() string {
	return fmt.Sprintf("[%.2f] %s", c.time, c.name)
}
