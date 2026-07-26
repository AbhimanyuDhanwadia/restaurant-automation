package automation

import (
	"context"
	"errors"
	"testing"
)

func TestExecuteWithRetry(t *testing.T) {
	attempts := 0
	err := ExecuteWithRetry(context.Background(), Job{Run: func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary")
		}
		return nil
	}}, RetryPolicy{MaxAttempts: 3}, nil)
	if err != nil || attempts != 3 {
		t.Fatalf("got err=%v attempts=%d", err, attempts)
	}
}
