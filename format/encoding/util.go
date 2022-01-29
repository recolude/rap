package encoding

import (
	"github.com/recolude/rap/format"
)

func CollectionDuration(collection format.CaptureCollection) float64 {
	return collection.End() - collection.Start()
}
