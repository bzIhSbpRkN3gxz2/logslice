package index

import (
	"fmt"
	"sync"
	"time"
)

// Manager caches built indexes in memory and handles eviction of stale
// on-disk index files before loading them.
type Manager struct {
	mu      sync.Mutex
	cache   map[string]*Index
	evict   EvictOptions
	build   BuildOptions
}

// NewManager creates a Manager with the provided options.
func NewManager(build BuildOptions, evict EvictOptions) *Manager {
	return &Manager{
		cache: make(map[string]*Index),
		evict: evict,
		build: build,
	}
}

// Get returns a cached or freshly built Index for the given log file.
// It evicts a stale on-disk index (at idxPath) before attempting to load it.
func (m *Manager) Get(logPath, idxPath string) (*Index, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if idx, ok := m.cache[logPath]; ok {
		return idx, nil
	}

	// Evict stale index from disk before loading.
	if err := EvictStale(idxPath, m.evict); err != nil {
		return nil, fmt.Errorf("manager: evict %s: %w", idxPath, err)
	}

	idx, err := Load(idxPath)
	if err != nil {
		// Fall back to building a new index.
		idx, err = Build(logPath, m.build)
		if err != nil {
			return nil, fmt.Errorf("manager: build index for %s: %w", logPath, err)
		}
		if saveErr := Save(idxPath, idx); saveErr != nil {
			// Non-fatal: log the error but continue with the in-memory index.
			_ = saveErr
		}
	}

	m.cache[logPath] = idx
	return idx, nil
}

// Invalidate removes the in-memory cached index for logPath.
func (m *Manager) Invalidate(logPath string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, logPath)
}

// CacheSize returns the number of indexes currently held in memory.
func (m *Manager) CacheSize() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.cache)
}

// defaultManagerMaxAge is used when constructing a Manager via NewManager with
// default eviction settings.
const defaultManagerMaxAge = 72 * time.Hour
