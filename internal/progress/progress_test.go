package progress_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/progress"
)

func TestIncCountsTotal(t *testing.T) {
	r := progress.New(io.Discard, 0)
	r.Inc(true)
	r.Inc(false)
	r.Inc(true)
	total, matched, skipped := r.Summary()
	if total != 3 {
		t.Fatalf("expected total=3, got %d", total)
	}
	if matched != 2 {
		t.Fatalf("expected matched=2, got %d", matched)
	}
	if skipped != 1 {
		t.Fatalf("expected skipped=1, got %d", skipped)
	}
}

func TestStopPrintsSummary(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(&buf, 0)
	r.Inc(true)
	r.Inc(true)
	r.Stop()
	out := buf.String()
	if !strings.Contains(out, "total=2") {
		t.Fatalf("expected total=2 in output, got: %s", out)
	}
	if !strings.Contains(out, "matched=2") {
		t.Fatalf("expected matched=2 in output, got: %s", out)
	}
	if !strings.Contains(out, "skipped=0") {
		t.Fatalf("expected skipped=0 in output, got: %s", out)
	}
}

func TestPeriodicReporting(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(&buf, 50*time.Millisecond)
	r.Start()
	r.Inc(true)
	time.Sleep(120 * time.Millisecond)
	r.Stop()
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// At least one periodic line + final summary
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 output lines, got %d: %s", len(lines), buf.String())
	}
}

func TestZeroIntervalNoBackground(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(&buf, 0)
	r.Start() // should not panic or block
	r.Inc(false)
	r.Stop()
	if buf.Len() == 0 {
		t.Fatal("expected at least a final summary line")
	}
}
