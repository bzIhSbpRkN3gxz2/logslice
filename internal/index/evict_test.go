package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempIndex(t *testing.T, dir string) string {
	t.Helper()
	p := filepath.Join(dir, "test.idx")
	if err := os.WriteFile(p, []byte("data"), 0o600); err != nil {
		t.Fatalf("writeTempIndex: %v", err)
	}
	return p
}

func TestEvictStaleRemovesOldFile(t *testing.T) {
	dir := t.TempDir()
	p := writeTempIndex(t, dir)

	// Back-date the file modification time.
	old := time.Now().Add(-100 * time.Hour)
	if err := os.Chtimes(p, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	opts := EvictOptions{MaxAge: 72 * time.Hour}
	if err := EvictStale(p, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("expected file to be removed")
	}
}

func TestEvictStaleKeepsFreshFile(t *testing.T) {
	dir := t.TempDir()
	p := writeTempIndex(t, dir)

	opts := EvictOptions{MaxAge: 72 * time.Hour}
	if err := EvictStale(p, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(p); err != nil {
		t.Fatalf("expected file to still exist: %v", err)
	}
}

func TestEvictStaleMissingFile(t *testing.T) {
	opts := DefaultEvictOptions()
	if err := EvictStale("/nonexistent/path/file.idx", opts); err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestEvictStaleZeroMaxAge(t *testing.T) {
	dir := t.TempDir()
	p := writeTempIndex(t, dir)

	// Back-date so it would normally be evicted.
	old := time.Now().Add(-200 * time.Hour)
	if err := os.Chtimes(p, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	opts := EvictOptions{MaxAge: 0}
	if err := EvictStale(p, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file should not have been removed: %v", err)
	}
}
