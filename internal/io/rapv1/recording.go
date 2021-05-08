package rapv1

import "github.com/recolude/rap/format"

type recordingV1 struct {
	id             string
	name           string
	captureStreams []format.CaptureStream
	recordings     []format.Recording
	metadata       map[string]string
}

func (rec recordingV1) ID() string {
	return rec.id
}

func (rec recordingV1) Name() string {
	return rec.name
}

func (rec recordingV1) Binaries() []format.Binary {
	return nil
}

func (rec recordingV1) BinaryReferences() []format.BinaryReference {
	return nil
}

func (rec recordingV1) Metadata() map[string]string {
	return rec.metadata
}

func (rec recordingV1) Recordings() []format.Recording {
	return rec.recordings
}

func (rec recordingV1) CaptureStreams() []format.CaptureStream {
	return rec.captureStreams
}
