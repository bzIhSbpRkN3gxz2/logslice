package output_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestWriterFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	w, err := output.New(output.Options{Path: path})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := w.WriteLine([]byte(l)); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}

	if w.Count() != 3 {
		t.Errorf("Count() = %d, want 3", w.Count())
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	want := "line one\nline two\nline three\n"
	if string(got) != want {
		t.Errorf("file content = %q, want %q", got, want)
	}
}

func TestWriterStdout(t *testing.T) {
	// Just ensure New("-") doesn't error and Close doesn't close real stdout.
	w, err := output.New(output.Options{Path: "-"})
	if err != nil {
		t.Fatalf("New stdout: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Errorf("Close stdout writer: %v", err)
	}
}

func TestWriterCustomBuffer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "buf.log")

	w, err := output.New(output.Options{Path: path, BufferSize: 128})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	payload := bytes.Repeat([]byte("x"), 200)
	if err := w.WriteLine(payload); err != nil {
		t.Fatalf("WriteLine large: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	got, _ := os.ReadFile(path)
	if len(got) != 201 { // 200 bytes + newline
		t.Errorf("file size = %d, want 201", len(got))
	}
}

func TestWriterInvalidPath(t *testing.T) {
	_, err := output.New(output.Options{Path: "/no/such/directory/out.log"})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
