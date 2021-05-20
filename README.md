# RAP

**RAP is in alpha, and APIs will be changing each update**

Recolude's official recording file format.

## Install

```
git clone https://github.com/recolude/rap
cd rap
go install ./cmd/rap-cli
```

## CLI Usage

```
NAME:
   RAP CLI - Utils around recolude file format

USAGE:
   rap-cli [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Eli Davis <eli@recolude.com>

COMMANDS:
   from-csv   Builds a recording from CSV
   json       Transforms a file to json
   summarize  Summarizes a file
   upgrade    Upgrades a file from v1 to v2
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Building Recordings Programmatically

With this new library you can create your own recordings programmatically. The below example creates a recording of the sin wave and then writes it to disk.

```golang
package main

import (
	"math"
	"os"
	"time"

	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/io"
	"github.com/recolude/rap/format/metadata"
)

func main() {
	iterations := 1000
	positions := make([]position.Capture, iterations)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		currentTime := float64(i)
		positions[i] = position.NewCapture(currentTime, 0, math.Sin(currentTime), 0)
	}
	duration := time.Since(start)

	rec := format.NewRecording(
		"",
		"Sin Wave Demo",
		[]format.CaptureCollection{
			position.NewCollection("Sin Wave", positions),
		},
		nil,
		metadata.NewBlock(map[string]metadata.Property{
			"iterations": metadata.NewIntProperty(int32(iterations)),
			"benchmark":  metadata.NewStringProperty(duration.String()),
		}),
		nil,
		nil,
	)

	f, _ := os.Create("sin demo.rap")
	recordingWriter := io.NewRecoludeWriter(f)
	recordingWriter.Write(rec)
}
```

## Testing Locally

You need to generate mocks before you can run parts of the test suite.

```
go generate ./...
```
