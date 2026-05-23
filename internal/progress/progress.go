// Package progress provides a simple progress reporter for tracking
// log line processing throughput and completion status.
package progress

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Reporter tracks processing progress and emits periodic status updates.
type Reporter struct {
	out      io.Writer
	total    int64
	matched  int64
	skipped  int64
	interval time.Duration
	stop     chan struct{}
}

// New creates a Reporter that writes status lines to out every interval.
// Pass interval=0 to disable periodic reporting.
func New(out io.Writer, interval time.Duration) *Reporter {
	return &Reporter{
		out:      out,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Inc records one processed line. matched indicates whether the line
// fell within the requested time range.
func (r *Reporter) Inc(matched bool) {
	atomic.AddInt64(&r.total, 1)
	if matched {
		atomic.AddInt64(&r.matched, 1)
	} else {
		atomic.AddInt64(&r.skipped, 1)
	}
}

// Start begins background reporting. Call Stop when processing is done.
func (r *Reporter) Start() {
	if r.interval == 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.print()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts background reporting and prints a final summary.
func (r *Reporter) Stop() {
	if r.interval != 0 {
		close(r.stop)
	}
	r.print()
}

// Summary returns the final counts.
func (r *Reporter) Summary() (total, matched, skipped int64) {
	return atomic.LoadInt64(&r.total),
		atomic.LoadInt64(&r.matched),
		atomic.LoadInt64(&r.skipped)
}

func (r *Reporter) print() {
	t, m, s := r.Summary()
	fmt.Fprintf(r.out, "progress: total=%d matched=%d skipped=%d\n", t, m, s)
}
