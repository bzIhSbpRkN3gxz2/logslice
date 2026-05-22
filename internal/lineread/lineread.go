// Package lineread provides utilities for reading log files line by line
// with support for extracting timestamps from each line.
package lineread

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
)

// Line represents a single log line with its parsed timestamp.
type Line struct {
	Raw       string
	Timestamp time.Time
	Valid     bool // true if a timestamp was successfully parsed
}

// Reader reads log lines and attempts to parse timestamps from each.
type Reader struct {
	scanner *bufio.Scanner
	parser  *timeparse.Parser
}

// NewReader creates a new Reader that reads from r.
// It uses a shared timeparse.Parser for efficient format caching.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(r),
		parser:  timeparse.NewParser(),
	}
}

// Next advances to the next line and returns it.
// Returns false when there are no more lines or an error occurs.
func (r *Reader) Next() (Line, bool) {
	if !r.scanner.Scan() {
		return Line{}, false
	}
	raw := r.scanner.Text()
	ts, err := r.parser.Parse(raw)
	return Line{
		Raw:       raw,
		Timestamp: ts,
		Valid:     err == nil,
	}, true
}

// Err returns any error encountered during scanning.
func (r *Reader) Err() error {
	return r.scanner.Err()
}
