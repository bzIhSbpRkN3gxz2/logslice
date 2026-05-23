package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempIndex(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "test.idx")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("write temp index: %v", err)
	}
	return path
}

func TestEvictStaleRemovesOldFile(t *testing.T) {
	dir := t.TempDir()
	path := writeTempIndex(t, dir)

	// Backdate the modification time by 96 hours.
	old := time.Now().Add(-96 * time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	evicted, err := EvictStale(path, time.Now(), DefaultEvictOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !evicted {
		t.Fatal("expected file to be evicted")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected file to be removed from disk")
	}
}

func TestEvictStaleKeepsFreshFile(t *testing.T) {
	dir := t.TempDir()
	path := writeTempIndex(t, dir)

	evicted, err := EvictStale(path, time.Now(), DefaultEvictOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evicted {
		t.Fatal("expected fresh file to be kept")
	}
}

func TestEvictStaleMissingFile(t *testing.T) {
	evicted, err := EvictStale("/nonexistent/path/index.idx", time.Now(), DefaultEvictOptions())
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if evicted {
		t.Fatal("expected false for missing file")
	}
}

func TestEvictStaleZeroMaxAge(t *testing.T) {
	dir := t.TempDir()
	path := writeTempIndex(t, dir)

	old := time.Now().Add(-200 * time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	opts := EvictOptions{MaxAge: 0}
	evicted, err := EvictStale(path, time.Now(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evicted {
		t.Fatal("zero MaxAge should disable eviction")
	}
}
