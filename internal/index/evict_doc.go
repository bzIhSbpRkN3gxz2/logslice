// Package index provides sparse positional indexes over structured log files.
//
// # Eviction
//
// EvictStale removes index files that have exceeded a configurable maximum age.
// This prevents stale indexes from being used after a log file has been rotated
// or truncated.
//
// Typical usage:
//
//	opts := index.DefaultEvictOptions()
//	if err := index.EvictStale("/var/cache/logslice/app.log.idx", opts); err != nil {
//		log.Printf("evict: %v", err)
//	}
package index
