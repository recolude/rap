package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
)

func RecordingFromCSV(in io.Reader) (format.Recording, error) {
	csvReader := csv.NewReader(in)

	nameIndex := -1
	idIndex := -1
	timeIndex := -1

	posXIndex := -1
	posYIndex := -1
	posZIndex := -1

	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	for i, column := range header {
		switch strings.TrimSpace(column) {
		case "time":
			timeIndex = i
			break

		case "id":
			idIndex = i
			break

		case "name":
			nameIndex = i
			break

		case "x":
			posXIndex = i
			break

		case "y":
			posYIndex = i
			break

		case "z":
			posZIndex = i
			break
		}
	}

	workingData := make(map[string]map[string][]position.Capture)

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		timeDirty := row[timeIndex]
		time, err := strconv.ParseFloat(strings.TrimSpace(timeDirty), 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse time entry: %w", err)
		}

		name := row[nameIndex]
		id := row[idIndex]

		x, err := strconv.ParseFloat(strings.TrimSpace(row[posXIndex]), 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse x entry: %w", err)
		}

		y, err := strconv.ParseFloat(strings.TrimSpace(row[posYIndex]), 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse y entry: %w", err)
		}

		z, err := strconv.ParseFloat(strings.TrimSpace(row[posZIndex]), 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse z entry: %w", err)
		}

		if workingData[id] == nil {
			workingData[id] = make(map[string][]position.Capture)
		}

		workingData[id][name] = append(workingData[id][name], position.NewCapture(time, x, y, z))
	}

	allRecordings := make([]format.Recording, 0)
	for id, mappings := range workingData {
		for name, captures := range mappings {
			allRecordings = append(
				allRecordings,
				format.NewRecording(
					id,
					name,
					[]format.CaptureCollection{
						position.NewStream("Position", captures),
					},
					nil,
					nil,
					nil,
				),
			)
		}
	}

	return format.NewRecording("", "", nil, allRecordings, nil, nil), nil
}
