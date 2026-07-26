package automation

import (
	"context"
	"sync"
)

type EventBus struct {
	mu          sync.RWMutex
	events      []Event
	subscribers map[int]chan Event
	nextID      int
	closed      bool
}

func NewEventBus() *EventBus { return &EventBus{subscribers: make(map[int]chan Event)} }

func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return ErrQueueClosed
	}
	b.events = append(b.events, event)
	channels := make([]chan Event, 0, len(b.subscribers))
	for _, subscriber := range b.subscribers {
		channels = append(channels, subscriber)
	}
	b.mu.Unlock()
	for _, subscriber := range channels {
		select {
		case subscriber <- event:
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

func (b *EventBus) Subscribe(buffer int) (<-chan Event, func()) {
	if buffer < 1 {
		buffer = 1
	}
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		closed := make(chan Event)
		close(closed)
		return closed, func() {}
	}
	id := b.nextID
	b.nextID++
	channel := make(chan Event, buffer)
	b.subscribers[id] = channel
	b.mu.Unlock()
	return channel, func() {
		b.mu.Lock()
		if existing, ok := b.subscribers[id]; ok {
			delete(b.subscribers, id)
			close(existing)
		}
		b.mu.Unlock()
	}
}

func (b *EventBus) Snapshot() []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]Event(nil), b.events...)
}

func (b *EventBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for id, subscriber := range b.subscribers {
		close(subscriber)
		delete(b.subscribers, id)
	}
}
