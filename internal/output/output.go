// Package output handles writing filtered log lines to destinations.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Writer wraps a destination and writes log lines to it.
type Writer struct {
	w       *bufio.Writer
	closer  io.Closer
	count   int64
}

// Options configures an output Writer.
type Options struct {
	// Path is the file path to write to. Use "-" for stdout.
	Path string
	// BufferSize overrides the default 64 KiB write buffer.
	BufferSize int
}

// New creates a Writer for the given options.
func New(opts Options) (*Writer, error) {
	const defaultBuf = 64 * 1024

	bufSize := opts.BufferSize
	if bufSize <= 0 {
		bufSize = defaultBuf
	}

	var (
		w      io.Writer
		closer io.Closer
	)

	if opts.Path == "" || opts.Path == "-" {
		w = os.Stdout
		closer = io.NopCloser(os.Stdout)
	} else {
		f, err := os.Create(opts.Path)
		if err != nil {
			return nil, fmt.Errorf("output: create %q: %w", opts.Path, err)
		}
		w = f
		closer = f
	}

	return &Writer{
		w:      bufio.NewWriterSize(w, bufSize),
		closer: closer,
	}, nil
}

// WriteLine writes a single raw log line followed by a newline.
func (w *Writer) WriteLine(line []byte) error {
	if _, err := w.w.Write(line); err != nil {
		return err
	}
	if err := w.w.WriteByte('\n'); err != nil {
		return err
	}
	w.count++
	return nil
}

// Count returns the number of lines written so far.
func (w *Writer) Count() int64 { return w.count }

// Close flushes the buffer and closes the underlying destination.
func (w *Writer) Close() error {
	if err := w.w.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	return w.closer.Close()
}
