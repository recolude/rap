package format_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/recolude/rap/format"
	"github.com/stretchr/testify/assert"
)

func Test_StringProperty(t *testing.T) {
	tests := map[string]struct {
		value string
		data  []byte
	}{
		"empty": {value: "", data: []byte{0}},
		"a":     {value: "a", data: []byte{1, 'a'}},
		"abcd":  {value: "abcd", data: []byte{4, 'a', 'b', 'c', 'd'}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := format.NewStringProperty(tc.value)
			assert.Equal(t, byte(0), prop.Code())
			assert.Equal(t, tc.value, prop.String())
			assert.Equal(t, tc.data, prop.Data())
		})
	}
}

func Test_IntProperty(t *testing.T) {
	tests := map[string]struct {
		value     int32
		stringVal string
	}{
		"0":    {value: 0, stringVal: "0"},
		"-10":  {value: -10, stringVal: "-10"},
		"3000": {value: 3000, stringVal: "3000"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := format.NewIntProperty(tc.value)
			assert.Equal(t, byte(1), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())

			var out int32
			binary.Read(bytes.NewBuffer(prop.Data()), binary.LittleEndian, &out)
			assert.Equal(t, tc.value, out)
		})
	}
}

func Test_Float32Property(t *testing.T) {
	tests := map[string]struct {
		value     float32
		stringVal string
	}{
		"0":    {value: 0, stringVal: "0.000000"},
		"-10":  {value: -10, stringVal: "-10.000000"},
		"3000": {value: 3000, stringVal: "3000.000000"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := format.NewFloat32Property(tc.value)
			assert.Equal(t, byte(2), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())

			var out float32
			binary.Read(bytes.NewBuffer(prop.Data()), binary.LittleEndian, &out)
			assert.Equal(t, tc.value, out)
		})
	}
}

func Test_BoolProperty(t *testing.T) {
	tests := map[string]struct {
		value     bool
		stringVal string
		byteVal   int
	}{
		"true":  {value: true, stringVal: "true", byteVal: 3},
		"false": {value: false, stringVal: "false", byteVal: 4},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := format.NewBoolProperty(tc.value)

			assert.Equal(t, byte(tc.byteVal), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())
			assert.Equal(t, tc.value, prop.Value())
			assert.Len(t, prop.Data(), 0)
		})
	}
}

func Test_ByteProperty(t *testing.T) {
	tests := map[string]struct {
		value     byte
		stringVal string
	}{
		"0":  {value: 0, stringVal: "0"},
		"1":  {value: 1, stringVal: "1"},
		"12": {value: 44, stringVal: "44"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := format.NewByteProperty(tc.value)

			assert.Equal(t, byte(5), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())
			assert.Equal(t, tc.value, prop.Value())
			assert.Equal(t, []byte{tc.value}, prop.Data())
		})
	}
}
