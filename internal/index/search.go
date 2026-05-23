package index

import (
	"time"
)

// SearchResult holds the byte offsets bounding a time range query.
type SearchResult struct {
	// StartOffset is the byte offset of the first line at or after Start.
	StartOffset int64
	// EndOffset is the byte offset just past the last line before or at End.
	// A value of -1 means "read to end of file".
	EndOffset int64
}

// Search returns the byte-offset window within the indexed file that may
// contain log lines falling in [start, end]. Either bound may be zero to
// indicate "unbounded".
//
// The returned offsets are conservative: callers must still filter individual
// lines, but they can seek directly to StartOffset and stop reading at
// EndOffset to avoid scanning the whole file.
func Search(idx *Index, start, end time.Time) SearchResult {
	var result SearchResult

	if !start.IsZero() {
		result.StartOffset = FloorOffset(idx, start)
	} else {
		result.StartOffset = 0
	}

	if !end.IsZero() {
		// Use the entry just after end as the ceiling so we don't miss lines
		// whose timestamp equals end exactly.
		result.EndOffset = ceilOffset(idx, end)
	} else {
		result.EndOffset = -1
	}

	return result
}

// ceilOffset returns the byte offset of the first index entry whose timestamp
// is strictly after t, or -1 if no such entry exists (meaning read to EOF).
func ceilOffset(idx *Index, t time.Time) int64 {
	if len(idx.Entries) == 0 {
		return -1
	}
	for _, e := range idx.Entries {
		if e.Time.After(t) {
			return e.Offset
		}
	}
	return -1
}
