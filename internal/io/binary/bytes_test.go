package binary_test

import (
	"bytes"
	"testing"

	"github.com/recolude/rap/internal/io/binary"
	"github.com/stretchr/testify/assert"
)

func Test_ByteArrayToBytes(t *testing.T) {
	tests := map[string]struct {
		input []byte
	}{
		"empty array":           {input: []byte{}},
		"single element array":  {input: []byte{0}},
		"multi element array":   {input: []byte{1, 0}},
		"multi element array 2": {input: []byte{0, 0}},
		"multi element array 3": {input: []byte{4, 10, 3, 88}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := binary.BytesArrayToBytes(tc.input)
			back, _, err := binary.ReadBytesArray(bytes.NewBuffer(got))
			if assert.NoError(t, err) {
				assert.Equal(t, tc.input, back)
			}
		})
	}
}
