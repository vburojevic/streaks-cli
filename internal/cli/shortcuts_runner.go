package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"streaks-cli/internal/shortcuts"
)

type runResult struct {
	Output   []byte
	Attempts int
	Duration time.Duration
}

func runShortcutWithRetry(ctx context.Context, name string, input []byte, opts *rootOptions) (runResult, error) {
	start := time.Now()
	if opts == nil {
		out, err := runShortcut(ctx, name, input, shortcuts.RunOptions{OutputType: shortcutsOutputType(opts)})
		return runResult{Output: out, Attempts: 1, Duration: time.Since(start)}, err
	}
	attempts := opts.retries + 1
	wait := opts.retryWait
	if wait <= 0 {
		wait = time.Second
	}
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		out, err := runShortcut(ctx, name, input, shortcuts.RunOptions{OutputType: shortcutsOutputType(opts)})
		if err == nil {
			return runResult{Output: out, Attempts: attempt, Duration: time.Since(start)}, nil
		}
		lastErr = err
		if attempt < attempts {
			select {
			case <-ctx.Done():
				return runResult{Attempts: attempt, Duration: time.Since(start)}, ctx.Err()
			case <-time.After(wait):
			}
			wait *= 2
		}
	}
	if lastErr == nil {
		lastErr = errors.New("shortcuts run failed")
	}
	return runResult{Attempts: attempts, Duration: time.Since(start)}, fmt.Errorf("shortcuts run failed after %d attempts: %w", attempts, lastErr)
}

func shortcutsOutputType(opts *rootOptions) string {
	if opts == nil {
		return "public.plain-text"
	}
	if strings.TrimSpace(opts.shortcutsOutput) == "" {
		return "public.plain-text"
	}
	return strings.TrimSpace(opts.shortcutsOutput)
}
