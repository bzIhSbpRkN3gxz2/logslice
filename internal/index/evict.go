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
	// MaxAge is the maximum age of an index file before it is removed.
	// A zero value disables eviction.
	MaxAge time.Duration
}

// EvictStale removes the index file at path if it is older than opts.MaxAge.
// If the file does not exist the call is a no-op.
// If opts.MaxAge is zero the call is a no-op.
func EvictStale(path string, opts EvictOptions) error {
	if opts.MaxAge == 0 {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if time.Since(info.ModTime()) > opts.MaxAge {
		return os.Remove(path)
	}
	return nil
}
