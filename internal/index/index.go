// Package index provides byte-offset indexing for large log files,
// enabling fast seeks to approximate time positions without full scans.
package index

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

// Entry records a timestamp and the byte offset of the line that carries it.
type Entry struct {
	Timestamp time.Time `json:"ts"`
	Offset    int64     `json:"offset"`
}

// Index is an ordered collection of sampled log-line positions.
type Index struct {
	Entries []Entry `json:"entries"`
}

// Add appends a new entry. Callers must add entries in ascending offset order.
func (idx *Index) Add(ts time.Time, offset int64) {
	idx.Entries = append(idx.Entries, Entry{Timestamp: ts, Offset: offset})
}

// FloorOffset returns the byte offset of the last entry whose timestamp is
// less than or equal to target. Returns 0 if no such entry exists.
func (idx *Index) FloorOffset(target time.Time) int64 {
	if len(idx.Entries) == 0 {
		return 0
	}
	i := sort.Search(len(idx.Entries), func(i int) bool {
		return idx.Entries[i].Timestamp.After(target)
	})
	if i == 0 {
		return 0
	}
	return idx.Entries[i-1].Offset
}

// Save serialises the index to path as JSON.
func (idx *Index) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(idx)
}

// Load deserialises an index from path.
func Load(path string) (*Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var idx Index
	if err := json.NewDecoder(f).Decode(&idx); err != nil {
		return nil, err
	}
	return &idx, nil
}
