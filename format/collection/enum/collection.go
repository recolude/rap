package enum

import (
	"github.com/recolude/rap/format"
)

type Collection struct {
	name        string
	enumMembers []string
	captures    []Capture
}

func NewCollection(name string, enumMembers []string, captures []Capture) Collection {
	return Collection{
		name:        name,
		enumMembers: enumMembers,
		captures:    captures,
	}
}

func (s Collection) Name() string {
	return s.name
}

func (s Collection) EnumMembers() []string {
	return s.enumMembers
}

func (s Collection) Captures() []format.Capture {
	returnVal := make([]format.Capture, len(s.captures))
	for i := range s.captures {
		returnVal[i] = s.captures[i]
	}
	return returnVal
}

func (Collection) Signature() string {
	return "recolude.enum"
}

func (c Collection) Slice(beginning, end float64) format.CaptureCollection {
	slicedCaptures := make([]Capture, 0)
	for _, c := range c.captures {
		if format.CaptureFallsWithin(c, beginning, end) {
			slicedCaptures = append(slicedCaptures, c)
		}
	}
	return NewCollection(c.Name(), c.EnumMembers(), slicedCaptures)
}
