package index

import (
	"os"
	"time"
)

// EvictOptions controls the behaviour of EvictStale.
type EvictOptions struct {
	// MaxAge is the maximum age of an index file before it is considered stale.
	// Files with a modification time older than MaxAge will be removed.
	MaxAge time.Duration
}

// DefaultEvictOptions returns sensible defaults for index eviction.
func DefaultEvictOptions() EvictOptions {
	return EvictOptions{
		MaxAge: 72 * time.Hour,
	}
}

// EvictStale removes the index file at path if it is older than opts.MaxAge
// relative to now. It returns true when the file was removed, false when it
// was kept or did not exist. Any unexpected filesystem error is returned.
func EvictStale(path string, now time.Time, opts EvictOptions) (bool, error) {
	if opts.MaxAge <= 0 {
		return false, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if now.Sub(info.ModTime()) <= opts.MaxAge {
		return false, nil
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return false, err
	}

	return true, nil
}
