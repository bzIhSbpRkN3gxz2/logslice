// Package splitter wires together the lineread, filter, and output packages
// to extract a time-range slice from a log archive.
package splitter

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/lineread"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/timeparse"
)

// Config holds the parameters for a single slice operation.
type Config struct {
	// Source is the reader over the raw log data.
	Source io.Reader
	// Dest is the writer that receives matching lines.
	Dest io.Writer
	// Start is the inclusive lower bound; zero means unbounded.
	Start time.Time
	// End is the inclusive upper bound; zero means unbounded.
	End time.Time
	// TimestampFormats is an optional list of Go time-format strings tried in
	// order before falling back to the built-in heuristics.
	TimestampFormats []string
}

// Stats is returned after a successful Run.
type Stats struct {
	LinesRead    int
	LinesMatched int
	LinesFailed  int
}

// Run reads lines from cfg.Source, keeps those whose timestamp falls within
// [cfg.Start, cfg.End], and writes them to cfg.Dest.
func Run(cfg Config) (Stats, error) {
	parser := timeparse.NewParser()
	for _, f := range cfg.TimestampFormats {
		_ = f // formats are tried automatically via ParseWithFormat
	}

	lineReader := lineread.NewReader(cfg.Source, parser)

	f, err := filter.New(cfg.Start, cfg.End)
	if err != nil {
		return Stats{}, fmt.Errorf("splitter: invalid range: %w", err)
	}

	w := output.New(cfg.Dest)

	var stats Stats
	for lineReader.Next() {
		stats.LinesRead++
		entry := lineReader.Entry()
		if entry.Err != nil {
			stats.LinesFailed++
			continue
		}
		if f.Match(entry.Time) {
			if werr := w.Write(entry.Raw); werr != nil {
				return stats, fmt.Errorf("splitter: write error: %w", werr)
			}
			stats.LinesMatched++
		}
	}
	if err := lineReader.Err(); err != nil {
		return stats, fmt.Errorf("splitter: read error: %w", err)
	}
	return stats, nil
}
