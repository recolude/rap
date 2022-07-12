package parsing_test

import (
	"testing"

	"github.com/recolude/rap/format/collection/enum"
	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
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

func Test_JSONObj_ReferencesNotArray_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": 6
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "references property should be an array")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferencesNotArrayPt2_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": {}
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "references property should be an array")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferencesArrayOfNonObjects_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [6]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "references property should be array of reference objects")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferenceObjWithoutName_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"uri": "something"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "reference requires name")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferenceObjWithoutURI_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"name": "something"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "reference requires uri")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferenceObjWithoutSize_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"name": "something",
				"uri": "something"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "reference requires a size property")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferenceObjFloatSize_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"name": "something",
				"uri": "something",
				"size": 1.2
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "reference size must be int")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferenceObjNegativeSize_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"name": "something",
				"uri": "something",
				"size": -20
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "reference size must be non-negative")
	assert.Nil(t, recording)
}

func Test_JSONObj_ReferenceObjMetadataNotObj_Errs(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"name": "something",
				"uri": "uri",
				"metadata": 7
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "metadata should be object")
	assert.Nil(t, recording)
}

func Test_JSONObj_SingleReferenceNoMetadata(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"references": [
			{
				"name": "ref 1",
				"uri": "file:///C:/dev/yo/mama.fat",
				"size": 200
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

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))
	assert.Equal(t, 0, len(recording.Recordings()))
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))

	assert.Equal(t, 1, len(recording.BinaryReferences()))
	assert.Equal(t, "ref 1", recording.BinaryReferences()[0].Name())
	assert.Equal(t, uint64(200), recording.BinaryReferences()[0].Size())
	assert.Equal(t, "file:///C:/dev/yo/mama.fat", recording.BinaryReferences()[0].URI())
}

func Test_JSONObj_EmptyCollectionsArray(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.Nil(t, err)
	assert.NotNil(t, recording)

	assert.Equal(t, "my id", recording.ID())
	assert.Equal(t, "my name", recording.Name())

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.CaptureCollections()))
	assert.Equal(t, 0, len(recording.Recordings()))
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))
}

func Test_JSONObj_PositionCollectionCaptures(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"type": "recolude.position",
				"name": "Some Positions",
				"captures": [
					{
						"time": 1.3,
						"data": {
							"x": 1.1,
							"y": 2.2,
							"z": 3.3
						}
					},
					{
						"time": 2.4,
						"data": {
							"x": 4.4,
							"y": 5.5,
							"z": 6.6
						}
					}
				]
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

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.Recordings()))
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))

	assert.Equal(t, 1, len(recording.CaptureCollections()))
	assert.Equal(t, "Some Positions", recording.CaptureCollections()[0].Name())
	assert.Equal(t, "recolude.position", recording.CaptureCollections()[0].Signature())
	assert.Equal(t, 2, len(recording.CaptureCollections()[0].Captures()))
	assert.Equal(t, 1.3, recording.CaptureCollections()[0].Captures()[0].Time())
	assert.Equal(t, 2.4, recording.CaptureCollections()[0].Captures()[1].Time())

	assert.Equal(t, "[1.30] - 1.10, 2.20, 3.30", recording.CaptureCollections()[0].Captures()[0].String())
	assert.Equal(t, "[2.40] - 4.40, 5.50, 6.60", recording.CaptureCollections()[0].Captures()[1].String())
}

func Test_JSONObj_PositionCollectionCaptureLackingZ_Errors(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"type": "recolude.position",
				"name": "Some Positions",
				"captures": [
					{
						"time": 1.3,
						"data": {
							"x": 1.1,
							"y": 2.2
						}
					},
					{
						"time": 2.4,
						"data": {
							"x": 4.4,
							"y": 5.5,
							"z": 6.6
						}
					}
				]
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "position capture requires z property")
	assert.Nil(t, recording)
}

func Test_JSONObj_CollectionWithoutName_Errors(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"type": "recolude.position",
				"captures": [
					
				]
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "collection requires name")
	assert.Nil(t, recording)
}

func Test_JSONObj_CollectionWithoutType_Errors(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"name": "recolude.position",
				"captures": [
					
				]
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "collection requires type")
	assert.Nil(t, recording)
}

func Test_JSONObj_CollectionWithoutCaptures_Errors(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"name": "recolude.position",
				"type": "recolude.position"
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "collection object requires captures property")
	assert.Nil(t, recording)
}

