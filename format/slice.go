package format

import "math"

type SliceOption func(options *sliceOptions)

type sliceOptions struct {
	beginning    float64
	end          float64
	keepBinaries bool
}

func BeginningOfSlice(beginning float64) SliceOption {
	return func(options *sliceOptions) {
		options.beginning = beginning
	}
}

func EndOfSlice(end float64) SliceOption {
	return func(options *sliceOptions) {
		options.end = end
	}
}

func KeepBinariesInSlice(shouldKeep bool) SliceOption {
	return func(options *sliceOptions) {
		options.keepBinaries = shouldKeep
	}
}

// Slice takes a recording and keeps information that satisfies the slice
// options provided.
func Slice(rec Recording, options ...SliceOption) Recording {
	finalOpts := &sliceOptions{
		beginning:    math.Inf(-1),
		end:          math.Inf(1),
		keepBinaries: true,
	}

	// Loop through each option
	for _, opt := range options {
		opt(finalOpts)
	}

	outRec := recording{
		id:               rec.ID(),
		name:             rec.Name(),
		metadata:         rec.Metadata(),
		binaryReferences: rec.BinaryReferences(),
	}

	if finalOpts.keepBinaries {
		outRec.binaries = rec.Binaries()
	}

	allChildRec := make([]Recording, len(rec.Recordings()))
	for i, child := range rec.Recordings() {
		allChildRec[i] = Slice(child, options...)
	}
	outRec.recordings = allChildRec

	allCollections := make([]CaptureCollection, len(rec.CaptureCollections()))
	for i, child := range rec.CaptureCollections() {
		allCollections[i] = child.Slice(finalOpts.beginning, finalOpts.end)
	}
	outRec.captureCollections = allCollections

	return outRec
}

// CaptureFallsWithin returns true whenever the captures time falls within the
// range: beginning <= time < end
func CaptureFallsWithin(c Capture, beginning, end float64) bool {
	return c.Time() >= beginning && c.Time() < end
}
