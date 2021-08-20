package rapv1

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	fmt "fmt"
	"io"
	"io/ioutil"
	math "math"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/enum"
	"github.com/recolude/rap/format/collection/euler"
	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/metadata"
)

func getNumberOfRecordings(file io.Reader) (int, int, error) {
	numberOfRecordings := make([]byte, 4)

	bytesRead, err := file.Read(numberOfRecordings)
	if err != nil {
		return -1, bytesRead, err
	}

	if bytesRead != 4 {
		return -1, bytesRead, fmt.Errorf("issue reading number of recordings, read %d bytes", bytesRead)
	}

	return int(binary.LittleEndian.Uint32(numberOfRecordings)), bytesRead, nil
}

func oldToNewEvents(oldEvents []*CustomEventCapture) event.Collection {
	customEventCaptures := make([]event.Capture, len(oldEvents))
	for eventIndex, customEvent := range oldEvents {
		eventDict := customEvent.GetData()
		dictToUse := make(map[string]metadata.Property)
		// Older files did not have dictionaries associated with their
		// custom events
		if eventDict == nil || len(eventDict) == 0 {
			dictToUse["value"] = metadata.NewStringProperty(customEvent.GetContents())
		} else {
			for key, val := range eventDict {
				floatVal, err := strconv.ParseFloat(val, 32)
				if err == nil {
					dictToUse[key] = metadata.NewFloat32Property(float32(floatVal))
				} else {
					dictToUse[key] = metadata.NewStringProperty(val)
				}
				dictToUse[key] = metadata.NewStringProperty(val)
			}
		}

		customEventCaptures[eventIndex] = event.NewCapture(
			float64(customEvent.Time),
			customEvent.GetName(),
			metadata.NewBlock(dictToUse),
		)
	}
	return event.NewCollection("Custom Event", customEventCaptures)
}

func convertMetadata(original map[string]string) metadata.Block {
	out := make(map[string]metadata.Property)

	for k, v := range original {
		floatVal, err := strconv.ParseFloat(v, 32)
		if err == nil {
			out[k] = metadata.NewFloat32Property(float32(floatVal))
		} else {
			out[k] = metadata.NewStringProperty(v)
		}
	}

	return metadata.NewBlock(out)
}

func protobufToStd(inRec *Recording) (format.Recording, error) {
	subjectRecordings := make([]format.Recording, len(inRec.GetSubjects()))
	lifecycleEnumMembers := []string{"START", "ENABLE", "DISABLE", "DESTROY"}
	for subjectIndex, rec := range inRec.GetSubjects() {
		positionCaptures := make([]position.Capture, len(rec.GetCapturedPositions()))
		rotationCaptures := make([]euler.Capture, len(rec.GetCapturedRotations()))
		lifeCycleCaptures := make([]enum.Capture, len(rec.GetLifecycleEvents()))

		for posIndex, pos := range rec.GetCapturedPositions() {
			positionCaptures[posIndex] = position.NewCapture(
				float64(pos.GetTime()),
				float64(pos.GetX()),
				float64(pos.GetY()),
				float64(pos.GetZ()),
			)
		}

		for rotIndex, rot := range rec.GetCapturedRotations() {
			rotationCaptures[rotIndex] = euler.NewEulerZXYCapture(
				float64(rot.GetTime()),
				float64(rot.GetX()),
				float64(rot.GetY()),
				float64(rot.GetZ()),
			)
		}

		for lifeIndex, lifeEvent := range rec.GetLifecycleEvents() {
			lifeCycleCaptures[lifeIndex] = enum.NewCapture(
				float64(lifeEvent.GetTime()),
				int(lifeEvent.Type),
			)
		}

		positionStream := position.NewCollection("Position", positionCaptures)
		rotationStream := euler.NewCollection("Rotation", rotationCaptures)
		lifeStream := enum.NewCollection("Life Cycle", lifecycleEnumMembers, lifeCycleCaptures)

		subjectRecordings[subjectIndex] = &recordingV1{
			id:   fmt.Sprint(rec.GetId()),
			name: rec.GetName(),
			captureStreams: []format.CaptureCollection{
				positionStream,
				rotationStream,
				oldToNewEvents(rec.GetCustomEvents()),
				lifeStream,
			},
			metadata: convertMetadata(rec.GetMetadata()),
		}
	}

	return &recordingV1{
		id:             inRec.GetName(),
		name:           inRec.GetName(),
		captureStreams: []format.CaptureCollection{oldToNewEvents(inRec.GetCustomEvents())},
		recordings:     subjectRecordings,
		metadata:       convertMetadata(inRec.GetMetadata()),
	}, nil
}

