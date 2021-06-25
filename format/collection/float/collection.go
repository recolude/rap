package float

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

func (s Collection) Name() string {
	return s.name
}

func (s Collection) Captures() []format.Capture {
	returnVal := make([]format.Capture, len(s.captures))
	for i := range s.captures {
		returnVal[i] = s.captures[i]
	}
	return returnVal
}

func (Collection) Signature() string {
	return "recolude.float"
}

func (c Collection) Slice(beginning, end float64) format.CaptureCollection {
	slicedCaptures := make([]Capture, 0)
	for _, c := range c.captures {
		if format.CaptureFallsWithin(c, beginning, end) {
			slicedCaptures = append(slicedCaptures, c)
		}
	}
	return NewCollection(c.Name(), slicedCaptures)
}
