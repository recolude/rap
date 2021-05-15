package enum_test

import (
	"fmt"
	"testing"

	"github.com/recolude/rap/format"
	enumCollection "github.com/recolude/rap/format/collection/enum"
	"github.com/recolude/rap/format/encoding/enum"
	"github.com/stretchr/testify/assert"
)

func Test_Singleenum(t *testing.T) {
	tests := map[string]struct {
		streamName  string
		enumMembers []string
		captures    []enumCollection.Capture
	}{
		"nil enums": {streamName: "", captures: nil},
		"0-enums":   {streamName: "empty stream", captures: []enumCollection.Capture{}},
		"1-enums": {
			streamName:  "ahhhh",
			enumMembers: []string{"a", "b"},
			captures: []enumCollection.Capture{
				enumCollection.NewCapture(1.2, 1),
			},
		},
	}

	storageTechniques := []struct {
		displayName    string
		technique      enumCollection.StorageTechnique
		timeTollerance float64
	}{
		{
			displayName:    "Raw64",
			technique:      enumCollection.Raw64,
			timeTollerance: 0,
		},
		{
			displayName:    "Raw32",
			technique:      enumCollection.Raw32,
			timeTollerance: 0.0005,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				collectionIn := enumCollection.NewCollection(tc.streamName, technique.technique, tc.enumMembers, tc.captures)

				encoder := enum.NewEncoder()

				// ACT ====================================================================
				header, collectionData, encodeErr := encoder.Encode([]format.CaptureCollection{collectionIn})
				collectionOut, decodeErr := encoder.Decode(header, collectionData[0])

				// ASSERT =================================================================
				assert.NoError(t, encodeErr)
				assert.NoError(t, decodeErr)
				assert.NotNil(t, collectionOut)
				assert.Len(t, collectionData, 1)
				assert.Equal(t, tc.streamName, collectionOut.Name())
				if assert.NotNil(t, collectionOut) {
					assert.Equal(t, collectionIn.Name(), collectionOut.Name())
					if assert.Len(t, collectionOut.Captures(), len(collectionIn.Captures())) {

						enmstr := collectionOut.(enumCollection.Collection)
						if assert.Len(t, enmstr.EnumMembers(), len(tc.enumMembers)) {
							for i, mem := range tc.enumMembers {
								assert.Equal(t, mem, enmstr.EnumMembers()[i])
							}
						}

						for i, c := range collectionOut.Captures() {
							enumCapture, ok := c.(enumCollection.Capture)
							if assert.True(t, ok) == false {
								break
							}

							assert.InDelta(t, tc.captures[i].Time(), enumCapture.Time(), technique.timeTollerance, "times are not equal: %.2f != %.2f", tc.captures[i].Time(), enumCapture.Time())
							assert.Equal(t, tc.captures[i].Value(), enumCapture.Value())

						}
					}
				}
			})
		}
	}
}
