package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/lineread"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/splitter"
	"github.com/yourorg/logslice/internal/timeparse"
)

const version = "0.1.0"

func main() {
	var (
		start   = flag.String("start", "", "start time (inclusive), e.g. 2024-01-15T10:00:00Z")
		end     = flag.String("end", "", "end time (inclusive), e.g. 2024-01-15T11:00:00Z")
		out     = flag.String("out", "", "output file path (default: stdout)")
		ver     = flag.Bool("version", false, "print version and exit")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice [flags] <logfile>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *ver {
		fmt.Printf("logslice %s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)

	parser := timeparse.NewParser()

	var startTime, endTime time.Time
	var err error

	if *start != "" {
		startTime, err = parser.Parse(*start)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logslice: invalid --start time: %v\n", err)
			os.Exit(1)
		}
	}
	if *end != "" {
		endTime, err = parser.Parse(*end)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logslice: invalid --end time: %v\n", err)
			os.Exit(1)
		}
	}

	f, err := filter.New(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}

	src, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: cannot open input: %v\n", err)
		os.Exit(1)
	}
	defer src.Close()

	reader := lineread.NewReader(src, parser)

	w, err := output.New(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: cannot open output: %v\n", err)
		os.Exit(1)
	}
	defer w.Close()

	if err := splitter.Run(reader, f, w); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}
