package metadata_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"
	"time"

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

func Test_MetadataBlockProperty_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	strProp := metadata.NewStringProperty("Meeee")
	float32Prop := metadata.NewFloat32Property(1234.5678)
	float32ArrProp := metadata.NewFloat32ArrayProperty([]float32{1.2, 3.4, 5.6})
	byteProp := metadata.NewByteProperty(123)
	byteProp2 := metadata.NewByteProperty(1)
	boolProp := metadata.NewBoolProperty(true)
	boolProp2 := metadata.NewBoolProperty(false)
	boolArrProp := metadata.NewBoolArrayProperty([]bool{true, false, false})
	timeProp := metadata.NewTimeProperty(time.Unix(1234567, 0))
	vecProp := metadata.NewVector2Property(999, 444)
	vec3Prop := metadata.NewVector3Property(999, 444, 222)
	int32Prop := metadata.NewIntProperty(888888)
	binArrProp := metadata.NewBinaryArrayProperty([]byte{1, 2, 3, 4, 5, 6})
	blockProp := metadata.NewMetadataProperty(metadata.NewBlock(
		map[string]metadata.Property{
			"STR":         strProp,
			"BYTE":        byteProp,
			"BYTE2":       byteProp2,
			"BOOL_ARR":    boolArrProp,
			"FLOAT32":     float32Prop,
			"FLOAT32_ARR": float32ArrProp,
			"BOOL_TRUE":   boolProp,
			"BOOL_FALSE":  boolProp2,
			"TIME":        timeProp,
			"VEC2":        vecProp,
			"VEC3":        vec3Prop,
			"INT32":       int32Prop,
			"BIN_ARR":     binArrProp,
		},
	))
	empty := metadata.NewMetadataProperty(metadata.EmptyBlock())

	// ACT ====================================================================
	jsonMarshal, errMarsh := blockProp.MarshalJSON()
	unmarshErr := empty.UnmarshalJSON(jsonMarshal)
	resultingDataBlock := empty.Block()

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(11), empty.Code())
	assert.Len(t, resultingDataBlock.Mapping(), 13)
	assert.Equal(t, strProp, resultingDataBlock.Mapping()["STR"])
	assert.Equal(t, byteProp, resultingDataBlock.Mapping()["BYTE"])
	assert.Equal(t, byteProp2, resultingDataBlock.Mapping()["BYTE2"])
	assert.Equal(t, boolArrProp, resultingDataBlock.Mapping()["BOOL_ARR"])
	assert.Equal(t, float32Prop, resultingDataBlock.Mapping()["FLOAT32"])
	assert.Equal(t, float32ArrProp, resultingDataBlock.Mapping()["FLOAT32_ARR"])
	assert.Equal(t, boolProp, resultingDataBlock.Mapping()["BOOL_TRUE"])
	assert.Equal(t, boolProp2, resultingDataBlock.Mapping()["BOOL_FALSE"])
	assert.Equal(t, timeProp, resultingDataBlock.Mapping()["TIME"])
	assert.Equal(t, vecProp, resultingDataBlock.Mapping()["VEC2"])
	assert.Equal(t, vec3Prop, resultingDataBlock.Mapping()["VEC3"])
	assert.Equal(t, int32Prop, resultingDataBlock.Mapping()["INT32"])
	assert.Equal(t, binArrProp, resultingDataBlock.Mapping()["BIN_ARR"])
	assert.Equal(t, `{
	"BIN_ARR": Type 5 Array of 7 elements;
	"BOOL_ARR": Type 3 Array of 4 elements;
	"BOOL_FALSE": false;
	"BOOL_TRUE": true;
	"BYTE": 123;
	"BYTE2": 1;
	"FLOAT32": 1234.567749;
	"FLOAT32_ARR": Type 2 Array of 3 elements;
	"INT32": 888888;
	"STR": Meeee;
	"TIME": 1234567000000000 ns;
	"VEC2": 999.000000, 444.000000;
	"VEC3": 999.000000, 444.000000, 222.000000;
}`, empty.String())
}

