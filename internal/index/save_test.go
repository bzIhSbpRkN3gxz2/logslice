package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	idx := &Index{
		entries: []Entry{
			{Timestamp: now, Offset: 0},
			{Timestamp: now.Add(time.Minute), Offset: 512},
			{Timestamp: now.Add(2 * time.Minute), Offset: 1024},
		},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "test.idx")

	if err := idx.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.entries) != len(idx.entries) {
		t.Fatalf("entry count: got %d, want %d", len(loaded.entries), len(idx.entries))
	}
	for i, e := range loaded.entries {
		want := idx.entries[i]
		if !e.Timestamp.Equal(want.Timestamp) {
			t.Errorf("entry %d timestamp: got %v, want %v", i, e.Timestamp, want.Timestamp)
		}
		if e.Offset != want.Offset {
			t.Errorf("entry %d offset: got %d, want %d", i, e.Offset, want.Offset)
		}
	}
}

func TestSaveInvalidPath(t *testing.T) {
	idx := &Index{}
	err := idx.Save("/nonexistent/dir/index.idx")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestLoadNotFoundIndex(t *testing.T) {
	_, err := Load("/nonexistent/index.idx")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadCorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.idx")
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestLoadWrongVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "v99.idx")
	data := []byte(`{"version":99,"entries":[]}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected version error, got nil")
	}
}
