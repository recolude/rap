package io

import (
	"io"

	"github.com/recolude/rap/format"
)

func GetRecoringVersion(file io.Reader) (int, int, error) {
	version := make([]byte, 1)

	bytesRead, err := file.Read(version)
	if err != nil {
		return 0, bytesRead, err
	}

	return int(version[0]), bytesRead, nil
}

func Load(in io.Reader) (format.Recording, int, error) {
	return NewReader(nil, in).Read()
}
