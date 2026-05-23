package checkpoint_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestSaveAndLoad(t *testing.T) {
	store := checkpoint.New(tempPath(t))

	want := checkpoint.State{
		Offset:   1024,
		LastTime: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}

	if err := store.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.Offset != want.Offset {
		t.Errorf("Offset: got %d, want %d", got.Offset, want.Offset)
	}
	if !got.LastTime.Equal(want.LastTime) {
		t.Errorf("LastTime: got %v, want %v", got.LastTime, want.LastTime)
	}
	if got.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set after Save")
	}
}

func TestLoadNotFound(t *testing.T) {
	store := checkpoint.New(tempPath(t))

	_, err := store.Load()
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRemove(t *testing.T) {
	path := tempPath(t)
	store := checkpoint.New(path)

	if err := store.Save(checkpoint.State{Offset: 42}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := store.Remove(); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error("expected file to be deleted after Remove")
	}
}

func TestRemoveIdempotent(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	// Remove on a non-existent file should not error.
	if err := store.Remove(); err != nil {
		t.Errorf("Remove on missing file: %v", err)
	}
}

func TestSaveOverwrites(t *testing.T) {
	store := checkpoint.New(tempPath(t))

	if err := store.Save(checkpoint.State{Offset: 100}); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	if err := store.Save(checkpoint.State{Offset: 200}); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Offset != 200 {
		t.Errorf("Offset: got %d, want 200", got.Offset)
	}
}
