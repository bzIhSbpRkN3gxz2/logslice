package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

var (
	t0 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
)

func TestFilterMatch(t *testing.T) {
	f, err := filter.New(filter.Range{Start: t1, End: t2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"before start", t0, false},
		{"at start", t1, true},
		{"inside", t1.Add(30 * time.Minute), true},
		{"at end (exclusive)", t2, false},
		{"after end", t3, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := f.Match(c.t); got != c.want {
				t.Errorf("Match(%v) = %v, want %v", c.t, got, c.want)
			}
		})
	}
}

func TestFilterUnboundedStart(t *testing.T) {
	f, _ := filter.New(filter.Range{End: t2})
	if !f.Match(t0) {
		t.Error("expected match before end with no start bound")
	}
	if f.Match(t2) {
		t.Error("expected no match at end boundary")
	}
}

func TestFilterUnboundedEnd(t *testing.T) {
	f, _ := filter.New(filter.Range{Start: t1})
	if f.Match(t0) {
		t.Error("expected no match before start with no end bound")
	}
	if !f.Match(t3) {
		t.Error("expected match after start with no end bound")
	}
}

func TestFilterInvalidRange(t *testing.T) {
	_, err := filter.New(filter.Range{Start: t2, End: t1})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
	if err != filter.ErrInvalidRange {
		t.Errorf("got %v, want ErrInvalidRange", err)
	}
}

func TestFilterPast(t *testing.T) {
	f, _ := filter.New(filter.Range{Start: t1, End: t2})
	if f.Past(t1) {
		t.Error("t1 should not be past end")
	}
	if !f.Past(t2) {
		t.Error("t2 should be past end (exclusive)")
	}
	if !f.Past(t3) {
		t.Error("t3 should be past end")
	}
}
