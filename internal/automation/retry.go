package automation

import (
	"context"
	"time"
)

type RetryPolicy struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

func (p RetryPolicy) normalized() RetryPolicy {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = 30 * time.Second
	}
	return p
}

func ExecuteWithRetry(ctx context.Context, job Job, policy RetryPolicy, onRetry func(int, error)) error {
	policy = policy.normalized()
	var err error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		err = job.Run(ctx)
		if err == nil {
			return nil
		}
		if attempt == policy.MaxAttempts {
			break
		}
		delay := policy.InitialDelay
		for n := 1; n < attempt; n++ {
			delay *= 2
		}
		if delay > policy.MaxDelay {
			delay = policy.MaxDelay
		}
		if onRetry != nil {
			onRetry(attempt+1, err)
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}
