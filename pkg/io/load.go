package io

import (
	"io"

	"github.com/recolude/rap/pkg/data"
)

func GetRecoringVersion(file io.Reader) (int, int, error) {
	version := make([]byte, 1)

	bytesRead, err := file.Read(version)
	if err != nil {
		return 0, bytesRead, err
	}

	return int(version[0]), bytesRead, nil
}

func Load(in io.Reader) (data.Recording, int, error) {
	return NewReader(nil, in).Read()
}
