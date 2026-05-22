// Package timeparse provides timestamp parsing with automatic format detection
// and format caching for high-throughput log processing.
package timeparse

import (
	"fmt"
	"strings"
	"time"
)

// knownFormats is the ordered list of formats tried during auto-detection.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan  2 15:04:05",
	"Jan 02 15:04:05",
}

// Parser detects and caches the timestamp format found in log lines.
type Parser struct {
	cachedFormat string
}

// NewParser returns a new Parser with no cached format.
func NewParser() *Parser {
	return &Parser{}
}

// Parse attempts to extract and parse a timestamp from the beginning of line.
// It tries the cached format first, then falls back to all known formats.
func (p *Parser) Parse(line string) (time.Time, error) {
	// Grab a reasonable prefix to avoid scanning the whole line.
	prefix := line
	if len(prefix) > 35 {
		prefix = prefix[:35]
	}
	// Strip trailing content after the first space following the timestamp.
	fields := strings.SplitN(prefix, " ", 2)
	timePart := fields[0]

	if p.cachedFormat != "" {
		if t, err := time.Parse(p.cachedFormat, timePart); err == nil {
			return t, nil
		}
	}

	for _, fmt := range knownFormats {
		if t, err := time.Parse(fmt, timePart); err == nil {
			p.cachedFormat = fmt
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("lineread: no known format matched %q", timePart)
}

// ParseWithFormat parses s using the explicit Go time layout format.
func ParseWithFormat(s, format string) (time.Time, error) {
	t, err := time.Parse(format, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("timeparse: ParseWithFormat: %w", err)
	}
	return t, nil
}
