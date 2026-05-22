// Package output provides buffered writers for logslice output destinations.
//
// A Writer can target either a named file or standard output (using "-" or an
// empty path string). All writes are buffered for performance; callers MUST
// call Close to flush pending data and release the underlying file handle.
//
// Typical usage:
//
//	w, err := output.New(output.Options{Path: "filtered.log"})
//	if err != nil { ... }
//	defer w.Close()
//
//	for _, line := range lines {
//		if err := w.WriteLine(line.Raw); err != nil { ... }
//	}
//
//	fmt.Printf("wrote %d lines\n", w.Count())
package output
