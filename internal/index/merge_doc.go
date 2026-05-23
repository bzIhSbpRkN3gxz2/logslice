// Package index provides time-based indexing for log files.
//
// # Merge
//
// Merge combines two Index values into a single sorted Index. It is useful
// when an incremental index built over a newly-appended portion of a log file
// must be unified with a previously persisted index.
//
// Behaviour:
//   - The result contains all entries from both inputs, sorted by timestamp.
//   - When both inputs contain an entry for the same timestamp, the entry with
//     the higher byte offset is kept. This ensures the merged index always
//     points at the most recent known position for a given moment in time.
//   - Neither input is modified; a new Index is returned.
//
// Example:
//
//	base, _ := index.Load(idxPath)
//	delta := index.Build(newSegment, index.DefaultBuildOptions())
//	unified := index.Merge(base, delta)
//	_ = index.Save(idxPath, unified)
package index
