package automation

import (
	"context"
	"errors"
	"testing"
)

func TestQueueRejectsFullAndClosed(t *testing.T) {
	q := NewQueue(1)
	if err := q.Enqueue(context.Background(), Job{ID: "one"}); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(q.Enqueue(context.Background(), Job{ID: "two"}), ErrQueueFull) {
		t.Fatal("expected full queue error")
	}
	q.Close()
	if !errors.Is(q.Enqueue(context.Background(), Job{ID: "three"}), ErrQueueClosed) {
		t.Fatal("expected closed queue error")
	}
}
