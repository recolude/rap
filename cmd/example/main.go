package main

import (
	"math"
	"os"
	"time"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
	rapio "github.com/recolude/rap/format/io"
)

func main() {
	iterations := 1000
	positions := make([]position.Capture, iterations)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		positions[i] = position.NewCapture(float64(i), 0, math.Sin(float64(i)), 0)
	}
	duration := time.Since(start)

	sinWavePositions := position.NewCollection("Sin Wave", positions)

	rec := format.NewRecording(
		"",
		"Sin Wave Demo",
		[]format.CaptureCollection{sinWavePositions},
		nil,
		format.NewMetadataBlock(map[string]format.Property{
			"iterations": format.NewIntProperty(int32(iterations)),
			"benchmark":  format.NewStringProperty(duration.String()),
		}),
		nil,
		nil,
	)

	f, _ := os.Create("sin demo.rap")
	recordingWriter := rapio.NewRecoludeWriter(f)
	recordingWriter.Write(rec)
}
