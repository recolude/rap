package rapv1

import (
	"github.com/recolude/rap/pkg/data"
)

type recordingV1 struct {
	name           string
	captureStreams []data.CaptureStream
	recordings     []data.Recording
	metadata       map[string]string
}

func (rec recordingV1) Name() string {
	return rec.name
}

func (rec recordingV1) Binaries() []data.Binary {
	return nil
}

func (rec recordingV1) BinaryReferences() []data.BinaryReference {
	return nil
}

func (rec recordingV1) Metadata() map[string]string {
	return rec.metadata
}

func (rec recordingV1) Recordings() []data.Recording {
	return rec.recordings
}

func (rec recordingV1) CaptureStreams() []data.CaptureStream {
	return rec.captureStreams
}
