package encoding

import "github.com/recolude/rap/pkg/data"

type Encoder interface {
	Accepts(data.CaptureStream) bool
	Decode(header []byte, streamData [][]byte) ([]data.CaptureStream, error)
	Encode([]data.CaptureStream) ([]byte, [][]byte, error)
	Version() uint
	Signature() string
}
