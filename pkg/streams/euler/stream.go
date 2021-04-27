package euler

import (
	"github.com/recolude/rap/pkg/data"
)

type Stream struct {
	name     string
	captures []Capture
}

func NewStream(name string, captures []Capture) Stream {
	return Stream{
		name:     name,
		captures: captures,
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
