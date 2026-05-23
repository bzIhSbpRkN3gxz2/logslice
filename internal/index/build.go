package index

import (
	"bufio"
	"io"

	"github.com/yourorg/logslice/internal/timeparse"
)

// BuildOptions controls how densely the index is sampled.
type BuildOptions struct {
	// SampleEvery instructs the builder to record one entry every N lines.
	// A value of 0 or 1 records every line (not recommended for large files).
	SampleEvery int
}

// DefaultBuildOptions returns sensible defaults.
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{SampleEvery: 1000}
}

// Build scans r, sampling line offsets according to opts, and returns a
// populated Index. parser is used to extract timestamps from each line.
func Build(r io.ReadSeeker, parser *timeparse.Parser, opts BuildOptions) (*Index, error) {
	if opts.SampleEvery <= 0 {
		opts.SampleEvery = 1
	}

	var (
		idx    Index
		offset int64
		lineNo int
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw := scanner.Bytes()
		lineNo++

		if lineNo%opts.SampleEvery == 0 {
			ts, err := parser.Parse(string(raw))
			if err == nil {
				idx.Add(ts, offset)
			}
		}

		// +1 for the newline consumed by the scanner.
		offset += int64(len(raw)) + 1
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &idx, nil
}
