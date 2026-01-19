package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
	"streaks-cli/internal/shortcuts"
)

type actionCmdOptions struct {
	task   string
	status string
	input  string
	dryRun bool
	stdin  bool
	trace  string
}

var runShortcut = shortcuts.Run
var loadConfig = config.Load
var shortcutExists = shortcuts.Exists

func addActionCommands(root *cobra.Command, defs []discovery.ActionDef, opts *rootOptions) {
	for _, def := range defs {
		if def.Transport != discovery.TransportShortcuts {
			continue
		}
		def := def
		cmdOpts := &actionCmdOptions{}
		cmd := &cobra.Command{
			Use:   def.ID,
			Short: def.Title,
			RunE: func(_ *cobra.Command, _ []string) error {
				return runActionCommand(context.Background(), def, cmdOpts, opts)
			},
		}
		cmd.Flags().StringVar(&cmdOpts.input, "input", "", "Raw input to pass to the wrapper (overrides --task/--status)")
		cmd.Flags().BoolVar(&cmdOpts.stdin, "stdin", false, "Force reading input from stdin")
		cmd.Flags().BoolVar(&cmdOpts.dryRun, "dry-run", false, "Print wrapper and payload without running")
		cmd.Flags().StringVar(&cmdOpts.trace, "trace", "", "Append JSON trace of input/output to a file")
		if def.RequiresTask {
			cmd.Flags().StringVar(&cmdOpts.task, "task", "", "Task name")
		}
		if len(def.ParamOptions) > 0 {
			cmd.Flags().StringVar(&cmdOpts.status, "status", "", "Status value for the action")
		}
		root.AddCommand(cmd)
	}
}

func runActionCommand(ctx context.Context, def discovery.ActionDef, cmdOpts *actionCmdOptions, opts *rootOptions) error {
	cfg, _, err := loadConfig(discovery.DefaultActionDefinitions())
	if err != nil {
		return err
	}
	wrapperName := cfg.Wrappers[def.ID]
	if wrapperName == "" {
		wrapperName = config.WrapperName(cfg.WrapperPrefix, def.ID)
	}
	exists, err := shortcutExists(ctx, wrapperName)
	if err != nil {
		return err
	}
	if !exists && cfg.WrapperPrefix == config.DefaultWrapperPrefix {
		legacy := config.WrapperName(config.LegacyWrapperPrefix, def.ID)
		legacyExists, legacyErr := shortcutExists(ctx, legacy)
		if legacyErr == nil && legacyExists {
			wrapperName = legacy
			exists = true
		}
	}
	if !exists {
		return exitError(ExitCodeWrappersMissing, fmt.Errorf("wrapper shortcut not found: %s", wrapperName))
	}
	input, err := buildActionInput(def, cmdOpts)
	if err != nil {
		return exitError(ExitCodeUsage, err)
	}
	if cmdOpts.dryRun {
		return printDryRun(opts, wrapperName, input)
	}

	if opts != nil && opts.noOutput {
		_, err := runShortcutWithRetry(ctx, wrapperName, input, opts)
		if err != nil {
			return exitError(ExitCodeActionFailed, err)
		}
		return nil
	}

	ctxRun := ctx
	if opts != nil && opts.timeout > 0 {
		var cancel context.CancelFunc
		ctxRun, cancel = context.WithTimeout(ctx, opts.timeout)
		defer cancel()
	}

	out, err := runShortcutWithRetry(ctxRun, wrapperName, input, opts)
	if err != nil {
		_ = appendTrace(cmdOpts.trace, traceEntry{Wrapper: wrapperName, Input: input, Error: err.Error()})
		return exitError(ExitCodeActionFailed, err)
	}
	_ = appendTrace(cmdOpts.trace, traceEntry{Wrapper: wrapperName, Input: input, Output: out})
	if opts.isJSON() || opts.isPlain() {
		var payload any
		if err := json.Unmarshal(out, &payload); err != nil {
			return exitError(ExitCodeActionFailed, fmt.Errorf("wrapper output is not valid JSON: %w", err))
		}
		if opts.isPlain() {
			return output.PrintJSON(os.Stdout, payload, false)
		}
		return output.PrintJSON(os.Stdout, payload, opts.pretty)
	}
	_, err = fmt.Fprint(os.Stdout, string(out))
	return err
}

func buildActionInput(def discovery.ActionDef, cmdOpts *actionCmdOptions) ([]byte, error) {
	if cmdOpts.input != "" {
		return []byte(cmdOpts.input), nil
	}

	if cmdOpts.stdin || !isTTY(os.Stdin) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			return data, nil
		}
	}

	payload := map[string]any{}
	if def.RequiresTask {
		if strings.TrimSpace(cmdOpts.task) == "" {
			return nil, errors.New("missing --task (or provide JSON via --input or stdin)")
		}
		payload["task"] = strings.TrimSpace(cmdOpts.task)
	}
	if cmdOpts.status != "" {
		payload["status"] = cmdOpts.status
	}
	if len(payload) == 0 {
		return nil, nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func printDryRun(opts *rootOptions, wrapper string, input []byte) error {
	payload := map[string]any{
		"dry_run": true,
		"wrapper": wrapper,
	}
	if input != nil {
		payload["input"] = json.RawMessage(input)
	}
	if opts != nil {
		if opts.isJSON() {
			return output.PrintJSON(os.Stdout, payload, opts.pretty)
		}
		if opts.isPlain() {
			return output.PrintJSON(os.Stdout, payload, false)
		}
	}
	if input != nil {
		fmt.Fprintf(os.Stdout, "Dry run: %s %s\n", wrapper, string(input))
		return nil
	}
	fmt.Fprintf(os.Stdout, "Dry run: %s\n", wrapper)
	return nil
}
