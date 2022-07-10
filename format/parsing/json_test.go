package parsing_test

import (
	"testing"

	"github.com/recolude/rap/format/parsing"

	"github.com/stretchr/testify/assert"
)

func Test_EmptyBytes_RaisesErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte{}

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Error(t, err)
	assert.Nil(t, recording)
}

func Test_InvalidJSON_RaisesErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{something}")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Error(t, err)
	assert.Nil(t, recording)
}

func Test_SingleJSONString_RaisesErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("\"something\"")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Error(t, err)
	assert.Nil(t, recording)
}

func Test_EmptyJSONObject_ThrowsErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{}")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recording object can not be empty")
	assert.Nil(t, recording)
}

func Test_JSONObj_ContainsOnlyID_ThrowsErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{ \"id\": \"something\" }")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recording requires name")
	assert.Nil(t, recording)
}

func Test_JSONObj_IDNotString_ThrowsErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{ \"id\": 5 }")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recording id must be string")
	assert.Nil(t, recording)
}

func Test_JSONObj_NameNotString_ThrowsErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{ \"id\": \"5\", \"name\": 5 }")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recording name must be string")
	assert.Nil(t, recording)
}

func Test_JSONObj_ContainsOnlyName_ThrowsErr(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{ \"name\": \"something\" }")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recording requires id")
	assert.Nil(t, recording)
}

func Test_JSONObj_JustNameAndID_Success(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte("{ \"id\": \"my id\", \"name\": \"my name\"  }")

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Nil(t, err)
	assert.NotNil(t, recording)

	assert.Equal(t, "my id", recording.ID())
	assert.Equal(t, "my name", recording.Name())
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))
	assert.Equal(t, 0, len(recording.Recordings()))
}

func Test_JSONObj_EmptySubRecordings_Success(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": []
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Nil(t, err)
	assert.NotNil(t, recording)

	assert.Equal(t, "my id", recording.ID())
	assert.Equal(t, "my name", recording.Name())
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))
	assert.Equal(t, 0, len(recording.Recordings()))
}

func Test_JSONObj_SingleSubRecording_Success(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": [
			{
				"id": "sub",
				"name": "subRecording"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Nil(t, err)
	assert.NotNil(t, recording)

	assert.Equal(t, "my id", recording.ID())
	assert.Equal(t, "my name", recording.Name())
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))

	if assert.Equal(t, 1, len(recording.Recordings())) {
		subRecording := recording.Recordings()[0]
		assert.NotNil(t, subRecording)
		assert.Equal(t, "sub", subRecording.ID())
		assert.Equal(t, "subRecording", subRecording.Name())
		assert.Equal(t, 0, len(subRecording.Metadata().Mapping()))
		assert.Equal(t, 0, len(subRecording.Binaries()))
		assert.Equal(t, 0, len(subRecording.BinaryReferences()))
		assert.Equal(t, 0, len(subRecording.CaptureCollections()))
	}
}

func Test_JSONObj_InvalidSubRecordingDefinition_Err(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": [
			{
				"name": "subRecording"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recording requires id")
	assert.Nil(t, recording)
}

func Test_JSONObj_SubRecordingsMustBeArray_Err(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": 9
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recordings property should be array")
	assert.Nil(t, recording)
}

func Test_JSONObj_SubRecordingsMustBeArray2_Err(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": {
			
		}
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recordings property should be array")
	assert.Nil(t, recording)
}

func Test_JSONObj_SubRecordingsMustBeArrayOfObjects_Err(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": [
			9	
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "recordings property should be array of recording objects")
	assert.Nil(t, recording)
}

func Test_JSONObj_MultipleSubRecording_Success(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"recordings": [
			{
				"id": "sub",
				"name": "subRecording"
			},
			{
				"id": "sub2",
				"name": "subRecording 2"
			},
			{
				"id": "sub3",
				"name": "subRecording 3"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Nil(t, err)
	assert.NotNil(t, recording)

	assert.Equal(t, "my id", recording.ID())
	assert.Equal(t, "my name", recording.Name())
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))

	if assert.Equal(t, 3, len(recording.Recordings())) {
		subRecording1 := recording.Recordings()[0]
		assert.NotNil(t, subRecording1)
		assert.Equal(t, "sub", subRecording1.ID())
		assert.Equal(t, "subRecording", subRecording1.Name())
		assert.Equal(t, 0, len(subRecording1.Metadata().Mapping()))
		assert.Equal(t, 0, len(subRecording1.Binaries()))
		assert.Equal(t, 0, len(subRecording1.BinaryReferences()))
		assert.Equal(t, 0, len(subRecording1.CaptureCollections()))

		subRecording2 := recording.Recordings()[1]
		assert.NotNil(t, subRecording2)
		assert.Equal(t, "sub2", subRecording2.ID())
		assert.Equal(t, "subRecording 2", subRecording2.Name())
		assert.Equal(t, 0, len(subRecording2.Metadata().Mapping()))
		assert.Equal(t, 0, len(subRecording2.Binaries()))
		assert.Equal(t, 0, len(subRecording2.BinaryReferences()))
		assert.Equal(t, 0, len(subRecording2.CaptureCollections()))

		subRecording3 := recording.Recordings()[2]
		assert.NotNil(t, subRecording3)
		assert.Equal(t, "sub3", subRecording3.ID())
		assert.Equal(t, "subRecording 3", subRecording3.Name())
		assert.Equal(t, 0, len(subRecording3.Metadata().Mapping()))
		assert.Equal(t, 0, len(subRecording3.Binaries()))
		assert.Equal(t, 0, len(subRecording3.BinaryReferences()))
		assert.Equal(t, 0, len(subRecording3.CaptureCollections()))
	}
}

func Test_JSONObj_MultiplePropertyMetadata(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"metadata": {
			"float": 3.2,
			"something": "weellp"
		}
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Nil(t, err)
	assert.NotNil(t, recording)

	assert.Equal(t, "my id", recording.ID())
	assert.Equal(t, "my name", recording.Name())

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))
	assert.Equal(t, 0, len(recording.Recordings()))

	assert.Equal(t, 2, len(recording.Metadata().Mapping()))
	assert.Equal(t, "weellp", recording.Metadata().Mapping()["something"].String())
	assert.Equal(t, "3.200000", recording.Metadata().Mapping()["float"].String())
}

func Test_JSONObj_MetadataNotObject_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"metadata": 6
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "metadata should be object")
	assert.Nil(t, recording)
}
