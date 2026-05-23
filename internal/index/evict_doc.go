// Package index provides an on-disk positional index that maps log-line
// timestamps to byte offsets within a log file, enabling fast time-range
// seeks without scanning the entire file.
//
// # Eviction
//
// Index files can become stale when the underlying log file is rotated or
// re-written. EvictStale checks the modification time of a cached index file
// and removes it when it exceeds a configurable MaxAge threshold.
//
// Typical usage:
//
//	evicted, err := index.EvictStale(cachePath, index.DefaultEvictOptions())
//	if err != nil {
//		log.Printf("evict: %v", err)
//	}
//	if evicted {
//		// rebuild the index on next access
//	}
package index
