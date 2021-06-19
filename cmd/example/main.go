package main

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/encoding"
	positionEncoder "github.com/recolude/rap/format/encoding/position"
	"github.com/recolude/rap/format/io"
	"github.com/recolude/rap/format/metadata"
)

func main() {
	iterations := 1000
	positions := make([]position.Capture, iterations)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		currentTime := float64(i)
		positions[i] = position.NewCapture(currentTime, 0, math.Sin(currentTime/10.0), 0)
	}
	duration := time.Since(start)

	rec := format.NewRecording(
		"",
		"Sin Wave Demo",
		[]format.CaptureCollection{
			position.NewCollection("Sin Wave", positions),
		},
		nil,
		metadata.NewBlock(map[string]metadata.Property{
			"iterations": metadata.NewIntProperty(iterations),
			"benchmark":  metadata.NewStringProperty(duration.String()),
		}),
		nil,
		nil,
	)

	log.Print(duration.String())

	f, _ := os.Create("sin demo.rap")
	recordingWriter := io.NewWriter(
		[]encoding.Encoder{
			positionEncoder.NewEncoder(positionEncoder.Oct24),
		},
		true,
		f,
		io.BST16,
	)

	// Writes a recording in 1,171 bytes
	recordingWriter.Write(rec)
}