func Test_MetadataArrayProperty_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	float32ArrProp := metadata.NewFloat32ArrayProperty([]float32{1.2, 3.4, 5.6})
	empty := metadata.NewFloat32ArrayProperty([]float32{})

	// ACT ====================================================================
	jsonMarshal, errMarsh := float32ArrProp.MarshalJSON()
	unmarshErr := empty.UnmarshalJSON(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(15), empty.Code())
	assert.Equal(t, float32ArrProp, empty)
}

func Test_MetadataVector3Property_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	vec3ArrProp := metadata.NewVector3Property(1.2, 3.4, 5.6)

	// ACT ====================================================================
	jsonMarshal, errMarsh := vec3ArrProp.MarshalJSON()
	unmarshVec, unmarshErr := metadata.UnmarshalNewVector3Property(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(7), unmarshVec.Code())
	assert.Equal(t, vec3ArrProp, unmarshVec)
	assert.Equal(t, "1.200000, 3.400000, 5.600000", unmarshVec.String())
}

func Test_MetadataVector2Property_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	vec3ArrProp := metadata.NewVector2Property(1.2, 3.4)

	// ACT ====================================================================
	jsonMarshal, errMarsh := vec3ArrProp.MarshalJSON()
	vec2Prop, unmarshErr := metadata.UnmarshalNewVector2Property(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(6), vec2Prop.Code())
	assert.Equal(t, vec3ArrProp, vec2Prop)
	assert.Equal(t, "1.200000, 3.400000", vec2Prop.String())
}

func Test_MetadataByteArrayProperty_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	binArrProp := metadata.NewBinaryArrayProperty([]byte{1, 2, 3, 4, 5})
	empty := metadata.NewBinaryArrayProperty([]byte{})

	// ACT ====================================================================
	jsonMarshal, errMarsh := binArrProp.MarshalJSON()
	unmarshErr := empty.UnmarshalJSON(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(18), empty.Code())
	assert.Equal(t, binArrProp, empty)
}

func Test_MetadataByteProperty_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	byteProp := metadata.NewByteProperty(33)

	// ACT ====================================================================
	jsonMarshal, errMarsh := byteProp.MarshalJSON()
	unmarshalByteProp, unmarshErr := metadata.UnmarshalNewByteProperty(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(5), unmarshalByteProp.Code())
	assert.Equal(t, byteProp, unmarshalByteProp)
	assert.Equal(t, "33", unmarshalByteProp.String())
}

func Test_MetadataBoolTrueProperty_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	byteProp := metadata.NewBoolProperty(true)

	// ACT ====================================================================
	jsonMarshal, errMarsh := byteProp.MarshalJSON()
	unmarshalByteProp, unmarshErr := metadata.UnmarshalNewBoolProperty(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(3), unmarshalByteProp.Code())
	assert.Equal(t, byteProp, unmarshalByteProp)
	assert.Equal(t, "true", unmarshalByteProp.String())
}

func Test_MetadataBoolFalseProperty_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	byteProp := metadata.NewBoolProperty(false)

	// ACT ====================================================================
	jsonMarshal, errMarsh := byteProp.MarshalJSON()
	unmarshalByteProp, unmarshErr := metadata.UnmarshalNewBoolProperty(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(4), unmarshalByteProp.Code())
	assert.Equal(t, byteProp, unmarshalByteProp)
	assert.Equal(t, "false", unmarshalByteProp.String())
}

func Test_Float32_MarshalJSON(t *testing.T) {
	// ARRANGE ================================================================
	floatProp := metadata.NewFloat32Property(556677.112233)

	// ACT ====================================================================
	jsonMarshal, errMarsh := floatProp.MarshalJSON()
	unmarshalByteProp, unmarshErr := metadata.UnmarshalNewFloat32Property(jsonMarshal)

	// ASSERT =================================================================
	assert.NoError(t, errMarsh)
	assert.NoError(t, unmarshErr)

	assert.Equal(t, byte(2), unmarshalByteProp.Code())
	assert.Equal(t, floatProp, unmarshalByteProp)
	assert.Equal(t, "556677.125000", unmarshalByteProp.String())
}
