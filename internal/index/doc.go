// Package index provides a sparse offset index for large log files.
//
// An Index maps sampled timestamps to byte offsets within a log file,
// enabling fast seeking to approximate positions before a linear scan.
//
// # Building an Index
//
// Use Build to create an index from an io.ReadSeeker:
//
//	idx, err := index.Build(f, index.DefaultBuildOptions())
//
// # Querying an Index
//
// FloorOffset returns the largest byte offset whose timestamp is less than
// or equal to the query time, giving a safe seek position:
//
//	offset, ok := idx.FloorOffset(t)
//
// # Persistence
//
// Indexes can be saved to and loaded from disk:
//
//	err = idx.Save(path)
//	idx, err = index.Load(path)
package index
