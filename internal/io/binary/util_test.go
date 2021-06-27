package binary_test

import (
	"testing"

	"github.com/recolude/rap/internal/io/binary"
	"github.com/stretchr/testify/assert"
)

func Test_FloatBST(t *testing.T) {
	tests := map[string]struct {
		start     float64
		duration  float64
		value     float64
		bufSize   int
		tolerance float64
	}{
		"[-1, 1][1]: 0": {start: -1, duration: 2, value: 0, bufSize: 1, tolerance: 0.01},
		"[-1, 1][2]: 0": {start: -1, duration: 2, value: 0, bufSize: 2, tolerance: 0.001},
		"[-1, 1][4]: 0": {start: -1, duration: 2, value: 0, bufSize: 4, tolerance: 0.000001},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			buf := make([]byte, tc.bufSize)
			binary.UnsignedFloatBSTToBytes(tc.value, tc.start, tc.duration, buf)
			back := binary.BytesToUnisngedFloatBST(tc.start, tc.duration, buf)
			assert.InDelta(t, tc.value, back, tc.tolerance)
		})
	}
}
