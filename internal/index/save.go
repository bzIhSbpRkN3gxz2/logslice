package index

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// diskEntry is the JSON-serialisable form of a single index entry.
type diskEntry struct {
	Timestamp time.Time `json:"ts"`
	Offset    int64     `json:"off"`
}

// diskIndex is the JSON-serialisable form of an Index.
type diskIndex struct {
	Version int         `json:"version"`
	Entries []diskEntry `json:"entries"`
}

const indexVersion = 1

// Save writes the index to the given file path as JSON.
func (idx *Index) Save(path string) error {
	di := diskIndex{
		Version: indexVersion,
		Entries: make([]diskEntry, len(idx.entries)),
	}
	for i, e := range idx.entries {
		di.Entries[i] = diskEntry{Timestamp: e.Timestamp, Offset: e.Offset}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("index save: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(di); err != nil {
		return fmt.Errorf("index save encode: %w", err)
	}
	return nil
}

// Load reads an index that was previously written by Save.
func Load(path string) (*Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("index load: %w", err)
	}
	defer f.Close()

	var di diskIndex
	if err := json.NewDecoder(f).Decode(&di); err != nil {
		return nil, fmt.Errorf("index load decode: %w", err)
	}
	if di.Version != indexVersion {
		return nil, fmt.Errorf("index load: unsupported version %d", di.Version)
	}

	idx := &Index{entries: make([]Entry, len(di.Entries))}
	for i, e := range di.Entries {
		idx.entries[i] = Entry{Timestamp: e.Timestamp, Offset: e.Offset}
	}
	return idx, nil
}
