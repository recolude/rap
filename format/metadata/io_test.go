package metadata_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/EliCDavis/vector"
	"github.com/recolude/rap/format/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_BasicIO(t *testing.T) {
	// ARRANGE ================================================================
	tests := map[string]metadata.Property{
		"int prop 77":   metadata.NewIntProperty(77),
		"int prop -100": metadata.NewIntProperty(100),
		"int prop -0":   metadata.NewIntProperty(0),
		"string prop":   metadata.NewStringProperty("dee"),
		"bool true":     metadata.NewBoolProperty(true),
		"bool false":    metadata.NewBoolProperty(false),
		"byte test":     metadata.NewByteProperty(22),
		"vec2 test":     metadata.NewVector2Property(1.2, 3.4),
		"vec3 test":     metadata.NewVector3Property(1.2, 3.4, 5.6),
		"block test": metadata.NewMetadataProperty(metadata.NewBlock(
			map[string]metadata.Property{
				"nested prop 1":    metadata.NewStringProperty("God kill me"),
				"nested prop 2":    metadata.NewStringProperty("ahhhh"),
				"nested prop 3":    metadata.NewIntProperty(666),
				"nested prop time": metadata.NewTimeProperty(time.Now()),
			},
		)),
		"String Array": metadata.NewStringArrayProperty([]string{"x", "y", "z"}),
		"Int Array":    metadata.NewIntArrayProperty([]int{1, 2, 3, 4}),
		"time array":   metadata.NewTimestampArrayProperty([]time.Time{time.Now(), time.Now().Add(1)}),
		"float array":  metadata.NewFloat32ArrayProperty([]float32{1.2, 3.4}),
		"vec2 array":   metadata.NewVector2ArrayProperty([]vector.Vector2{vector.NewVector2(1, 2), vector.NewVector2(3, 4)}),
		"vec3 array":   metadata.NewVector3ArrayProperty([]vector.Vector3{vector.NewVector3(1, 2, 3), vector.NewVector3(4, 5, 6)}),
		"metadata array": metadata.NewMetadataArrayProperty([]metadata.Block{
			metadata.EmptyBlock(),
			metadata.NewBlock(map[string]metadata.Property{
				"ahh": metadata.NewBoolProperty(true),
			}),
			metadata.EmptyBlock(),
		}),
		"byte array": metadata.NewBinaryArrayProperty([]byte{1, 2, 3, 4}),
		"bool array": metadata.NewBoolArrayProperty([]bool{true, false, true}),
	}

	// ACT ====================================================================
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bufferData := bytes.Buffer{}

			_, err := metadata.WriteProprty(&bufferData, tc)
			assert.NoError(t, err)
			propBack, err := metadata.ReadProperty(bytes.NewReader(bufferData.Bytes()))
			assert.NoError(t, err)

			assert.Equal(t, tc, propBack)
		})
	}
}
