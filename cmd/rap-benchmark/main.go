package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

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

func printRecording(out io.Writer, recording format.Recording, depth int) {
	fmt.Fprintf(out, "Name: %s\n", recording.Name())
	fmt.Fprintf(out, "Collections: %d\n", len(recording.CaptureCollections()))
	for _, collection := range recording.CaptureCollections() {
		fmt.Fprintf(out, "  Name: %s\n", collection.Name())
		fmt.Fprintf(out, "  Signature: %s\n", collection.Signature())
		fmt.Fprintf(out, "  Captures: %d\n", collection.Length())
	}
	fmt.Fprintf(out, "Sub Recordings: %d\n", len(recording.Recordings()))
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: true,
						Usage:    "File to run benchmark on",
					},
				},
				Aliases: []string{"r"},
				Usage:   "Run benchmark",
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					timeLoadStarted := time.Now()
					file, err := os.Open(fileToLoad)
					if err != nil {
						panic(err)
					}
					recording, originalBytesRead, err := rapio.Load(file)
					if err != nil {
						panic(err)
					}
					fmt.Fprintf(c.App.Writer, "Time to read in: %s\n", time.Now().Sub(timeLoadStarted))
					printRecording(c.App.Writer, recording, 0)

					fmt.Fprintf(c.App.Writer, "Original Size: %s\n\n", kb(originalBytesRead))

					encoders := []encoding.Encoder{
						event.NewEncoder(),
						position.NewEncoder(position.Oct48),
						euler.NewEncoder(euler.Raw32),
						enum.NewEncoder(),
					}

					recBuf := bytes.Buffer{}
					recordingWriter := rapio.NewWriter(encoders, true, &recBuf, rapio.BST16)
					recordingReader := rapio.NewReader(encoders, &recBuf)

					_, err = recordingWriter.Write(recording)
					if err != nil {
						panic(err)
					}

					newSize := kb((len(recBuf.Bytes())))

					recBack, _, err := recordingReader.Read()
					if err != nil {
						panic(err)
					}
					printRecording(c.App.Writer, recBack, 0)

					fmt.Fprintf(c.App.Writer, "New Size: %s\n", newSize)

					return nil
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
