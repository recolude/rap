package position

import (
	"github.com/recolude/rap/format"
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

func (s Stream) Captures() []format.Capture {
	returnVal := make([]format.Capture, len(s.captures))
	for i := range s.captures {
		returnVal[i] = s.captures[i]
	}
	return returnVal
}

func (Stream) Signature() string {
	return "recolude.position"
}
