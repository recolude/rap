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
	return "recolude.position"
}
