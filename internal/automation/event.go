package automation

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventOrderReceived   EventType = "order.received"
	EventOrderNormalized EventType = "order.normalized"
	EventOrderValidated  EventType = "order.validated"
	EventOrderStored     EventType = "order.stored"
	EventOrderQueued     EventType = "order.queued"
	EventJobRetrying     EventType = "job.retrying"
	EventJobFailed       EventType = "job.failed"
)

type Event struct {
	ID        string            `json:"id"`
	Type      EventType         `json:"type"`
	OrderID   string            `json:"order_id"`
	CreatedAt time.Time         `json:"created_at"`
	Attempt   int               `json:"attempt,omitempty"`
	Payload   map[string]string `json:"payload,omitempty"`
}

func NewEvent(eventType EventType, orderID string, payload map[string]string) Event {
	return Event{ID: uuid.NewString(), Type: eventType, OrderID: orderID, CreatedAt: time.Now().UTC(), Payload: payload}
}
