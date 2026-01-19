package cli

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func runShortcutWithRetry(ctx context.Context, name string, input []byte, opts *rootOptions) ([]byte, error) {
	if opts == nil {
		return runShortcut(ctx, name, input)
	}
	attempts := opts.retries + 1
	wait := opts.retryWait
	if wait <= 0 {
		wait = time.Second
	}
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		out, err := runShortcut(ctx, name, input)
		if err == nil {
			return out, nil
		}
		lastErr = err
		if attempt < attempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
			wait *= 2
		}
	}
	if lastErr == nil {
		lastErr = errors.New("shortcuts run failed")
	}
	return nil, fmt.Errorf("shortcuts run failed after %d attempts: %w", attempts, lastErr)
}
