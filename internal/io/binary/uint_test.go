package binary_test

import (
	"bytes"
	"testing"

	"github.com/recolude/rap/internal/io/binary"
	"github.com/stretchr/testify/assert"
)

func Test_UintArrayToBytes(t *testing.T) {
	tests := map[string]struct {
		input []uint
	}{
		"empty array":         {input: []uint{}},
		"one element array":   {input: []uint{0}},
		"one element array 1": {input: []uint{1}},
		"multi element array": {input: []uint{1, 2, 3, 4, 5}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := binary.UintArrayToBytes(tc.input)
			back, _, err := binary.ReadUintArray(bytes.NewBuffer(got))
			if assert.NoError(t, err) && assert.Equal(t, len(tc.input), len(back)) {
				for i, correct := range tc.input {
					assert.Equal(t, correct, back[i])
				}
			}
		})
	}
}
