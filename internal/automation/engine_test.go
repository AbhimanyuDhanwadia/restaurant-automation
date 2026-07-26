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
