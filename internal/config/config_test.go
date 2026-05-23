package config_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestParseDefaults(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Source != "-" {
		t.Errorf("expected default source '-', got %q", cfg.Source)
	}
	if len(cfg.Outputs) != 1 || cfg.Outputs[0] != "-" {
		t.Errorf("expected default output ['-'], got %v", cfg.Outputs)
	}
	if cfg.ProgressInterval != 5*time.Second {
		t.Errorf("expected 5s progress interval, got %v", cfg.ProgressInterval)
	}
}

func TestParseTimeRange(t *testing.T) {
	args := []string{
		"-start", "2024-01-01T00:00:00Z",
		"-end", "2024-01-02T00:00:00Z",
	}
	cfg, err := config.Parse(newFS(), args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Start.Year() != 2024 {
		t.Errorf("unexpected start year: %d", cfg.Start.Year())
	}
	if cfg.End.Day() != 2 {
		t.Errorf("unexpected end day: %d", cfg.End.Day())
	}
}

func TestParseInvalidStartTime(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"-start", "not-a-time"})
	if err == nil {
		t.Fatal("expected error for invalid -start")
	}
}

func TestParseEndBeforeStart(t *testing.T) {
	args := []string{
		"-start", "2024-06-01T00:00:00Z",
		"-end", "2024-01-01T00:00:00Z",
	}
	_, err := config.Parse(newFS(), args)
	if err == nil {
		t.Fatal("expected error when end is before start")
	}
}

func TestParseMissingSourceFile(t *testing.T) {
	_, err := config.Parse(newFS(), []string{"-src", "/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestParseExistingSourceFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := config.Parse(newFS(), []string{"-src", filepath.Clean(f.Name())})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Source != f.Name() {
		t.Errorf("source mismatch: got %q", cfg.Source)
	}
}

func TestParseVersionFlag(t *testing.T) {
	cfg, err := config.Parse(newFS(), []string{"-version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Version {
		t.Error("expected Version=true")
	}
}