func Test_JSONObj_CollectionCapturesAsNonArray_Errors(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"name": "recolude.position",
				"type": "recolude.position",
				"captures": {}
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "collection's captures property must be an array")
	assert.Nil(t, recording)
}

func Test_JSONObj_InvalidCollectionType_Errors(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"name": "recolude.position",
				"type": "recolude.unknownType",
				"captures": []
			}
		]
	}`)

	// ACT ====================================================================
	recording, err := parsing.FromJSON(payload)

	// ASSERT =================================================================
	assert.EqualError(t, err, "unrecognized collection type: 'recolude.unknownType'")
	assert.Nil(t, recording)
}

func Test_JSONObj_RotationCollectionCaptures(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"type": "recolude.euler",
				"name": "Some Rotations",
				"captures": [
					{
						"time": 1.3,
						"data": {
							"x": 1.1,
							"y": 2.2,
							"z": 3.3
						}
					},
					{
						"time": 2.4,
						"data": {
							"x": 4.4,
							"y": 5.5,
							"z": 6.6
						}
					}
				]
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

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.Recordings()))
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))

	assert.Equal(t, 1, len(recording.CaptureCollections()))
	assert.Equal(t, "Some Rotations", recording.CaptureCollections()[0].Name())
	assert.Equal(t, "recolude.euler", recording.CaptureCollections()[0].Signature())
	assert.Equal(t, 2, len(recording.CaptureCollections()[0].Captures()))
	assert.Equal(t, 1.3, recording.CaptureCollections()[0].Captures()[0].Time())
	assert.Equal(t, 2.4, recording.CaptureCollections()[0].Captures()[1].Time())

	assert.Equal(t, "[1.30] Rotation - 1.10, 2.20, 3.30", recording.CaptureCollections()[0].Captures()[0].String())
	assert.Equal(t, "[2.40] Rotation - 4.40, 5.50, 6.60", recording.CaptureCollections()[0].Captures()[1].String())
}

func Test_JSONObj_EventCollectionCaptures(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"type": "recolude.event",
				"name": "Some Custom EVents",
				"captures": [
					{
						"time": 1.3,
						"data": {
							"name": "Some event",
							"metadata": {
								"some key": 12
							}
						}
					}
				]
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

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.Recordings()))
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))

	assert.Equal(t, 1, len(recording.CaptureCollections()))
	assert.Equal(t, "Some Custom EVents", recording.CaptureCollections()[0].Name())
	assert.Equal(t, "recolude.event", recording.CaptureCollections()[0].Signature())
	assert.Equal(t, 1, len(recording.CaptureCollections()[0].Captures()))
	assert.Equal(t, 1.3, recording.CaptureCollections()[0].Captures()[0].Time())

	event, isEvent := recording.CaptureCollections()[0].Captures()[0].(event.Capture)
	assert.True(t, isEvent)

	assert.Equal(t, "Some event", event.Name())
	assert.Len(t, event.Metadata().Mapping(), 1)
	assert.Equal(
		t,
		metadata.NewIntProperty(12),
		event.Metadata().Mapping()["some key"],
	)
}

func Test_JSONObj_EnumCollectionCaptures(t *testing.T) {
	// ARRANGE ================================================================
	payload := []byte(`{ 
		"id": "my id", 
		"name": "my name",
		"collections": [
			{
				"type": "recolude.enum",
				"name": "Some Enums",
				"captures": [
					{
						"time": 1.3,
						"data": "Test"
					}
				]
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

	assert.Equal(t, 0, len(recording.Binaries()))
	assert.Equal(t, 0, len(recording.Recordings()))
	assert.Equal(t, 0, len(recording.Metadata().Mapping()))
	assert.Equal(t, 0, len(recording.BinaryReferences()))

	assert.Equal(t, 1, len(recording.CaptureCollections()))
	assert.Equal(t, "Some Enums", recording.CaptureCollections()[0].Name())
	assert.Equal(t, "recolude.enum", recording.CaptureCollections()[0].Signature())
	assert.Equal(t, 1, len(recording.CaptureCollections()[0].Captures()))
	assert.Equal(t, 1.3, recording.CaptureCollections()[0].Captures()[0].Time())

	capture, isEnum := recording.CaptureCollections()[0].Captures()[0].(enum.Capture)
	assert.True(t, isEnum)
	assert.Equal(t, 0, capture.Value())
}
