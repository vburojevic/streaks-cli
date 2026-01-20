package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

import "streaks-cli/internal/output"

func runNamedShortcut(ctx context.Context, actionID, name string, input []byte, cmdOpts *actionCmdOptions, opts *rootOptions) error {
	result, err := runShortcutOnce(ctx, name, input, opts)
	if err != nil {
		_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Error: err.Error()})
		if isShortcutNotFound(err) {
			return exitError(ExitCodeShortcutMissing, err)
		}
		return exitError(ExitCodeActionFailed, err)
	}
	_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Output: result.Output})
	return emitActionOutput(actionID, name, input, result, opts)
}

func runCandidateShortcuts(ctx context.Context, candidates []string, actionID string, input []byte, cmdOpts *actionCmdOptions, opts *rootOptions) error {
	if len(candidates) == 0 {
		return exitError(ExitCodeShortcutMissing, fmt.Errorf("no matching Streaks shortcut found for action %s", actionID))
	}
	for _, name := range candidates {
		result, err := runShortcutOnce(ctx, name, input, opts)
		if err != nil {
			if isShortcutNotFound(err) {
				continue
			}
			_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Error: err.Error()})
			return exitError(ExitCodeActionFailed, err)
		}
		_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Output: result.Output})
		return emitActionOutput(actionID, name, input, result, opts)
	}
	return exitError(ExitCodeShortcutMissing, fmt.Errorf("no matching Streaks shortcut found for action %s; expected one of: %s", actionID, strings.Join(candidates, ", ")))
}

func runShortcutOnce(ctx context.Context, name string, input []byte, opts *rootOptions) (runResult, error) {
	if opts != nil && opts.noOutput {
		result, err := runShortcutWithRetry(ctx, name, input, opts)
		result.Output = nil
		return result, err
	}
	ctxRun := ctx
	if opts != nil && opts.timeout > 0 {
		var cancel context.CancelFunc
		ctxRun, cancel = context.WithTimeout(ctx, opts.timeout)
		defer cancel()
	}
	return runShortcutWithRetry(ctxRun, name, input, opts)
}

func emitActionOutput(actionID, shortcutName string, input []byte, result runResult, opts *rootOptions) error {
	if opts != nil && opts.noOutput {
		return nil
	}
	if opts != nil && opts.isAgent() {
		return output.PrintJSON(os.Stdout, buildActionEnvelope(actionID, shortcutName, input, result), false)
	}
	_, err := fmt.Fprint(os.Stdout, string(result.Output))
	return err
}

func isShortcutNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "find shortcut") {
		return true
	}
	if strings.Contains(msg, "couldn’t find shortcut") || strings.Contains(msg, "couldn't find shortcut") {
		return true
	}
	return false
}

func normalizeShortcutOutput(out []byte, shortcutName string) any {
	var payload any
	if err := json.Unmarshal(out, &payload); err != nil {
		payload = map[string]any{
			"raw":      strings.TrimSpace(string(out)),
			"format":   "text",
			"shortcut": shortcutName,
		}
	}
	return payload
}

type actionEnvelope struct {
	OK         bool               `json:"ok"`
	Timestamp  string             `json:"timestamp"`
	Action     actionEnvelopeInfo `json:"action"`
	Shortcut   actionShortcutInfo `json:"shortcut"`
	Attempts   int                `json:"attempts"`
	DurationMS int64              `json:"duration_ms"`
	Input      any                `json:"input,omitempty"`
	Result     any                `json:"result,omitempty"`
}

type actionEnvelopeInfo struct {
	ID string `json:"id"`
}

type actionShortcutInfo struct {
	Name string `json:"name"`
}

func buildActionEnvelope(actionID, shortcutName string, input []byte, result runResult) actionEnvelope {
	envelope := actionEnvelope{
		OK:         true,
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Action:     actionEnvelopeInfo{ID: actionID},
		Shortcut:   actionShortcutInfo{Name: shortcutName},
		Attempts:   result.Attempts,
		DurationMS: result.Duration.Milliseconds(),
		Result:     normalizeShortcutOutput(result.Output, shortcutName),
	}
	if len(input) > 0 {
		envelope.Input = normalizeInput(input)
	}
	return envelope
}

func normalizeInput(input []byte) any {
	var payload any
	if err := json.Unmarshal(input, &payload); err == nil {
		return payload
	}
	return strings.TrimSpace(string(input))
}
