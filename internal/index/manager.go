// Package index provides offset-based indexing for fast log range queries.
package index

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Manager coordinates index lifecycle: build, cache, load, and eviction.
type Manager struct {
	cacheDir string
	mu       sync.Mutex
	cache    map[string]*Index
	maxAge   time.Duration
}

// NewManager creates a Manager that stores index files under cacheDir.
func NewManager(cacheDir string, maxAge time.Duration) *Manager {
	return &Manager{
		cacheDir: cacheDir,
		cache:    make(map[string]*Index),
		maxAge:   maxAge,
	}
}

// indexPath returns the canonical path for the index of logFile.
func (m *Manager) indexPath(logFile string) string {
	base := filepath.Base(logFile)
	return filepath.Join(m.cacheDir, base+".idx")
}

// Get returns a cached or freshly built index for logFile.
// If the index file exists and is fresh it is loaded; otherwise it is rebuilt.
func (m *Manager) Get(logFile string, opts BuildOptions) (*Index, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if idx, ok := m.cache[logFile]; ok {
		return idx, nil
	}

	idxPath := m.indexPath(logFile)
	if info, err := os.Stat(idxPath); err == nil {
		age := time.Since(info.ModTime())
		if m.maxAge == 0 || age <= m.maxAge {
			idx, err := Load(idxPath)
			if err == nil {
				m.cache[logFile] = idx
				return idx, nil
			}
		}
	}

	f, err := os.Open(logFile)
	if err != nil {
		return nil, fmt.Errorf("manager: open log file: %w", err)
	}
	defer f.Close()

	idx, err := Build(f, opts)
	if err != nil {
		return nil, fmt.Errorf("manager: build index: %w", err)
	}

	if err := os.MkdirAll(m.cacheDir, 0o755); err == nil {
		_ = Save(idxPath, idx) // best-effort persist
	}

	m.cache[logFile] = idx
	return idx, nil
}

// Invalidate removes the in-memory and on-disk index for logFile.
func (m *Manager) Invalidate(logFile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cache, logFile)
	idxPath := m.indexPath(logFile)
	if err := os.Remove(idxPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("manager: invalidate: %w", err)
	}
	return nil
}
