package index

import (
	"testing"
	"time"
)

func makeTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestFloorOffsetEmpty(t *testing.T) {
	var idx Index
	if got := idx.FloorOffset(makeTime("2024-01-01T00:00:00Z")); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestFloorOffsetExact(t *testing.T) {
	var idx Index
	idx.Add(makeTime("2024-01-01T00:00:00Z"), 0)
	idx.Add(makeTime("2024-01-01T01:00:00Z"), 500)
	idx.Add(makeTime("2024-01-01T02:00:00Z"), 1000)

	got := idx.FloorOffset(makeTime("2024-01-01T01:00:00Z"))
	if got != 500 {
		t.Fatalf("expected 500, got %d", got)
	}
}

func TestFloorOffsetBetween(t *testing.T) {
	var idx Index
	idx.Add(makeTime("2024-01-01T00:00:00Z"), 0)
	idx.Add(makeTime("2024-01-01T02:00:00Z"), 800)

	got := idx.FloorOffset(makeTime("2024-01-01T01:00:00Z"))
	if got != 0 {
		t.Fatalf("expected 0 (floor entry), got %d", got)
	}
}

func TestFloorOffsetBeforeAll(t *testing.T) {
	var idx Index
	idx.Add(makeTime("2024-01-01T06:00:00Z"), 999)

	got := idx.FloorOffset(makeTime("2024-01-01T00:00:00Z"))
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSaveLoad(t *testing.T) {
	path := t.TempDir() + "/test.idx"

	var orig Index
	orig.Add(makeTime("2024-06-01T10:00:00Z"), 0)
	orig.Add(makeTime("2024-06-01T11:00:00Z"), 2048)

	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Fatalf("entry count mismatch: want %d got %d", len(orig.Entries), len(loaded.Entries))
	}
	for i, e := range orig.Entries {
		if !loaded.Entries[i].Timestamp.Equal(e.Timestamp) || loaded.Entries[i].Offset != e.Offset {
			t.Errorf("entry %d mismatch", i)
		}
	}
}

func TestLoadNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/index.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
