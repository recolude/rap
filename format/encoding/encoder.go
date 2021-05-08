package encoding

import "github.com/recolude/rap/format"

type Encoder interface {
	Accepts(format.CaptureStream) bool
	Decode(header []byte, streamData []byte) (format.CaptureStream, error)
	Encode([]format.CaptureStream) ([]byte, [][]byte, error)
	Version() uint
	Signature() string
}
