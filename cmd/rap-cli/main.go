package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/encoding"
	"github.com/recolude/rap/format/encoding/enum"
	"github.com/recolude/rap/format/encoding/euler"
	"github.com/recolude/rap/format/encoding/event"
	"github.com/recolude/rap/format/encoding/position"
	rapio "github.com/recolude/rap/format/io"
	"github.com/recolude/rap/format/parsing"
	"github.com/urfave/cli/v2"
)

func kb(byteCount int) string {
	return fmt.Sprintf("%d kb", byteCount/1024)
}

func BuildApp(in io.Reader, out io.Writer, errOut io.Writer) *cli.App {
	return &cli.App{
		Name:  "RAP CLI",
		Usage: "Utils around recolude file format",
		Authors: []*cli.Author{
			{
				Name:  "Eli Davis",
				Email: "eli@recolude.com",
			},
		},
		Version:   "1.0.0",
		Reader:    in,
		Writer:    out,
		ErrWriter: errOut,
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

					fi, err := os.Stat(fileToLoad)
					if err != nil {
						return err
					}
					size := fi.Size()

					file, err := os.Open(fileToLoad)
					if err != nil {
						return err
					}

					recording, _, err := rapio.Load(file)
					if err != nil {
						return err
					}

					printSummary(c.App.Writer, recording, size)
					return nil
				},
			},
			{
				Name: "to-json",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: false,
						Usage:    "File to turn to JSON",
					},
				},
				Usage: "Transforms a file to json",
				Action: func(c *cli.Context) error {
					var recording format.Recording

					if c.IsSet("file") {
						file, err := os.Open(c.String("file"))
						if err != nil {
							return err
						}
						recording, _, err = rapio.Load(file)
						if err != nil {
							return err
						}
					} else {
						var err error
						recording, _, err = rapio.Load(c.App.Reader)
						if err != nil {
							return err
						}
					}

					if recording == nil {
						return errors.New("can not build json from nil recording")
					}

					return toJson(c.App.Writer, recording, 0)
				},
			},
			{
				Name: "from-json",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "in",
						Aliases:  []string{"i"},
						Required: false,
						Usage:    "JSON file to build recording from",
					},
					&cli.StringFlag{
						Name:     "out",
						Aliases:  []string{"o"},
						Required: false,
						Usage:    "file to write recording too",
					},
				},
				Usage: "Transforms json to RAP",
				Action: func(c *cli.Context) error {
					jsonStream := c.App.Reader
					if c.IsSet("in") {
						file, err := os.Open(c.String("in"))
						if err != nil {
							return err
						}
						jsonStream = file
					}

					jsonData, err := ioutil.ReadAll(jsonStream)
					if err != nil {
						return err
					}

					builtRecording, err := parsing.FromJSON(jsonData)
					if err != nil {
						return err
					}

					rapStream := c.App.Writer
					if c.IsSet("out") {
						outPath := c.String("out")

						if _, err := os.Stat(outPath); err == nil {
							e := os.Remove(outPath)
							if e != nil {
								return e
							}
						}

						file, err := os.Create(outPath)
						if err != nil {
							return err
						}
						rapStream = file
					}

					encoders := []encoding.Encoder{
						event.NewEncoder(),
						position.NewEncoder(position.Oct24),
						euler.NewEncoder(euler.Raw16),
						enum.NewEncoder(),
					}

					recordingWriter := rapio.NewWriter(encoders, true, rapStream, rapio.BST16)
					_, err = recordingWriter.Write(builtRecording)
					return err
				},
			},
			{
				Name: "upgrade",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "File to upgrade",
					},
				},
				Usage: "Upgrades a file from v1 to v2",
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					file, err := os.Open(fileToLoad)
					if err != nil {
						return err
					}

					recording, _, err := rapio.Load(file)
					if err != nil {
						return err
					}

					encoders := []encoding.Encoder{
						event.NewEncoder(),
						position.NewEncoder(position.Oct24),
						euler.NewEncoder(euler.Raw16),
						enum.NewEncoder(),
					}

					recordingWriter := rapio.NewWriter(encoders, true, c.App.Writer, rapio.BST16)
					_, err = recordingWriter.Write(recording)
					return err
				},
			},
			{
				Name: "from-csv",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "File to turn to upgrade",
					},
				},
				Usage: "Builds a recording from CSV",
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					csvStream, err := os.Open(fileToLoad)
					if err != nil {
						return err
					}

					recording, err := RecordingFromCSV(csvStream)
					if err != nil {
						return err
					}

					encoders := []encoding.Encoder{
						position.NewEncoder(position.Raw64),
					}

					recordingWriter := rapio.NewWriter(encoders, true, c.App.Writer, rapio.Raw64)
					_, err = recordingWriter.Write(recording)
					return err
				},
			},
		},
	}
}

func main() {
	app := BuildApp(os.Stdin, os.Stdout, os.Stderr)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
