// Package timeparse provides utilities for parsing timestamps from log lines.
package timeparse

import (
	"fmt"
	"time"
)

// Common log timestamp formats ordered by specificity.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// Parser holds a cached format for repeated parsing of same-format logs.
type Parser struct {
	cachedFormat string
}

// NewParser returns a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// Parse attempts to extract a time.Time from the given raw timestamp string.
// It tries the cached format first, then falls back to all known formats.
func (p *Parser) Parse(raw string) (time.Time, error) {
	if p.cachedFormat != "" {
		if t, err := time.Parse(p.cachedFormat, raw); err == nil {
			return t, nil
		}
	}

	for _, format := range knownFormats {
		if t, err := time.Parse(format, raw); err == nil {
			p.cachedFormat = format
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("timeparse: unrecognized timestamp format: %q", raw)
}

// ParseWithFormat parses a timestamp using an explicit format string.
func ParseWithFormat(raw, format string) (time.Time, error) {
	t, err := time.Parse(format, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("timeparse: failed to parse %q with format %q: %w", raw, format, err)
	}
	return t, nil
}
