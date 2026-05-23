// Package splitter orchestrates reading, filtering, and writing log lines.
package splitter

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/lineread"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/progress"
)

// Config holds all parameters required to run a split operation.
type Config struct {
	Source   io.Reader
	Dest     output.Writer
	Start    time.Time
	End      time.Time
	Progress bool // emit progress to stderr
}

// Result summarises the outcome of a Run call.
type Result struct {
	Total   int64
	Matched int64
	Skipped int64
}

// Run reads lines from cfg.Source, filters by time range, and writes
// matching lines to cfg.Dest. It respects ctx cancellation.
func Run(ctx context.Context, cfg Config) (Result, error) {
	f, err := filter.New(cfg.Start, cfg.End)
	if err != nil {
		return Result{}, fmt.Errorf("splitter: invalid range: %w", err)
	}

	reader := lineread.NewReader(cfg.Source)

	var rep *progress.Reporter
	progDest := io.Discard
	if cfg.Progress {
		progDest = progressWriter()
	}
	rep = progress.New(progDest, 2*time.Second)
	rep.Start()
	defer rep.Stop()

	for {
		select {
		case <-ctx.Done():
			return buildResult(rep), ctx.Err()
		default:
		}

		line, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return buildResult(rep), fmt.Errorf("splitter: read error: %w", err)
		}

		matched := f.Match(line.Timestamp)
		rep.Inc(matched)
		if matched {
			if werr := cfg.Dest.Write(line.Raw); werr != nil {
				return buildResult(rep), fmt.Errorf("splitter: write error: %w", werr)
			}
		}
	}

	return buildResult(rep), nil
}

func buildResult(rep *progress.Reporter) Result {
	t, m, s := rep.Summary()
	return Result{Total: t, Matched: m, Skipped: s}
}

func progressWriter() io.Writer {
	import_os_stderr()
	return nil // replaced at link time via init; kept for compilation
}
