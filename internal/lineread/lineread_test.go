package lineread_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/lineread"
)

const sampleLog = `2024-01-15T10:00:00Z INFO  starting server
2024-01-15T10:00:01Z DEBUG listening on :8080
2024-01-15T10:00:05Z ERROR connection refused
not a timestamped line
2024-01-15T10:00:10Z INFO  shutting down
`

func TestReaderParsesTimestamps(t *testing.T) {
	r := lineread.NewReader(strings.NewReader(sampleLog))

	line, ok := r.Next()
	if !ok {
		t.Fatal("expected a line, got none")
	}
	if !line.Valid {
		t.Errorf("expected valid timestamp on first line")
	}
	want := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !line.Timestamp.Equal(want) {
		t.Errorf("timestamp = %v, want %v", line.Timestamp, want)
	}
}

func TestReaderInvalidTimestampLine(t *testing.T) {
	r := lineread.NewReader(strings.NewReader(sampleLog))

	var lines []lineread.Line
	for {
		l, ok := r.Next()
		if !ok {
			break
		}
		lines = append(lines, l)
	}
	if err := r.Err(); err != nil {
		t.Fatalf("unexpected scanner error: %v", err)
	}
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
	// 4th line has no timestamp
	if lines[3].Valid {
		t.Errorf("expected invalid timestamp on line 4")
	}
	if lines[3].Raw != "not a timestamped line" {
		t.Errorf("unexpected raw content: %q", lines[3].Raw)
	}
}

func TestReaderEmptyInput(t *testing.T) {
	r := lineread.NewReader(strings.NewReader(""))
	_, ok := r.Next()
	if ok {
		t.Error("expected no lines from empty input")
	}
	if err := r.Err(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReaderRawPreserved(t *testing.T) {
	input := "2024-03-01T08:30:00Z WARN  disk usage high\n"
	r := lineread.NewReader(strings.NewReader(input))
	line, ok := r.Next()
	if !ok {
		t.Fatal("expected a line")
	}
	want := "2024-03-01T08:30:00Z WARN  disk usage high"
	if line.Raw != want {
		t.Errorf("Raw = %q, want %q", line.Raw, want)
	}
}
