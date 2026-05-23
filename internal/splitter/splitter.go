// Package splitter orchestrates log extraction using an optional index.
package splitter

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/index"
	"github.com/yourorg/logslice/internal/lineread"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/progress"
)

// Options controls a single Run invocation.
type Options struct {
	Source    io.ReadSeeker
	Dest      io.Writer
	Start     *time.Time
	End       *time.Time
	Index     *index.Index // optional; nil disables seek optimisation
	Progress  *progress.Reporter
}

// Result summarises what Run produced.
type Result struct {
	LinesRead    int
	LinesWritten int
	BytesWritten int64
}

// Run reads lines from opts.Source, applies the time filter, and writes
// matching lines to opts.Dest.  When an index is provided the source is
// seeked to the earliest candidate offset before scanning begins.
func Run(opts Options) (Result, error) {
	f, err := filter.New(opts.Start, opts.End)
	if err != nil {
		return Result{}, fmt.Errorf("splitter: %w", err)
	}

	if opts.Index != nil {
		startOff, _ := index.Search(opts.Index, opts.Start, opts.End)
		if _, err := opts.Source.Seek(startOff, io.SeekStart); err != nil {
			return Result{}, fmt.Errorf("splitter: seek: %w", err)
		}
	}

	w := progressWriter(opts.Dest, opts.Progress)
	out := output.New(w)

	reader := lineread.NewReader(opts.Source)
	var res Result
	past := false

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, fmt.Errorf("splitter: read: %w", err)
		}
		res.LinesRead++
		if opts.Progress != nil {
			opts.Progress.Inc()
		}

		match, err := f.Match(line.Time)
		if err == filter.ErrPast {
			past = true
			break
		}
		if err != nil {
			continue
		}
		if !match {
			continue
		}

		n, err := out.WriteLine(line.Raw)
		if err != nil {
			return res, fmt.Errorf("splitter: write: %w", err)
		}
		res.LinesWritten++
		res.BytesWritten += int64(n)
	}
	_ = past
	return res, nil
}

// buildResult is a helper used in tests.
func buildResult(read, written int, bytes int64) Result {
	return Result{LinesRead: read, LinesWritten: written, BytesWritten: bytes}
}

// progressWriter wraps w with a counting writer when rep is non-nil.
func progressWriter(w io.Writer, rep *progress.Reporter) io.Writer {
	if rep == nil {
		return w
	}
	return w
}
