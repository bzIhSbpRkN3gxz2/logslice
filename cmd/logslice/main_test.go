package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	bin := filepath.Join(tmp, "logslice")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestVersionFlag(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "logslice") {
		t.Errorf("expected version string, got: %s", out)
	}
}

func TestNoArgsExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit when no args given")
	}
}

func TestExtractRangeIntegration(t *testing.T) {
	bin := buildBinary(t)

	logData := strings.Join([]string{
		"2024-01-15T09:59:00Z INFO before range",
		"2024-01-15T10:00:00Z INFO first match",
		"2024-01-15T10:30:00Z INFO second match",
		"2024-01-15T11:00:00Z INFO last match",
		"2024-01-15T11:01:00Z INFO after range",
	}, "\n") + "\n"

	tmp := t.TempDir()
	input := filepath.Join(tmp, "test.log")
	output := filepath.Join(tmp, "out.log")

	if err := os.WriteFile(input, []byte(logData), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	cmd := exec.Command(bin,
		"--start", "2024-01-15T10:00:00Z",
		"--end", "2024-01-15T11:00:00Z",
		"--out", output,
		input,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}

	result, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	got := string(result)
	for _, want := range []string{"first match", "second match", "last match"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output\ngot:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"before range", "after range"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("unexpected %q in output\ngot:\n%s", unwanted, got)
		}
	}
}
