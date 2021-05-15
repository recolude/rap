# RAP

**RAP is in alpha, and APIs will be changing each update**

Recolude's official recording file format.

## Install

```
git clone https://github.com/recolude/rap
cd rap
go install ./cmd/rap-cli
```

## Usage

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

## Testing Locally

You need to generate mocks before you can run parts of the test suite.

```
go generate ./...
```
