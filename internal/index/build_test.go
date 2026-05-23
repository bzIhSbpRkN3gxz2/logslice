package index

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/timeparse"
)

const sampleLogs = `2024-01-01T00:00:00Z INFO  server started
2024-01-01T00:01:00Z DEBUG request received
2024-01-01T00:02:00Z INFO  processed ok
2024-01-01T00:03:00Z WARN  slow query detected
2024-01-01T00:04:00Z ERROR disk nearly full
`

func TestBuildSamplesLines(t *testing.T) {
	parser := timeparse.NewParser()
	r := strings.NewReader(sampleLogs)

	idx, err := Build(r, parser, BuildOptions{SampleEvery: 2})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(idx.Entries) == 0 {
		t.Fatal("expected at least one index entry")
	}
	// Every 2nd line should be sampled (lines 2, 4).
	if len(idx.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(idx.Entries))
	}
}

func TestBuildDefaultOptions(t *testing.T) {
	parser := timeparse.NewParser()
	r := strings.NewReader(sampleLogs)

	opts := DefaultBuildOptions()
	if opts.SampleEvery != 1000 {
		t.Fatalf("unexpected default SampleEvery: %d", opts.SampleEvery)
	}

	// With SampleEvery=1000 and only 5 lines, no entries should be recorded.
	idx, err := Build(r, parser, opts)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Fatalf("expected 0 entries with high sample rate, got %d", len(idx.Entries))
	}
}

func TestBuildOffsetsIncrease(t *testing.T) {
	parser := timeparse.NewParser()
	r := strings.NewReader(sampleLogs)

	idx, err := Build(r, parser, BuildOptions{SampleEvery: 1})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	for i := 1; i < len(idx.Entries); i++ {
		if idx.Entries[i].Offset <= idx.Entries[i-1].Offset {
			t.Errorf("offsets not strictly increasing at index %d", i)
		}
	}
}

func TestBuildZeroSampleEveryDefaultsToOne(t *testing.T) {
	parser := timeparse.NewParser()
	r := strings.NewReader(sampleLogs)

	idx, err := Build(r, parser, BuildOptions{SampleEvery: 0})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	// All 5 lines have valid timestamps, so expect 5 entries.
	if len(idx.Entries) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(idx.Entries))
	}
}
