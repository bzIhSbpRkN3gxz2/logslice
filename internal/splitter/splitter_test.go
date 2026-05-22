package splitter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/splitter"
)

const sampleLog = `2024-01-15T10:00:00Z level=info msg="startup"
2024-01-15T10:05:00Z level=info msg="request received"
2024-01-15T10:10:00Z level=warn msg="slow query"
2024-01-15T10:15:00Z level=error msg="connection refused"
2024-01-15T10:20:00Z level=info msg="shutdown"
`

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestRunExtractsRange(t *testing.T) {
	var out bytes.Buffer
	stats, err := splitter.Run(splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   &out,
		Start:  mustParse("2024-01-15T10:05:00Z"),
		End:    mustParse("2024-01-15T10:10:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.LinesMatched != 2 {
		t.Errorf("expected 2 matched lines, got %d", stats.LinesMatched)
	}
	if stats.LinesRead != 5 {
		t.Errorf("expected 5 lines read, got %d", stats.LinesRead)
	}
	if !strings.Contains(out.String(), "slow query") {
		t.Errorf("expected 'slow query' in output, got: %s", out.String())
	}
}

func TestRunUnboundedRange(t *testing.T) {
	var out bytes.Buffer
	stats, err := splitter.Run(splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.LinesMatched != 5 {
		t.Errorf("expected all 5 lines matched, got %d", stats.LinesMatched)
	}
}

func TestRunEmptySource(t *testing.T) {
	var out bytes.Buffer
	stats, err := splitter.Run(splitter.Config{
		Source: strings.NewReader(""),
		Dest:   &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.LinesRead != 0 || stats.LinesMatched != 0 {
		t.Errorf("expected zero stats, got %+v", stats)
	}
}

func TestRunInvalidRange(t *testing.T) {
	var out bytes.Buffer
	_, err := splitter.Run(splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   &out,
		Start:  mustParse("2024-01-15T10:20:00Z"),
		End:    mustParse("2024-01-15T10:00:00Z"),
	})
	if err == nil {
		t.Fatal("expected error for inverted range, got nil")
	}
}
