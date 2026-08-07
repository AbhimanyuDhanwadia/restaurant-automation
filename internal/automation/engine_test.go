package automation

import (
	"context"
	"testing"
	"time"
)

func TestEngineProcessesOrderPipeline(t *testing.T) {
	engine := NewEngine(1, 10, RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	if err := engine.SubmitOrder(context.Background(), "ORD-100"); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for len(engine.Events()) < 5 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if len(engine.Events()) != 5 {
		t.Fatalf("got %d events, want 5", len(engine.Events()))
	}
	if got := engine.Events()[4].Type; got != EventOrderQueued {
		t.Fatalf("last event = %s", got)
	}
}

func TestEngineRecordsOrderLifecycleEvents(t *testing.T) {
	engine := NewEngine(1, 10, RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()

	cases := []struct {
		status string
		want   EventType
	}{
		{status: "preparing", want: EventKitchenAccepted},
		{status: "ready", want: EventOrderReady},
		{status: "delivered", want: EventOrderDelivered},
		{status: "cancelled", want: EventOrderCancelled},
	}
	for _, test := range cases {
		if err := engine.RecordOrderStatus(context.Background(), "ORD-200", test.status); err != nil {
			t.Fatal(err)
		}
	}
	if err := engine.RecordOrderStatus(context.Background(), "ORD-200", "received"); err != nil {
		t.Fatal(err)
	}

	events := engine.Events()
	if len(events) != len(cases) {
		t.Fatalf("events = %#v, want %d lifecycle events", events, len(cases))
	}
	for index, test := range cases {
		if events[index].Type != test.want || events[index].Payload["status"] != test.status {
			t.Fatalf("event %d = %#v, want %s for %s", index, events[index], test.want, test.status)
		}
	}
}
