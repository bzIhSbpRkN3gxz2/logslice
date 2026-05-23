// Package checkpoint provides resumable processing support by persisting
// the last successfully processed byte offset and timestamp for a log file.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrNotFound is returned when no checkpoint file exists for the given key.
var ErrNotFound = errors.New("checkpoint: no checkpoint found")

// State holds the persisted position within a log file.
type State struct {
	// Offset is the byte offset of the last successfully processed line.
	Offset int64 `json:"offset"`
	// LastTime is the timestamp of the last successfully processed log entry.
	LastTime time.Time `json:"last_time"`
	// UpdatedAt records when this checkpoint was written.
	UpdatedAt time.Time `json:"updated_at"`
}

// Store persists and retrieves checkpoint State values.
type Store struct {
	path string
}

// New returns a Store that reads and writes checkpoints to the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes the given State to the store, overwriting any previous value.
func (s *Store) Save(state State) error {
	state.UpdatedAt = time.Now().UTC()

	f, err := os.CreateTemp("", "checkpoint-*.tmp")
	if err != nil {
		return err
	}
	tmpName := f.Name()

	if err := json.NewEncoder(f).Encode(state); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, s.path)
}

// Load reads the persisted State from the store.
// Returns ErrNotFound if the checkpoint file does not exist.
func (s *Store) Load() (State, error) {
	f, err := os.Open(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, ErrNotFound
		}
		return State{}, err
	}
	defer f.Close()

	var state State
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return State{}, err
	}
	return state, nil
}

// Remove deletes the checkpoint file. It is a no-op if the file does not exist.
func (s *Store) Remove() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
