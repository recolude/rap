package format

import (
	"fmt"
)

type ValidateOption func(options *validateOptions)

type validateOptions struct {
	requireChronologicalCapture bool
}

func RequireChronologicalCapture(required bool) ValidateOption {
	return func(options *validateOptions) {
		options.requireChronologicalCapture = required
	}
}

// Validate ensures that a given recording meets all criteria specified
func Validate(rec Recording, options ...ValidateOption) error {
	finalOpts := &validateOptions{
		requireChronologicalCapture: true,
	}

	// Loop through each option
	for _, opt := range options {
		opt(finalOpts)
	}

	for _, child := range rec.Recordings() {
		err := Validate(child, options...)
		if err != nil {
			return err
		}
	}

	if finalOpts.requireChronologicalCapture {
		for _, col := range rec.CaptureCollections() {
			if col.Length() < 2 {
				continue
			}

			for i := 1; i < col.Length(); i++ {
				if col.CaptureAt(i).Time() < col.CaptureAt(i-1).Time() {
					return fmt.Errorf("[%s] %s: %s capture collection violates chronological event validator", rec.ID(), rec.Name(), col.Name())
				}
			}
		}
	}

	return nil
}
