package data

import "io"

type Capture interface {
	Time() float64
	String() string
}

type Binary interface {
	Name() string
	Data() io.Reader
	Write(io.Writer) (int, error)
}

type CaptureStream interface {
	Name() string
	Captures() []Capture
}

type Recording interface {
	Name() string
	CaptureStreams() []CaptureStream
	Recordings() []Recording
	Metadata() map[string]string
	Binaries() []Binary
}
