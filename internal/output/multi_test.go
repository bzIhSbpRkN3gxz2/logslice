package output_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestMultiWriterFansOut(t *testing.T) {
	dir := t.TempDir()

	paths := []string{
		filepath.Join(dir, "a.log"),
		filepath.Join(dir, "b.log"),
	}

	var writers []*output.Writer
	for _, p := range paths {
		w, err := output.New(output.Options{Path: p})
		if err != nil {
			t.Fatalf("New %q: %v", p, err)
		}
		writers = append(writers, w)
	}

	mw := output.NewMulti(writers...)

	lines := [][]byte{[]byte("alpha"), []byte("beta"), []byte("gamma")}
	for _, l := range lines {
		if err := mw.WriteLine(l); err != nil {
			t.Fatalf("MultiWriter.WriteLine: %v", err)
		}
	}

	if err := mw.Close(); err != nil {
		t.Fatalf("MultiWriter.Close: %v", err)
	}

	want := "alpha\nbeta\ngamma\n"
	for _, p := range paths {
		got, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("ReadFile %q: %v", p, err)
		}
		if string(got) != want {
			t.Errorf("%q content = %q, want %q", p, got, want)
		}
	}
}

func TestMultiWriterWritersAccessor(t *testing.T) {
	dir := t.TempDir()
	w1, _ := output.New(output.Options{Path: filepath.Join(dir, "x.log")})
	w2, _ := output.New(output.Options{Path: filepath.Join(dir, "y.log")})
	mw := output.NewMulti(w1, w2)
	defer mw.Close()

	if len(mw.Writers()) != 2 {
		t.Errorf("Writers() len = %d, want 2", len(mw.Writers()))
	}
}
