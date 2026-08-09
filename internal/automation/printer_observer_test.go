package automation

import (
	"context"
	"errors"
	"testing"

	"github.com/restaurantautomation/api/internal/printers"
)

func TestPrinterEventObserverPublishesLifecycleEvents(t *testing.T) {
	engine := NewEngine(1, 10, RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()

	observer := NewPrinterEventObserver(engine)
	ticket := printers.Ticket{OrderID: "ORD-700", Destination: "kitchen"}
	observer.Printing(context.Background(), ticket, 1)
	observer.Printed(context.Background(), ticket, 1)
	observer.Failed(context.Background(), ticket, 3, errors.New("printer offline"))

	events := engine.Events()
	if len(events) != 3 {
		t.Fatalf("events = %#v, want three printer lifecycle events", events)
	}
	if events[0].Type != EventPrintStarted || events[1].Type != EventPrintFinished || events[2].Type != EventPrintFailed {
		t.Fatalf("event types = %#v, want print lifecycle sequence", events)
	}
	if events[2].Attempt != 3 || events[2].Payload["destination"] != "kitchen" || events[2].Payload["error"] != "printer offline" {
		t.Fatalf("failure event = %#v", events[2])
	}
}
