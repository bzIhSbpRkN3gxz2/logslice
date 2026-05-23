package index

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager coordinates building, caching, loading, and evicting index files
// for a given log file path.
type Manager struct {
	cacheDir string
	buildOpts BuildOptions
	evictOpts EvictOptions
}

// NewManager creates a Manager that stores index files under cacheDir.
func NewManager(cacheDir string, build BuildOptions, evict EvictOptions) *Manager {
	return &Manager{
		cacheDir:  cacheDir,
		buildOpts: build,
		evictOpts: evict,
	}
}

// indexPath returns the cache path for the given log file.
func (m *Manager) indexPath(logPath string) string {
	base := filepath.Base(logPath) + ".idx"
	return filepath.Join(m.cacheDir, base)
}

// Get returns a loaded Index for logPath, rebuilding it when stale or absent.
func (m *Manager) Get(logPath string) (*Index, error) {
	idxPath := m.indexPath(logPath)

	evicted, err := EvictStale(idxPath, m.evictOpts)
	if err != nil {
		return nil, fmt.Errorf("evict index: %w", err)
	}

	if !evicted {
		if idx, err := Load(idxPath); err == nil {
			return idx, nil
		}
	}

	return m.rebuild(logPath, idxPath)
}

// Invalidate removes the cached index for logPath, if present.
func (m *Manager) Invalidate(logPath string) error {
	p := m.indexPath(logPath)
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("invalidate index: %w", err)
	}
	return nil
}

func (m *Manager) rebuild(logPath, idxPath string) (*Index, error) {
	f, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("open log for indexing: %w", err)
	}
	defer f.Close()

	if err := os.MkdirAll(m.cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}

	idx, err := Build(f, m.buildOpts)
	if err != nil {
		return nil, fmt.Errorf("build index: %w", err)
	}

	if err := Save(idxPath, idx); err != nil {
		return nil, fmt.Errorf("save index: %w", err)
	}

	return idx, nil
}
