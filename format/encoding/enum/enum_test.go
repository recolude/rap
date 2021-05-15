package enum_test

import (
	"fmt"
	"testing"

	"github.com/recolude/rap/format"
	enumStream "github.com/recolude/rap/format/collection/enum"
	"github.com/recolude/rap/format/encoding/enum"
	"github.com/stretchr/testify/assert"
)

func Test_Singleenum(t *testing.T) {
	tests := map[string]struct {
		streamName  string
		enumMembers []string
		captures    []enumStream.Capture
	}{
		"nil enums": {streamName: "", captures: nil},
		"0-enums":   {streamName: "empty stream", captures: []enumStream.Capture{}},
		"1-enums": {
			streamName:  "ahhhh",
			enumMembers: []string{"a", "b"},
			captures: []enumStream.Capture{
				enumStream.NewCapture(1.2, 1),
			},
		},
	}

	storageTechniques := []struct {
		displayName    string
		technique      enum.StorageTechnique
		timeTollerance float64
	}{
		{
			displayName:    "Raw64",
			technique:      enum.Raw64,
			timeTollerance: 0,
		},
		{
			displayName:    "Raw32",
			technique:      enum.Raw32,
			timeTollerance: 0.0005,
		},
	}

	for name, tc := range tests {
		for _, technique := range storageTechniques {
			t.Run(fmt.Sprintf("%s/%s", name, technique.displayName), func(t *testing.T) {
				streamIn := enumStream.NewStream(tc.streamName, tc.enumMembers, tc.captures)

				encoder := enum.NewEncoder(technique.technique)

				// ACT ====================================================================
				header, streamsData, encodeErr := encoder.Encode([]format.CaptureCollection{streamIn})
				streamOut, decodeErr := encoder.Decode(header, streamsData[0])

				// ASSERT =================================================================
				assert.NoError(t, encodeErr)
				assert.NoError(t, decodeErr)
				assert.NotNil(t, streamOut)
				assert.Len(t, streamsData, 1)
				assert.Equal(t, tc.streamName, streamOut.Name())
				if assert.NotNil(t, streamOut) {
					assert.Equal(t, streamIn.Name(), streamOut.Name())
					if assert.Len(t, streamOut.Captures(), len(streamIn.Captures())) {

						enmstr := streamOut.(enumStream.Stream)
						if assert.Len(t, enmstr.EnumMembers(), len(tc.enumMembers)) {
							for i, mem := range tc.enumMembers {
								assert.Equal(t, mem, enmstr.EnumMembers()[i])
							}
						}

						for i, c := range streamOut.Captures() {
							enumCapture, ok := c.(enumStream.Capture)
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
