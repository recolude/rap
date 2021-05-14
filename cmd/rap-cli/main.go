package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	"github.com/recolude/rap/format/encoding/enum"
	"github.com/recolude/rap/format/encoding/euler"
	"github.com/recolude/rap/format/encoding/event"
	"github.com/recolude/rap/format/encoding/position"
	rapio "github.com/recolude/rap/format/io"
	"github.com/urfave/cli/v2"
)

func kb(byteCount int) string {
	return fmt.Sprintf("%d kb", byteCount/1024)
}

func toJson(out io.Writer, recording format.Recording) {
	fmt.Fprintf(out, "{ \"name\": \"%s\", ", recording.Name())

	fmt.Fprint(out, "\"streams\": [")
	for i, stream := range recording.CaptureStreams() {
		fmt.Fprintf(out, "{ \"name\": \"%s\", ", stream.Name())
		fmt.Fprintf(out, " \"signature\" : \"%s\" }", stream.Signature())
		if i < len(recording.CaptureStreams())-1 {
			fmt.Fprintf(out, ",")
		}
	}
	fmt.Fprint(out, "],")

	fmt.Fprint(out, "\"recordings\": [")
	for i, rec := range recording.Recordings() {
		if rec == nil {
			fmt.Fprintf(out, "null")
		} else {
			toJson(out, rec)
		}
		if i < len(recording.Recordings())-1 {
			fmt.Fprintf(out, ",")
		}
	}
	fmt.Fprint(out, "]")

	fmt.Fprint(out, "}")
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "summarize",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "File to turn to summarize",
					},
				},
				Usage: "Summarizes a file",
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					file, err := os.Open(fileToLoad)
					if err != nil {
						panic(err)
					}
					recording, _, err := rapio.Load(file)
					if err != nil {
						panic(err)
					}
					printSummary(c.App.Writer, recording)
					return nil
				},
			},
			{
				Name: "json",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "File to turn to JSON",
					},
				},
				Usage: "Transforms a file to json",
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					file, err := os.Open(fileToLoad)
					if err != nil {
						panic(err)
					}
					recording, _, err := rapio.Load(file)
					if err != nil {
						panic(err)
					}
					toJson(c.App.Writer, recording)

					return nil
				},
			},
			{
				Name: "upgrade",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "File to turn to upgrade",
					},
				},
				Usage: "Upgrades a file from v1 to v2",
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					file, err := os.Open(fileToLoad)
					if err != nil {
						panic(err)
					}
					recording, _, err := rapio.Load(file)
					if err != nil {
						panic(err)
					}

					encoders := []encoding.Encoder{
						event.NewEncoder(event.Raw32),
						position.NewEncoder(position.Oct24),
						euler.NewEncoder(euler.Raw16),
						enum.NewEncoder(enum.Raw32),
					}

					recordingWriter := rapio.NewWriter(encoders, c.App.Writer)
					_, err = recordingWriter.Write(recording)
					return err
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
