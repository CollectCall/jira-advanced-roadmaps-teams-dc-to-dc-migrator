package app

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type progressTracker struct {
	total   int
	current int
	label   string
	done    chan struct{}
	mu      sync.Mutex
	active  bool
}

func newProgressTracker(total int) *progressTracker {
	if !isInteractiveTerminal() {
		return &progressTracker{total: total}
	}
	return &progressTracker{total: total, active: true}
}

func (p *progressTracker) Start(stepLabel string) {
	if !p.active {
		return
	}
	p.mu.Lock()
	p.current++
	p.label = stepLabel
	done := make(chan struct{})
	p.done = done
	step := p.current
	total := p.total
	p.mu.Unlock()

	go func() {
		frames := []string{"-", "\\", "|", "/"}
		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-done:
				p.render(step, total, stepLabel, "done")
				fmt.Fprintln(os.Stdout)
				return
			case <-ticker.C:
				p.render(step, total, stepLabel, frames[i%len(frames)])
				i++
			}
		}
	}()
}

func (p *progressTracker) End() {
	if !p.active {
		return
	}
	p.mu.Lock()
	done := p.done
	p.done = nil
	p.mu.Unlock()
	if done != nil {
		close(done)
	}
}

func (p *progressTracker) render(step, total int, label, spinner string) {
	width := 20
	filled := 0
	if total > 0 {
		filled = int(float64(step) / float64(total) * float64(width))
	}
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", width-filled)
	fmt.Fprintf(os.Stdout, "\r[%s] %d/%d %s %s", bar, step, total, spinner, label)
}

func (p *progressTracker) Finish() {
	p.End()
}
