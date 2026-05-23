package splitter_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/splitter"
)

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

const sampleLog = `2024-01-10T10:00:00Z INFO startup
2024-01-10T10:01:00Z INFO request received
2024-01-10T10:02:00Z WARN slow query
2024-01-10T10:03:00Z INFO request received
2024-01-10T10:04:00Z ERROR timeout
`

func TestRunExtractsRange(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf)
	res, err := splitter.Run(context.Background(), splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   w,
		Start:  mustParse("2024-01-10T10:01:00Z"),
		End:    mustParse("2024-01-10T10:03:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 3 {
		t.Fatalf("expected 3 matched lines, got %d", res.Matched)
	}
	if res.Total != 5 {
		t.Fatalf("expected 5 total lines, got %d", res.Total)
	}
}

func TestRunUnboundedRange(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf)
	res, err := splitter.Run(context.Background(), splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   w,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 5 {
		t.Fatalf("expected all 5 lines matched, got %d", res.Matched)
	}
}

func TestRunEmptySource(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf)
	res, err := splitter.Run(context.Background(), splitter.Config{
		Source: strings.NewReader(""),
		Dest:   w,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 0 {
		t.Fatalf("expected 0 total, got %d", res.Total)
	}
}

func TestRunInvalidRange(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf)
	_, err := splitter.Run(context.Background(), splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   w,
		Start:  mustParse("2024-01-10T10:05:00Z"),
		End:    mustParse("2024-01-10T10:01:00Z"),
	})
	if err == nil {
		t.Fatal("expected error for inverted range, got nil")
	}
}

func TestRunContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	var buf bytes.Buffer
	w, _ := output.New(&buf)
	_, err := splitter.Run(ctx, splitter.Config{
		Source: strings.NewReader(sampleLog),
		Dest:   w,
	})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
