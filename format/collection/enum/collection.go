package enum

import (
	"github.com/recolude/rap/format"
)

type StorageTechnique int

const (
	// Raw64 encodes time with 64 bit precision
	Raw64 StorageTechnique = iota

	// Raw32 encodes time with 32 bit precision
	Raw32
)

type Collection struct {
	name             string
	enumMembers      []string
	captures         []Capture
	storageTechnique StorageTechnique
}

func NewCollection(name string, storageTechnique StorageTechnique, enumMembers []string, captures []Capture) Collection {
	return Collection{
		name:             name,
		enumMembers:      enumMembers,
		captures:         captures,
		storageTechnique: storageTechnique,
	}
}

func (s Collection) StorageTechnique() StorageTechnique {
	return s.storageTechnique
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
