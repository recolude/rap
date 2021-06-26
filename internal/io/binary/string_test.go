package binary_test

import (
	"bytes"
	"testing"

	"github.com/recolude/rap/internal/io/binary"
	"github.com/stretchr/testify/assert"
)

func Test_StringToBytes(t *testing.T) {
	tests := map[string]struct {
		input string
	}{
		"empty string":              {input: ""},
		"single letter":             {input: "a"},
		"word":                      {input: "soomething"},
		"long word above 256 chars": {input: "somethingggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"},
		"special character":         {input: "ABC€"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := binary.StringToBytes(tc.input)
			back, _, err := binary.ReadString(bytes.NewBuffer(got))
			if assert.NoError(t, err) {
				assert.Equal(t, tc.input, back)
			}
		})
	}
}

func Test_StringArrayToBytes(t *testing.T) {
	tests := map[string]struct {
		input []string
	}{
		"empty array":                    {input: []string{}},
		"array with single empty string": {input: []string{""}},
		"array of empty strings":         {input: []string{"", "", "", ""}},
		"array of single letters":        {input: []string{"a", "b"}},
		"array of words":                 {input: []string{"apple", "bannanna", "kiwi"}},
		"array of words 2":               {input: []string{"apple", "", "kiwi"}},
		"array of words with empty":      {input: []string{"apple", "", ""}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := binary.StringArrayToBytes(tc.input)
			back, _, err := binary.ReadStringArray(bytes.NewBuffer(got))
			if assert.NoError(t, err) && assert.Equal(t, len(tc.input), len(back)) {
				for i, correct := range tc.input {
					assert.Equal(t, correct, back[i])
				}
			}
		})
	}
}
