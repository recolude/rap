package position

import (
	"github.com/recolude/rap/format"
)

type Collection struct {
	name     string
	captures []Capture
}

func NewCollection(name string, captures []Capture) Collection {
	return Collection{
		name:     name,
		captures: captures,
	}
}

func (c Collection) Name() string {
	return c.name
}

func (Collection) Signature() string {
	return "recolude.position"
}

func (c Collection) Captures() []format.Capture {
	returnVal := make([]format.Capture, len(c.captures))
	for i := range c.captures {
		returnVal[i] = c.captures[i]
	}
	return returnVal
}

func (c Collection) Slice(beginning, end float64) format.CaptureCollection {
	slicedCaptures := make([]Capture, 0)
	for _, c := range c.captures {
		if c.Time() >= beginning && c.Time() < end {
			slicedCaptures = append(slicedCaptures, c)
		}
	}
	return NewCollection(c.Name(), slicedCaptures)
}
