package encoding

import "github.com/recolude/rap/format"

type Encoder interface {
	Accepts(format.CaptureCollection) bool
	Decode(header []byte, streamData []byte, times []float64) (format.CaptureCollection, error)
	Encode([]format.CaptureCollection) ([]byte, [][]byte, error)
	Version() uint
	Signature() string
}
