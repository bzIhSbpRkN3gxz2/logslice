package index

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, lines string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatalf("create temp log: %v", err)
	}
	if _, err := f.WriteString(lines); err != nil {
		t.Fatalf("write temp log: %v", err)
	}
	f.Close()
	return f.Name()
}

const sampleLog = "2024-01-01T00:00:00Z line one\n" +
	"2024-01-01T00:01:00Z line two\n" +
	"2024-01-01T00:02:00Z line three\n"

func TestManagerBuildsAndCaches(t *testing.T) {
	logFile := writeTempLog(t, sampleLog)
	cacheDir := t.TempDir()

	m := NewManager(cacheDir, time.Hour)
	idx, err := m.Get(logFile, DefaultBuildOptions())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(idx.Entries) == 0 {
		t.Fatal("expected non-empty index")
	}

	// Second call should return cached value without rebuilding.
	idx2, err := m.Get(logFile, DefaultBuildOptions())
	if err != nil {
		t.Fatalf("Get (cached): %v", err)
	}
	if idx != idx2 {
		t.Error("expected same pointer from cache")
	}
}

func TestManagerPersistsIndexFile(t *testing.T) {
	logFile := writeTempLog(t, sampleLog)
	cacheDir := t.TempDir()

	m := NewManager(cacheDir, time.Hour)
	if _, err := m.Get(logFile, DefaultBuildOptions()); err != nil {
		t.Fatalf("Get: %v", err)
	}

	expected := filepath.Join(cacheDir, filepath.Base(logFile)+".idx")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("index file not persisted: %v", err)
	}
}

func TestManagerLoadsFromDisk(t *testing.T) {
	logFile := writeTempLog(t, sampleLog)
	cacheDir := t.TempDir()

	m1 := NewManager(cacheDir, time.Hour)
	if _, err := m1.Get(logFile, DefaultBuildOptions()); err != nil {
		t.Fatalf("first Get: %v", err)
	}

	// New manager — no in-memory cache, should load from disk.
	m2 := NewManager(cacheDir, time.Hour)
	idx, err := m2.Get(logFile, DefaultBuildOptions())
	if err != nil {
		t.Fatalf("second Get: %v", err)
	}
	if len(idx.Entries) == 0 {
		t.Error("expected entries loaded from disk")
	}
}

func TestManagerInvalidate(t *testing.T) {
	logFile := writeTempLog(t, sampleLog)
	cacheDir := t.TempDir()

	m := NewManager(cacheDir, time.Hour)
	if _, err := m.Get(logFile, DefaultBuildOptions()); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if err := m.Invalidate(logFile); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	if _, ok := m.cache[logFile]; ok {
		t.Error("cache entry should be removed after invalidation")
	}
}

func TestManagerMissingLogFile(t *testing.T) {
	m := NewManager(t.TempDir(), time.Hour)
	_, err := m.Get("/nonexistent/file.log", DefaultBuildOptions())
	if err == nil {
		t.Fatal("expected error for missing log file")
	}
	if !strings.Contains(err.Error(), "open log file") {
		t.Errorf("unexpected error: %v", err)
	}
}
