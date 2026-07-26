package automation

import (
	"context"
	"errors"
	"sync"
)

var ErrQueueClosed = errors.New("automation queue is closed")
var ErrQueueFull = errors.New("automation queue is full")

type Job struct {
	ID          string
	Name        string
	OrderID     string
	Attempt     int
	MaxAttempts int
	Run         func(context.Context) error
}

type Queue struct {
	mu     sync.RWMutex
	jobs   chan Job
	closed bool
}

func NewQueue(capacity int) *Queue {
	if capacity < 1 {
		capacity = 1
	}
	return &Queue{jobs: make(chan Job, capacity)}
}
func (q *Queue) Enqueue(ctx context.Context, job Job) error {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.closed {
		return ErrQueueClosed
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.jobs <- job:
		return nil
	default:
		return ErrQueueFull
	}
}
func (q *Queue) Dequeue(ctx context.Context) (Job, error) {
	select {
	case <-ctx.Done():
		return Job{}, ctx.Err()
	case job, ok := <-q.jobs:
		if !ok {
			return Job{}, ErrQueueClosed
		}
		return job, nil
	}
}
func (q *Queue) Len() int      { return len(q.jobs) }
func (q *Queue) Capacity() int { return cap(q.jobs) }
func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	close(q.jobs)
}
