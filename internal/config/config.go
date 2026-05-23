// Package config parses and validates CLI flags and environment variables
// into a single Config struct consumed by the splitter pipeline.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

// Config holds all runtime options for a logslice run.
type Config struct {
	// Source is the path to the input log file ("-" for stdin).
	Source string
	// Outputs is one or more destination paths ("-" for stdout).
	Outputs []string
	// Start is the inclusive lower bound of the time range (zero = unbounded).
	Start time.Time
	// End is the inclusive upper bound of the time range (zero = unbounded).
	End time.Time
	// TimeFormat is an explicit Go time layout; empty enables auto-detection.
	TimeFormat string
	// ProgressInterval controls how often progress is printed (0 = disabled).
	ProgressInterval time.Duration
	// Version requests version output and early exit.
	Version bool
}

// Parse reads os.Args using the provided FlagSet and returns a validated Config.
// fs must not have been parsed yet; callers pass flag.CommandLine for production
// or a fresh flag.NewFlagSet for tests.
func Parse(fs *flag.FlagSet, args []string) (*Config, error) {
	var (
		src      string
		start    string
		end      string
		fmt_     string
		progress time.Duration
		version  bool
	)

	fs.StringVar(&src, "src", "-", "source log file (- for stdin)")
	fs.StringVar(&start, "start", "", "start time (RFC3339 or auto-detected)")
	fs.StringVar(&end, "end", "", "end time (RFC3339 or auto-detected)")
	fs.StringVar(&fmt_, "fmt", "", "explicit Go time layout for log timestamps")
	fs.DurationVar(&progress, "progress", 5*time.Second, "progress report interval (0 to disable)")
	fs.BoolVar(&version, "version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Source:           src,
		Outputs:          fs.Args(),
		TimeFormat:       fmt_,
		ProgressInterval: progress,
		Version:          version,
	}

	if len(cfg.Outputs) == 0 {
		cfg.Outputs = []string{"-"}
	}

	var parseErr error
	if start != "" {
		cfg.Start, parseErr = time.Parse(time.RFC3339, start)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid -start: %w", parseErr)
		}
	}
	if end != "" {
		cfg.End, parseErr = time.Parse(time.RFC3339, end)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid -end: %w", parseErr)
		}
	}

	if !cfg.Start.IsZero() && !cfg.End.IsZero() && cfg.End.Before(cfg.Start) {
		return nil, errors.New("config: -end must not be before -start")
	}

	if cfg.Source != "-" {
		if _, err := os.Stat(cfg.Source); err != nil {
			return nil, fmt.Errorf("config: source file not accessible: %w", err)
		}
	}

	return cfg, nil
}
