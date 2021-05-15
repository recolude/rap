package encoding

import "github.com/recolude/rap/format"

type Encoder interface {
	Accepts(format.CaptureCollection) bool
	Decode(header []byte, streamData []byte) (format.CaptureCollection, error)
	Encode([]format.CaptureCollection) ([]byte, [][]byte, error)
	Version() uint
	Signature() string
}
