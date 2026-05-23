/*
Package index — Manager

Manager provides a high-level coordinator for the index lifecycle.
It abstracts build, persist, load, and eviction behind a single Get call
so that callers (e.g. the splitter) never need to reason about whether an
index already exists on disk.

Usage:

	m := index.NewManager("/var/cache/logslice", 24*time.Hour)

	// Obtain an index, building it if necessary.
	idx, err := m.Get("/var/log/app.log", index.DefaultBuildOptions())

	// Search the index for a time range.
	start, end := index.Search(idx, &from, &to)

	// Discard a stale index so the next Get rebuilds it.
	_ = m.Invalidate("/var/log/app.log")

Thread safety:

All public methods on Manager are safe for concurrent use; an internal
mutex serialises cache reads and writes.
*/
package index
