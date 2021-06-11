package io

import (
	"io"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	"github.com/recolude/rap/format/encoding/enum"
	"github.com/recolude/rap/format/encoding/euler"
	"github.com/recolude/rap/format/encoding/event"
	"github.com/recolude/rap/format/encoding/position"
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
	return NewReader([]encoding.Encoder{
		event.NewEncoder(),
		position.NewEncoder(position.Oct48),
		euler.NewEncoder(euler.Raw32),
		enum.NewEncoder(),
	}, in).Read()
}
