package format

import "math"

func RecordingStart(rec Recording) float64 {
	min := math.Inf(0)

	for _, rec := range rec.Recordings() {
		childMin := RecordingStart(rec)
		if childMin < min {
			min = childMin
		}
	}

	for _, coll := range rec.CaptureCollections() {
		if coll.Length() > 0 {
			cap := coll.CaptureAt(0)
			if cap.Time() < min {
				min = cap.Time()
			}
		}
	}

	return min
}

func RecordingEnd(rec Recording) float64 {
	max := math.Inf(-1)

	for _, rec := range rec.Recordings() {
		childMax := RecordingEnd(rec)
		if childMax > max {
			max = childMax
		}
	}

	for _, coll := range rec.CaptureCollections() {
		if coll.Length() > 0 {
			cap := coll.CaptureAt(coll.Length() - 1)
			if cap.Time() > max {
				max = cap.Time()
			}
		}
	}

	return max
}

func RecordingDuration(rec Recording) float64 {
	min := RecordingStart(rec)
	max := RecordingEnd(rec)
	return max - min
}
