package io

import (
	"fmt"
	"io"

	"github.com/recolude/rap/internal/io/rapv1"
	"github.com/recolude/rap/pkg/data"
)

func GetRecoringVersion(file io.Reader) (int, error) {
	version := make([]byte, 1)

	_, err := file.Read(version)
	if err != nil {
		return 0, err
	}

	return int(version[0]), nil
}

// func GetRecordingInformation(file io.Reader) (*Recording, error) {
// 	version, err := GetRecoringVersion(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if version != 1 {
// 		return nil, fmt.Errorf("unsupported version: %d", version)
// 	}

// 	numberOfRecordings, err := GetNumberOfRecordings(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if numberOfRecordings != 1 {
// 		return nil, errors.New("Can only upload one recording at a time, recieved " + string(numberOfRecordings))
// 	}

// 	return readRecording(file)
// }

func Load(in io.Reader) (data.Recording, error) {
	if in == nil {
		panic("Attempting to load recording from nil reader")
	}

	version, err := GetRecoringVersion(in)
	if err != nil {
		return nil, err
	}

	if version == 1 {
		return rapv1.ReadRecording(in)
	}

	return nil, fmt.Errorf("Unrecognized file version: %d", version)
}
