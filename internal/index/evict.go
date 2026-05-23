package index

import (
	"os"
	"time"
)

// DefaultEvictOptions returns EvictOptions with sensible defaults.
func DefaultEvictOptions() EvictOptions {
	return EvictOptions{
		MaxAge: 72 * time.Hour,
	}
}

// EvictOptions controls the behaviour of EvictStale.
type EvictOptions struct {
	// MaxAge is the maximum age of an index file before it is considered stale.
	// A zero value disables eviction.
	MaxAge time.Duration
}

// EvictStale removes the index file at path if it is older than opts.MaxAge.
// It returns (true, nil) when the file was removed, (false, nil) when it was
// kept or did not exist, and (false, err) on any unexpected error.
func EvictStale(path string, opts EvictOptions) (evicted bool, err error) {
	if opts.MaxAge == 0 {
		return false, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if time.Since(info.ModTime()) <= opts.MaxAge {
		return false, nil
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
