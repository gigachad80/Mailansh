
package progress

import (
	"fmt"
	"sync"
	"time"

	"mailansh/pkg/color"
)

// Reporter handles progress reporting for long-running operations
type Reporter struct {
	total     int
	processed int
	mu        sync.Mutex
	ticker    *time.Ticker
	done      chan bool
	quiet     bool
}

// NewReporter creates a new progress reporter
func NewReporter(total int, quiet bool) *Reporter {
	pr := &Reporter{
		total:  total,
		ticker: time.NewTicker(2 * time.Second),
		done:   make(chan bool),
		quiet:  quiet,
	}
	if !quiet {
		go pr.start()
	}
	return pr
}

// start begins the progress reporting loop
func (pr *Reporter) start() {
	for {
		select {
		case <-pr.ticker.C:
			pr.mu.Lock()
			if pr.total > 0 && !pr.quiet {
				percentage := float64(pr.processed) / float64(pr.total) * 100
				fmt.Printf("\r%s[PROGRESS] %d/%d commits (%.1f%%)%s", color.Blue, pr.processed, pr.total, percentage, color.Reset)
			}
			pr.mu.Unlock()
		case <-pr.done:
			pr.ticker.Stop()
			return
		}
	}
}

// Add increments the processed counter
func (pr *Reporter) Add(count int) {
	pr.mu.Lock()
	pr.processed += count
	pr.mu.Unlock()
}

// Stop stops the progress reporter
func (pr *Reporter) Stop() {
	if !pr.quiet {
		close(pr.done)
		fmt.Print("\n")
	}
}
