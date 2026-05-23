// Package index provides an in-memory sparse index over a log file, mapping
// sampled timestamps to byte offsets. The index allows logslice to seek
// directly into large files rather than scanning from the beginning.
//
// # Building
//
// Use [Build] to create an index from an io.ReadSeeker. The default sampling
// rate (one entry per 1 000 lines) balances memory usage against seek
// precision.
//
// # Persisting
//
// [Save] writes the index to a file and [Load] restores it. The binary format
// is versioned so incompatible changes are detected at load time.
//
// # Searching
//
// [Search] accepts an *Index and a half-open time window [start, end] and
// returns a [SearchResult] with conservative byte-offset bounds. Callers
// should seek to StartOffset and read until EndOffset (or EOF when EndOffset
// is -1), then apply line-level filtering for precise results.
package index
