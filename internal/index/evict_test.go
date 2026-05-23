package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempIndex(t *testing.T, modTime time.Time) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.idx")
	if err := os.WriteFile(p, []byte("data"), 0o600); err != nil {
		t.Fatalf("write temp index: %v", err)
	}
	if err := os.Chtimes(p, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	return p
}

func TestEvictStaleRemovesOldFile(t *testing.T) {
	old := time.Now().Add(-100 * time.Hour)
	p := writeTempIndex(t, old)
	opts := EvictOptions{MaxAge: 72 * time.Hour}

	evicted, err := EvictStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !evicted {
		t.Fatal("expected file to be evicted")
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("file should have been removed")
	}
}

func TestEvictStaleKeepsFreshFile(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour)
	p := writeTempIndex(t, recent)
	opts := EvictOptions{MaxAge: 72 * time.Hour}

	evicted, err := EvictStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evicted {
		t.Fatal("expected file to be kept")
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file should still exist: %v", err)
	}
}

func TestEvictStaleMissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nonexistent.idx")
	opts := DefaultEvictOptions()

	evicted, err := EvictStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if evicted {
		t.Fatal("should not report eviction for missing file")
	}
}

func TestEvictStaleZeroMaxAge(t *testing.T) {
	old := time.Now().Add(-200 * time.Hour)
	p := writeTempIndex(t, old)
	opts := EvictOptions{MaxAge: 0}

	evicted, err := EvictStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evicted {
		t.Fatal("zero MaxAge should disable eviction")
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file should still exist: %v", err)
	}
}
