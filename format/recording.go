package format

import (
	"io"

	"github.com/recolude/rap/format/metadata"
)

//go:generate mockgen -destination=../internal/mocks/recording.go -package=mocks github.com/recolude/rap/format Recording,CaptureCollection,Binary,BinaryReference

type Capture interface {
	Time() float64
	String() string
}

type Binary interface {
	Name() string
	Data() io.Reader
	Size() uint64
	Metadata() metadata.Block
}

type BinaryReference interface {
	Name() string
	URI() string
	Size() uint64
	Metadata() metadata.Block
}

type CaptureCollection interface {
	Name() string
	Captures() []Capture
	Signature() string
	Slice(beginning, end float64) CaptureCollection
}

type Recording interface {
	ID() string
	Name() string
	CaptureCollections() []CaptureCollection
	Recordings() []Recording
	Metadata() metadata.Block
	Binaries() []Binary
	BinaryReferences() []BinaryReference
}

func NewRecording(
	id string,
	name string,
	captureCollections []CaptureCollection,
	recordings []Recording,
	metadata metadata.Block,
	binaries []Binary,
	binaryReferences []BinaryReference,
) recording {
	return recording{
		id:                 id,
		name:               name,
		recordings:         recordings,
		captureCollections: captureCollections,
		metadata:           metadata,
		binaries:           binaries,
		binaryReferences:   binaryReferences,
	}
}

type recording struct {
	id                 string
	name               string
	captureCollections []CaptureCollection
	recordings         []Recording
	metadata           metadata.Block
	binaries           []Binary
	binaryReferences   []BinaryReference
}

func (r recording) ID() string {
	return r.id
}

func (r recording) Name() string {
	return r.name
}

func (r recording) CaptureCollections() []CaptureCollection {
	return r.captureCollections
}

func (r recording) Recordings() []Recording {
	return r.recordings
}

func (r recording) Metadata() metadata.Block {
	return r.metadata
}

func (r recording) Binaries() []Binary {
	return r.binaries
}

func (r recording) BinaryReferences() []BinaryReference {
	return r.binaryReferences
}
