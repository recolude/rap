package main

import (
	"fmt"
	"io"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/enum"
	"github.com/recolude/rap/format/collection/euler"
	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/collection/position"
)

type summary struct {
	positionCaptureCount int
	eventCaptureCount    int
	eulerCaptureCount    int
	enumCaptureCount     int
	otherCaptureCount    int
}

func (s summary) Combine(other summary) summary {
	return summary{
		positionCaptureCount: s.positionCaptureCount + other.positionCaptureCount,
		eventCaptureCount:    s.eventCaptureCount + other.eventCaptureCount,
		eulerCaptureCount:    s.eulerCaptureCount + other.eulerCaptureCount,
		enumCaptureCount:     s.enumCaptureCount + other.enumCaptureCount,
		otherCaptureCount:    s.otherCaptureCount + other.otherCaptureCount,
	}
}

func summarize(recording format.Recording) summary {
	curSummary := summary{}
	for _, collection := range recording.CaptureCollections() {
		switch v := collection.(type) {
		case event.Collection:
			curSummary.eventCaptureCount += v.Length()
		case position.Collection:
			curSummary.positionCaptureCount += v.Length()
		case enum.Collection:
			curSummary.enumCaptureCount += v.Length()
		case euler.Collection:
			curSummary.eulerCaptureCount += v.Length()
		default:
			curSummary.otherCaptureCount += collection.Length()
		}
	}

	for _, rec := range recording.Recordings() {
		curSummary = curSummary.Combine(summarize(rec))
	}

	return curSummary
}

func printSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%db", size)
	}

	editedSize := float64(size) / 1024.0
	if editedSize < 1024 {
		return fmt.Sprintf("%.2fkb", editedSize)
	}

	editedSize = editedSize / 1024.0
	return fmt.Sprintf("%.2fmb", editedSize)
}

func printSummary(out io.Writer, recording format.Recording, size int64) {
	displayName := recording.Name()
	if displayName == "" {
		displayName = "[No Name]"
	}

	fmt.Fprintf(out, "Name:                    %s\n", displayName)
	fmt.Fprintf(out, "File Size:               %s\n", printSize(size))
	fmt.Fprintf(out, "Duration:                %.2fs\n", format.RecordingDuration(recording))
	fmt.Fprintf(out, "Sub Recordings:          %d\n", len(recording.Recordings()))

	recSummary := summarize(recording)
	fmt.Fprintf(out, "Total Position Captures: %d\n", recSummary.positionCaptureCount)
	fmt.Fprintf(out, "Total Euler Captures:    %d\n", recSummary.eulerCaptureCount)
	fmt.Fprintf(out, "Total Event Captures:    %d\n", recSummary.eventCaptureCount)
	fmt.Fprintf(out, "Total Enum Captures:     %d\n", recSummary.enumCaptureCount)
	fmt.Fprintf(out, "Total Other Captures:    %d\n", recSummary.otherCaptureCount)
}
