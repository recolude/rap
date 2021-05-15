package encoding

import (
	"math"

	"github.com/recolude/rap/format"
)

func CollectionDuration(collection format.CaptureCollection) float64 {
	startingTime := math.Inf(1)
	endingTime := math.Inf(-1)

	for _, capture := range collection.Captures() {
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
