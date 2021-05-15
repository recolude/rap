package main

import (
	"fmt"
	"io"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/streams/enum"
	"github.com/recolude/rap/format/streams/euler"
	"github.com/recolude/rap/format/streams/event"
	"github.com/recolude/rap/format/streams/position"
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
	for _, stream := range recording.CaptureStreams() {
		switch v := stream.(type) {
		case event.Stream:
			curSummary.eventCaptureCount += len(v.Captures())
		case position.Stream:
			curSummary.positionCaptureCount += len(v.Captures())
		case enum.Stream:
			curSummary.enumCaptureCount += len(v.Captures())
		case euler.Stream:
			curSummary.eulerCaptureCount += len(v.Captures())
		default:
			curSummary.otherCaptureCount += len(stream.Captures())
		}
	}

	for _, rec := range recording.Recordings() {
		curSummary = curSummary.Combine(summarize(rec))
	}

	return curSummary
}

func printSummary(out io.Writer, recording format.Recording) {
	displayName := recording.Name()
	if displayName == "" {
		displayName = "[No Name]"
	}

	fmt.Fprintf(out, "Name: %s\n", displayName)
	fmt.Fprintf(out, "Sub Recordings: %d\n", len(recording.Recordings()))

	recSummary := summarize(recording)
	fmt.Fprintf(out, "Total Position Captures: %d\n", recSummary.positionCaptureCount)
	fmt.Fprintf(out, "Total Euler Captures:    %d\n", recSummary.eulerCaptureCount)
	fmt.Fprintf(out, "Total Event Captures:    %d\n", recSummary.eventCaptureCount)
	fmt.Fprintf(out, "Total Enum Captures:     %d\n", recSummary.enumCaptureCount)
	fmt.Fprintf(out, "Total Other Captures:    %d\n", recSummary.otherCaptureCount)
}
