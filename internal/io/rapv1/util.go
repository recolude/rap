package rapv1

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	fmt "fmt"
	"io"
	"io/ioutil"
	math "math"

	"github.com/golang/protobuf/proto"
	"github.com/recolude/rap/pkg/data"
	"github.com/recolude/rap/pkg/streams/enum"
	"github.com/recolude/rap/pkg/streams/euler"
	"github.com/recolude/rap/pkg/streams/event"
	"github.com/recolude/rap/pkg/streams/position"
)

func getNumberOfRecordings(file io.Reader) (int, error) {
	numberOfRecordings := make([]byte, 4)

	bytesRead, err := file.Read(numberOfRecordings)
	if err != nil {
		return -1, err
	}

	if bytesRead != 4 {
		return -1, fmt.Errorf("issue reading number of recordings, read %d bytes", bytesRead)
	}

	return int(binary.LittleEndian.Uint32(numberOfRecordings)), nil
}

func oldToNewEvents(oldEvents []*CustomEventCapture) event.Stream {
	customEventCaptures := make([]event.Capture, 0)
	for _, customEvent := range oldEvents {
		dictToUse := customEvent.GetData()

		// Older files did not have dictionaries associated with their
		// custom events
		if dictToUse == nil || len(dictToUse) == 0 {
			dictToUse = map[string]string{
				"value": customEvent.GetContents(),
			}
		}
		customEventCaptures = append(
			customEventCaptures,
			event.NewCapture(
				float64(customEvent.Time),
				customEvent.GetName(),
				dictToUse,
			),
		)
	}
	return event.NewStream("Custom Event", customEventCaptures)
}

func protobufToStd(inRec *Recording) (data.Recording, error) {
	subjectRecordings := make([]data.Recording, 0)
	for _, rec := range inRec.GetSubjects() {
		positionCaptures := make([]position.Capture, 0)
		rotationCaptures := make([]euler.Capture, 0)
		lifeCycleCaptures := make([]enum.Capture, 0)

		for _, pos := range rec.GetCapturedPositions() {
			positionCaptures = append(
				positionCaptures,
				position.NewCapture(
					float64(pos.GetTime()),
					float64(pos.GetX()),
					float64(pos.GetY()),
					float64(pos.GetZ()),
				),
			)
		}

		for _, rot := range rec.GetCapturedRotations() {
			rotationCaptures = append(
				rotationCaptures,
				euler.NewEulerZXYCapture(
					float64(rot.GetTime()),
					float64(rot.GetX()),
					float64(rot.GetY()),
					float64(rot.GetZ()),
				),
			)
		}

		for _, lifeEvent := range rec.GetLifecycleEvents() {
			lifeCycleCaptures = append(
				lifeCycleCaptures,
				enum.NewCapture(
					float64(lifeEvent.GetTime()),
					int(lifeEvent.Type),
				),
			)
		}

		positionStream := position.NewStream("Position", positionCaptures)
		rotationStream := euler.NewStream("Rotation", rotationCaptures)
		lifeStream := enum.NewStream("Life Cycle", []string{"START", "ENABLE", "DISABLE", "DESTROY"}, lifeCycleCaptures)

		subjectRecordings = append(subjectRecordings, &recordingV1{
			name: rec.GetName(),
			captureStreams: []data.CaptureStream{
				positionStream,
				rotationStream,
				oldToNewEvents(rec.GetCustomEvents()),
				lifeStream,
			},
			metadata: rec.GetMetadata(),
		})
	}

	return &recordingV1{
		name:           inRec.GetName(),
		captureStreams: []data.CaptureStream{oldToNewEvents(inRec.GetCustomEvents())},
		recordings:     subjectRecordings,
		metadata:       inRec.GetMetadata(),
	}, nil
}

func ReadRecording(file io.Reader) (data.Recording, error) {
	numberOfRecordings, err := getNumberOfRecordings(file)
	if err != nil {
		return nil, err
	}

	if numberOfRecordings != 1 {
		return nil, fmt.Errorf("Can only upload one recording at a time, recieved %d", numberOfRecordings)
	}

	recordingSize := make([]byte, 8)

	bytesRead, err := file.Read(recordingSize)
	if err != nil {
		return nil, err
	}

	if bytesRead != 8 {
		return nil, fmt.Errorf("Issue reading recording size, read %d bytes", bytesRead)
	}

	compressedSize := int64(binary.LittleEndian.Uint64(recordingSize))
	compressedBytes := make([]byte, compressedSize)
	bytesRead, err = file.Read(compressedBytes)

	if err != nil {
		return nil, err
	}

	if int64(bytesRead) != compressedSize {
		return nil, fmt.Errorf("Issue reading recording size, read %d bytes out of %d", bytesRead, compressedSize)
	}

	deflateReader := flate.NewReader(bytes.NewReader(compressedBytes))

	uncompresseRecording, err := ioutil.ReadAll(deflateReader)
	if err != nil {
		return nil, err
	}

	recording := &Recording{}

	err = proto.Unmarshal(uncompresseRecording, recording)
	if err != nil {
		return nil, err
	}

	return protobufToStd(recording)
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
