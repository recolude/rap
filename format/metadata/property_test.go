package metadata_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/recolude/rap/format/metadata"
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
			prop := metadata.NewStringProperty(tc.value)
			assert.Equal(t, byte(0), prop.Code())
			assert.Equal(t, tc.value, prop.String())
			assert.Equal(t, tc.data, prop.Data())

			b, err := json.Marshal(prop)
			assert.Nil(t, err)

			var sp metadata.StringProperty
			assert.Nil(t, json.Unmarshal(b, &sp))
			assert.Equal(t, prop, sp)
		})
	}
}

func Test_IntProperty(t *testing.T) {
	tests := map[string]struct {
		value     int
		stringVal string
	}{
		"0":    {value: 0, stringVal: "0"},
		"-10":  {value: -10, stringVal: "-10"},
		"3000": {value: 3000, stringVal: "3000"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := metadata.NewIntProperty(tc.value)
			assert.Equal(t, byte(1), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())

			var out int32
			binary.Read(bytes.NewBuffer(prop.Data()), binary.LittleEndian, &out)
			assert.Equal(t, int32(tc.value), out)

			b, err := json.Marshal(prop)
			assert.Nil(t, err)

			var ip metadata.Int32Property
			assert.Nil(t, json.Unmarshal(b, &ip))
			assert.Equal(t, prop, ip)
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
			prop := metadata.NewFloat32Property(tc.value)
			assert.Equal(t, byte(2), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())

			var out float32
			binary.Read(bytes.NewBuffer(prop.Data()), binary.LittleEndian, &out)
			assert.Equal(t, tc.value, out)

			b, err := json.Marshal(prop)
			assert.Nil(t, err)

			var fp metadata.Float32Property
			assert.Nil(t, json.Unmarshal(b, &fp))
			assert.Equal(t, prop, fp)
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
			prop := metadata.NewBoolProperty(tc.value)

			assert.Equal(t, byte(tc.byteVal), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())
			assert.Equal(t, tc.value, prop.Value())
			assert.Len(t, prop.Data(), 0)

			b, err := json.Marshal(prop)
			assert.Nil(t, err)

			var bp metadata.BoolProperty
			assert.Nil(t, json.Unmarshal(b, &bp))
			assert.Equal(t, prop, bp)
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
		"44": {value: 44, stringVal: "44"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prop := metadata.NewByteProperty(tc.value)

			assert.Equal(t, byte(5), prop.Code())
			assert.Equal(t, tc.stringVal, prop.String())
			assert.Equal(t, tc.value, prop.Value())
			assert.Equal(t, []byte{tc.value}, prop.Data())

			b, err := json.Marshal(prop)
			assert.Nil(t, err)

			var bp metadata.ByteProperty
			assert.Nil(t, json.Unmarshal(b, &bp))
			assert.Equal(t, prop, bp)
		})
	}
}
