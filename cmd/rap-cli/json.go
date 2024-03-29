package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/metadata"
)

func toJson(out io.Writer, recording format.Recording, depth int) error {
	indentationBuilder := strings.Builder{}
	for i := 0; i < depth; i++ {
		indentationBuilder.WriteString("\t")
	}
	indentation := indentationBuilder.String()
	subIndentation := indentation + "\t"

	fmt.Fprintf(out, "%s{\n", indentation)

	// Write ID
	fmt.Fprintf(out, "%s\"id\": \"%s\",\n", subIndentation, recording.ID())

	// Write Name
	fmt.Fprintf(out, "%s\"name\": \"%s\",\n", subIndentation, recording.Name())

	// Write Metadata
	metadataJSONData, err := metadata.NewMetadataProperty(recording.Metadata()).MarshalJSON()
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\"metadata\": %s,\n", subIndentation, string(metadataJSONData))

	// Write Collections
	fmt.Fprintf(out, "%s\"collections\": [", subIndentation)
	if len(recording.CaptureCollections()) > 0 {
		fmt.Fprint(out, "\n")
	}

	subsubIndentation := subIndentation + "\t"
	for i, collection := range recording.CaptureCollections() {
		fmt.Fprintf(out, "%s{\n", subsubIndentation)
		fmt.Fprintf(out, "%s\t\"name\": \"%s\",\n", subsubIndentation, collection.Name())
		fmt.Fprintf(out, "%s\t\"signature\" : \"%s\",\n", subsubIndentation, collection.Signature())
		fmt.Fprintf(out, "%s\t\"count\" : %d", subsubIndentation, collection.Length())
		if collection.Signature() == "recolude.event" {
			fmt.Fprintf(out, ",\n%s\t\"captures\": [\n", subsubIndentation)
			eventCollection, ok := collection.(event.Collection)
			if ok {
				for capIndex := 0; capIndex < eventCollection.Length(); capIndex++ {
					event := eventCollection.CaptureAt(capIndex).(event.Capture)

					eventJSONData, err := metadata.NewMetadataProperty(event.Metadata()).MarshalJSON()
					if err != nil {
						return err
					}

					fmt.Fprintf(out, "%s\t\t{\n", subsubIndentation)
					fmt.Fprintf(out, "%s\t\t\t\"time\": %f,\n", subsubIndentation, event.Time())
					fmt.Fprintf(out, "%s\t\t\t\"name\": \"%s\",\n", subsubIndentation, event.Name())
					fmt.Fprintf(out, "%s\t\t\t\"data\": %s\n", subsubIndentation, string(eventJSONData))
					fmt.Fprintf(out, "%s\t\t}", subsubIndentation)
					if capIndex < eventCollection.Length()-1 {
						fmt.Fprintf(out, ",\n")
					} else {
						fmt.Fprintf(out, "\n")
					}
				}
			}
			fmt.Fprintf(out, "%s\t]\n", subsubIndentation)
		} else {
			fmt.Fprint(out, "\n")
		}
		fmt.Fprintf(out, "%s}", subsubIndentation)

		if i < len(recording.CaptureCollections())-1 {
			fmt.Fprintf(out, ",\n")
		}
	}
	if len(recording.CaptureCollections()) > 0 {
		fmt.Fprint(out, "\n")
		fmt.Fprintf(out, "%s],\n", subIndentation)
	} else {
		fmt.Fprint(out, "],\n")
	}

	// Write Recordings
	fmt.Fprintf(out, "%s\"recordings\": [", subIndentation)
	if len(recording.Recordings()) > 0 {
		fmt.Fprint(out, "\n")
	}

	for i, rec := range recording.Recordings() {
		if rec == nil {
			fmt.Fprintf(out, "null")
		} else {
			toJson(out, rec, depth+2)
		}
		if i < len(recording.Recordings())-1 {
			fmt.Fprintf(out, ",\n")
		}
	}

	if len(recording.Recordings()) > 0 {
		fmt.Fprint(out, "\n")
		fmt.Fprintf(out, "%s]\n", subIndentation)
	} else {
		fmt.Fprint(out, "]\n")
	}

	fmt.Fprintf(out, "%s}", indentation)
	return nil
}
