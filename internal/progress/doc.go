// Package progress provides lightweight, concurrency-safe progress reporting
// for logslice processing pipelines.
//
// Usage:
//
//	rep := progress.New(os.Stderr, 2*time.Second)
//	rep.Start()
//	defer rep.Stop()
//
//	for _, line := range lines {
//		matched := filter.Match(line)
//		rep.Inc(matched)
//	}
//
// When interval is 0 no background goroutine is started; Stop() still
// prints a final one-line summary to the provided writer.
package progress
