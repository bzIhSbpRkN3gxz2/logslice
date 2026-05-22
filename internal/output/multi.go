package output

import "fmt"

// MultiWriter fans out writes to multiple Writer instances simultaneously.
// It is useful when a single run must produce several output files (e.g.
// splitting a log into hourly shards).
type MultiWriter struct {
	writers []*Writer
}

// NewMulti creates a MultiWriter that writes to all provided Writers.
func NewMulti(writers ...*Writer) *MultiWriter {
	w := make([]*Writer, len(writers))
	copy(w, writers)
	return &MultiWriter{writers: w}
}

// WriteLine writes line to every underlying Writer.
// If any writer fails the error is returned immediately; subsequent writers
// in the slice are still attempted so that partial failures are minimised.
func (m *MultiWriter) WriteLine(line []byte) error {
	var first error
	for _, w := range m.writers {
		if err := w.WriteLine(line); err != nil && first == nil {
			first = fmt.Errorf("output: multiwriter: %w", err)
		}
	}
	return first
}

// Close closes all underlying Writers and returns the first error encountered.
func (m *MultiWriter) Close() error {
	var first error
	for _, w := range m.writers {
		if err := w.Close(); err != nil && first == nil {
			first = fmt.Errorf("output: multiwriter close: %w", err)
		}
	}
	return first
}

// Writers returns the slice of underlying Writers.
func (m *MultiWriter) Writers() []*Writer { return m.writers }
