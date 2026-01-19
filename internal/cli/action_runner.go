package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

import "streaks-cli/internal/output"

func runNamedShortcut(ctx context.Context, name string, input []byte, cmdOpts *actionCmdOptions, opts *rootOptions) error {
	out, err := runShortcutOnce(ctx, name, input, opts)
	if err != nil {
		_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Error: err.Error()})
		return exitError(ExitCodeActionFailed, err)
	}
	_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Output: out})
	return emitActionOutput(name, out, opts)
}

func runCandidateShortcuts(ctx context.Context, candidates []string, actionID string, input []byte, cmdOpts *actionCmdOptions, opts *rootOptions) error {
	if len(candidates) == 0 {
		return exitError(ExitCodeShortcutMissing, fmt.Errorf("no matching Streaks shortcut found for action %s", actionID))
	}
	for _, name := range candidates {
		out, err := runShortcutOnce(ctx, name, input, opts)
		if err != nil {
			if isShortcutNotFound(err) {
				continue
			}
			_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Error: err.Error()})
			return exitError(ExitCodeActionFailed, err)
		}
		_ = appendTrace(cmdOpts.trace, traceEntry{Shortcut: name, Input: input, Output: out})
		return emitActionOutput(name, out, opts)
	}
	return exitError(ExitCodeShortcutMissing, fmt.Errorf("no matching Streaks shortcut found for action %s; expected one of: %s", actionID, strings.Join(candidates, ", ")))
}

func runShortcutOnce(ctx context.Context, name string, input []byte, opts *rootOptions) ([]byte, error) {
	if opts != nil && opts.noOutput {
		_, err := runShortcutWithRetry(ctx, name, input, opts)
		return nil, err
	}
	ctxRun := ctx
	if opts != nil && opts.timeout > 0 {
		var cancel context.CancelFunc
		ctxRun, cancel = context.WithTimeout(ctx, opts.timeout)
		defer cancel()
	}
	return runShortcutWithRetry(ctxRun, name, input, opts)
}

func emitActionOutput(shortcutName string, out []byte, opts *rootOptions) error {
	if opts != nil && opts.noOutput {
		return nil
	}
	if opts != nil && (opts.isJSON() || opts.isPlain()) {
		var payload any
		if err := json.Unmarshal(out, &payload); err != nil {
			payload = map[string]any{
				"raw":      strings.TrimSpace(string(out)),
				"format":   "text",
				"shortcut": shortcutName,
			}
		}
		if opts.isPlain() {
			return output.PrintJSON(os.Stdout, payload, false)
		}
		return output.PrintJSON(os.Stdout, payload, opts.pretty)
	}
	_, err := fmt.Fprint(os.Stdout, string(out))
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
