package encoding

import (
	"math"

	"github.com/recolude/rap/pkg/data"
)

func StreamDuration(stream data.CaptureStream) float64 {
	startingTime := math.Inf(1)
	endingTime := math.Inf(-1)

	for _, capture := range stream.Captures() {
		if capture.Time() < startingTime {
			startingTime = capture.Time()
		}
		if capture.Time() > endingTime {
			endingTime = capture.Time()
		}
	}

	duration := endingTime - startingTime

	return duration
}
