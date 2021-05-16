package format

import "io"

//go:generate mockgen -destination=../internal/mocks/recording.go -package=mocks github.com/recolude/rap/format Recording,CaptureCollection

type Capture interface {
	Time() float64
	String() string
}

type Binary interface {
	Name() string
	Data() io.Reader
	Size() uint64
	Metadata() Metadata
	Write(io.Writer) (int, error)
}

type BinaryReference interface {
	Name() string
	URI() string
	Size() uint64
	Metadata() Metadata
}

type CaptureCollection interface {
	Name() string
	Captures() []Capture
	Signature() string
}

type Recording interface {
	ID() string
	Name() string
	CaptureCollections() []CaptureCollection
	Recordings() []Recording
	Metadata() Metadata
	Binaries() []Binary
	BinaryReferences() []BinaryReference
}

func NewRecording(
	id string,
	name string,
	captureCollections []CaptureCollection,
	recordings []Recording,
	metadata Metadata,
	binaries []Binary,
) recording {
	return recording{
		id:                 id,
		name:               name,
		recordings:         recordings,
		captureCollections: captureCollections,
		metadata:           metadata,
		binaries:           binaries,
	}
}

type recording struct {
	id                 string
	name               string
	captureCollections []CaptureCollection
	recordings         []Recording
	metadata           Metadata
	binaries           []Binary
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

func (r recording) Metadata() Metadata {
	return r.metadata
}

func (r recording) Binaries() []Binary {
	return r.binaries
}

func (r recording) BinaryReferences() []BinaryReference {
	return nil
}
