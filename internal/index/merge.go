package index

import "sort"

// Merge combines two Index values into a single sorted Index.
// Duplicate timestamps are deduplicated; the entry with the larger offset
// wins so that the merged index always points at the latest known position
// for a given timestamp.
func Merge(a, b Index) Index {
	if len(a.Entries) == 0 {
		return b
	}
	if len(b.Entries) == 0 {
		return a
	}

	// Collect all entries into a map keyed by timestamp (UnixNano) so that
	// duplicates are resolved deterministically.
	type key = int64
	seen := make(map[key]Entry, len(a.Entries)+len(b.Entries))

	for _, e := range a.Entries {
		seen[e.Time.UnixNano()] = e
	}
	for _, e := range b.Entries {
		k := e.Time.UnixNano()
		if existing, ok := seen[k]; !ok || e.Offset > existing.Offset {
			seen[k] = e
		}
	}

	merged := make([]Entry, 0, len(seen))
	for _, e := range seen {
		merged = append(merged, e)
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Time.Before(merged[j].Time)
	})

	return Index{Entries: merged}
}
