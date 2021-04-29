package data

import "io"

//go:generate mockgen -destination=../../internal/mocks/recording.go -package=mocks github.com/recolude/rap/pkg/data Recording,CaptureStream

type Capture interface {
	Time() float64
	String() string
}

type Binary interface {
	Name() string
	Data() io.Reader
	Size() uint64
	Write(io.Writer) (int, error)
}

type CaptureStream interface {
	Name() string
	Captures() []Capture
	Signature() string
}

type Recording interface {
	Name() string
	CaptureStreams() []CaptureStream
	Recordings() []Recording
	Metadata() map[string]string
	Binaries() []Binary
}

func NewRecording(
	name string, captureStreams []CaptureStream,
	recordings []Recording,
	metadata map[string]string,
	binaries []Binary,
) recording {
	return recording{
		name:           name,
		recordings:     recordings,
		captureStreams: captureStreams,
		metadata:       metadata,
		binaries:       binaries,
	}
}

type recording struct {
	name           string
	captureStreams []CaptureStream
	recordings     []Recording
	metadata       map[string]string
	binaries       []Binary
}

func (r recording) Name() string {
	return r.name
}

func (r recording) CaptureStreams() []CaptureStream {
	return r.captureStreams
}

func (r recording) Recordings() []Recording {
	return r.recordings
}

func (r recording) Metadata() map[string]string {
	return r.metadata
}

func (r recording) Binaries() []Binary {
	return r.binaries
}
