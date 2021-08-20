package io_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/recolude/rap/format"
	rapio "github.com/recolude/rap/format/io"
)

var result format.Recording

func BenchmarkLoadSimple(b *testing.B) {
	var r format.Recording
	for i := 0; i < b.N; i++ {
		file, err := os.Open(filepath.Join(v1DirectoryTestData, "Demo 38subj v1.rap")) // 1000 Particles over 300 seconds events.rap
		if err != nil {
			panic(err)
		}
		r, _, err = rapio.Load(file)
		if err != nil {
			panic(err)
		}
	}
	result = r
}