func ReadRecording(file io.Reader) (format.Recording, int, error) {
	numberOfRecordings, bytesReadNumberRec, err := getNumberOfRecordings(file)
	if err != nil {
		return nil, bytesReadNumberRec, err
	}

	if numberOfRecordings != 1 {
		return nil, bytesReadNumberRec, fmt.Errorf("Can only upload one recording at a time, received %d", numberOfRecordings)
	}

	recordingSize := make([]byte, 8)

	bytesRead := bytesReadNumberRec
	bytesReadFileSize, err := file.Read(recordingSize)
	bytesRead += bytesReadFileSize
	if err != nil {
		return nil, bytesRead, err
	}

	if bytesReadFileSize != 8 {
		return nil, bytesRead, fmt.Errorf("Issue reading recording size, read %d bytes", bytesRead)
	}

	compressedSize := int64(binary.LittleEndian.Uint64(recordingSize))
	compressedBytes := make([]byte, compressedSize)
	compressedBytesRead, err := file.Read(compressedBytes)
	bytesRead += compressedBytesRead
	if err != nil {
		return nil, bytesRead, err
	}

	if int64(compressedBytesRead) != compressedSize {
		return nil, bytesRead, fmt.Errorf("Issue reading recording size, read %d bytes out of %d", bytesRead, compressedSize)
	}

	deflateReader := flate.NewReader(bytes.NewReader(compressedBytes))

	uncompresseRecording, err := ioutil.ReadAll(deflateReader)
	if err != nil {
		return nil, bytesRead, err
	}

	recording := &Recording{}

	err = proto.Unmarshal(uncompresseRecording, recording)
	if err != nil {
		return nil, bytesRead, err
	}

	rec, err := protobufToStd(recording)

	return rec, bytesRead, err
}

func GetStartOfRecording(recording Recording) float64 {
	min := math.Inf(1)

	for _, e := range recording.CustomEvents {
		if float64(e.Time) < min {
			min = float64(e.Time)
		}
	}

	for _, subj := range recording.Subjects {
		for _, e := range subj.CustomEvents {
			if float64(e.Time) < min {
				min = float64(e.Time)
			}
		}

		for _, e := range subj.CapturedPositions {
			if float64(e.Time) < min {
				min = float64(e.Time)
			}
		}

		for _, e := range subj.CapturedRotations {
			if float64(e.Time) < min {
				min = float64(e.Time)
			}
		}

		for _, e := range subj.LifecycleEvents {
			if float64(e.Time) < min {
				min = float64(e.Time)
			}
		}
	}

	return min
}

func GetEndOfRecording(recording Recording) float64 {
	max := math.Inf(-1)

	for _, e := range recording.CustomEvents {
		if float64(e.Time) > max {
			max = float64(e.Time)
		}
	}

	for _, subj := range recording.Subjects {
		for _, e := range subj.CustomEvents {
			if float64(e.Time) > max {
				max = float64(e.Time)
			}
		}

		for _, e := range subj.CapturedPositions {
			if float64(e.Time) > max {
				max = float64(e.Time)
			}
		}

		for _, e := range subj.CapturedRotations {
			if float64(e.Time) > max {
				max = float64(e.Time)
			}
		}

		for _, e := range subj.LifecycleEvents {
			if float64(e.Time) > max {
				max = float64(e.Time)
			}
		}
	}

	return max
}

func GetRecordingDuration(recording Recording) float64 {
	return GetEndOfRecording(recording) - GetStartOfRecording(recording)
}

func CountRecordingCustomEvenets(recording Recording) int {
	events := len(recording.CustomEvents)
	for _, subj := range recording.Subjects {
		events += len(subj.CustomEvents)
	}
	return events
}
