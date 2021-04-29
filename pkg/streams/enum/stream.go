package enum

import (
	"github.com/recolude/rap/pkg/data"
)

type Stream struct {
	name        string
	enumMembers []string
	captures    []Capture
}

func NewStream(name string, enumMembers []string, captures []Capture) Stream {
	return Stream{
		name:        name,
		enumMembers: enumMembers,
		captures:    captures,
	}
}

func (s Stream) Name() string {
	return s.name
}

func (s Stream) Captures() []data.Capture {
	returnVal := make([]data.Capture, len(s.captures))
	for i := range s.captures {
		returnVal[i] = s.captures[i]
	}
	return returnVal
}

func (Stream) Signature() string {
	return "recolude.enum"
}
