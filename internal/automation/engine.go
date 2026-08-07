package automation

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type Stats struct {
	QueueDepth    int    `json:"queue_depth"`
	QueueCapacity int    `json:"queue_capacity"`
	Workers       int    `json:"workers"`
	Events        int    `json:"events"`
	Retried       uint64 `json:"retried"`
	Failed        uint64 `json:"failed"`
}

type Engine struct {
	bus       *EventBus
	queue     *Queue
	scheduler *Scheduler
	policy    RetryPolicy
	workers   int
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	started   bool
	mu        sync.Mutex
	retried   atomic.Uint64
	failed    atomic.Uint64
}

func NewEngine(workers, queueCapacity int, policy RetryPolicy) *Engine {
	if workers < 1 {
		workers = 1
	}
	return &Engine{bus: NewEventBus(), queue: NewQueue(queueCapacity), scheduler: NewScheduler(), policy: policy, workers: workers}
}
func (e *Engine) Start(parent context.Context) {
	e.mu.Lock()
	if e.started {
		e.mu.Unlock()
		return
	}
	e.ctx, e.cancel = context.WithCancel(parent)
	e.started = true
	e.mu.Unlock()
	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go e.worker()
	}
}
func (e *Engine) SubmitOrder(ctx context.Context, orderID string) error {
	if err := e.bus.Publish(ctx, NewEvent(EventOrderReceived, orderID, nil)); err != nil {
		return err
	}
	return e.queue.Enqueue(ctx, Job{ID: uuid.NewString(), Name: "normalize-order", OrderID: orderID, MaxAttempts: e.policy.MaxAttempts, Run: func(jobCtx context.Context) error {
		for _, eventType := range []EventType{EventOrderNormalized, EventOrderValidated, EventOrderStored, EventOrderQueued} {
			if err := e.bus.Publish(jobCtx, NewEvent(eventType, orderID, nil)); err != nil {
				return err
			}
		}
		return nil
	}})
}

// RecordOrderStatus publishes the operational lifecycle event that corresponds
// to a durable order-status transition.
func (e *Engine) RecordOrderStatus(ctx context.Context, orderID, status string) error {
	eventType, ok := lifecycleEventType(status)
	if !ok {
		return nil
	}
	return e.bus.Publish(ctx, NewEvent(eventType, orderID, map[string]string{"status": status}))
}

func lifecycleEventType(status string) (EventType, bool) {
	switch status {
	case "preparing":
		return EventKitchenAccepted, true
	case "ready":
		return EventOrderReady, true
	case "delivered":
		return EventOrderDelivered, true
	case "cancelled":
		return EventOrderCancelled, true
	default:
		return "", false
	}
}

func (e *Engine) worker() {
	defer e.wg.Done()
	for {
		job, err := e.queue.Dequeue(e.ctx)
		if err != nil {
			return
		}
		err = ExecuteWithRetry(e.ctx, job, e.policy, func(attempt int, retryErr error) {
			e.retried.Add(1)
			_ = e.bus.Publish(e.ctx, NewEvent(EventJobRetrying, job.OrderID, map[string]string{"attempt": strconv.Itoa(attempt), "error": retryErr.Error()}))
		})
		if err != nil {
			e.failed.Add(1)
			_ = e.bus.Publish(e.ctx, NewEvent(EventJobFailed, job.OrderID, map[string]string{"error": err.Error()}))
		}
	}
}
func (e *Engine) Events() []Event                             { return e.bus.Snapshot() }
func (e *Engine) Subscribe(buffer int) (<-chan Event, func()) { return e.bus.Subscribe(buffer) }
func (e *Engine) Stats() Stats {
	return Stats{QueueDepth: e.queue.Len(), QueueCapacity: e.queue.Capacity(), Workers: e.workers, Events: len(e.Events()), Retried: e.retried.Load(), Failed: e.failed.Load()}
}
func (e *Engine) Close() {
	e.mu.Lock()
	if !e.started {
		e.mu.Unlock()
		return
	}
	e.cancel()
	e.scheduler.Stop()
	e.queue.Close()
	e.bus.Close()
	e.mu.Unlock()
	e.wg.Wait()
}
