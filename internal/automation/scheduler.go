package automation

import (
	"context"
	"sync"
	"time"
)

type Scheduler struct {
	mu      sync.Mutex
	stopped bool
	timers  []*time.Timer
}

func NewScheduler() *Scheduler { return &Scheduler{} }
func (s *Scheduler) Schedule(ctx context.Context, delay time.Duration, job func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return
	}
	s.timers = append(s.timers, time.AfterFunc(delay, func() {
		select {
		case <-ctx.Done():
		default:
			job()
		}
	}))
}
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopped = true
	for _, timer := range s.timers {
		timer.Stop()
	}
	s.timers = nil
}
