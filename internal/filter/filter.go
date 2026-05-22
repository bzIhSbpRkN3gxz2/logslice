package filter

import "time"

// Range represents a time window [Start, End) for log filtering.
type Range struct {
	Start time.Time
	End   time.Time
}

// Filter decides whether log lines fall within a time range.
type Filter struct {
	r Range
}

// New creates a Filter for the given time range.
// Returns an error if Start is after End.
func New(r Range) (*Filter, error) {
	if !r.Start.IsZero() && !r.End.IsZero() && r.Start.After(r.End) {
		return nil, ErrInvalidRange
	}
	return &Filter{r: r}, nil
}

// Match reports whether t falls within the filter's range.
// A zero Start or End is treated as unbounded on that side.
func (f *Filter) Match(t time.Time) bool {
	if !f.r.Start.IsZero() && t.Before(f.r.Start) {
		return false
	}
	if !f.r.End.IsZero() && !t.Before(f.r.End) {
		return false
	}
	return true
}

// Past reports whether t is strictly before the filter's start.
// Useful for early-exit when iterating sorted log lines.
func (f *Filter) Past(t time.Time) bool {
	return !f.r.End.IsZero() && !t.Before(f.r.End)
}
