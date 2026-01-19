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
	task     string
	status   string
	input    string
	dryRun   bool
	stdin    bool
	trace    string
	shortcut string
}

var runShortcut = shortcuts.RunWithOptions
var loadConfig = config.Load
var shortcutExists = shortcuts.Exists
var listShortcuts = shortcuts.List
var discover = discovery.Discover

type shortcutKind int

const (
	shortcutWrapper shortcutKind = iota
	shortcutDirect
	shortcutExplicit
)

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
		cmd.Flags().StringVar(&cmdOpts.shortcut, "shortcut", "", "Run a specific shortcut by name/identifier (overrides auto-detection)")
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

	shortcutName, shortcutKind, err := resolveActionShortcut(ctx, def, cmdOpts, cfg)
	if err != nil {
		return err
	}
	input, err := buildActionInput(def, cmdOpts)
	if err != nil {
		return exitError(ExitCodeUsage, err)
	}
	if cmdOpts.dryRun {
		return printDryRun(opts, shortcutName, input)
	}

	if opts != nil && opts.noOutput {
		_, err := runShortcutWithRetry(ctx, shortcutName, input, opts, shortcutRunOptions(shortcutKind))
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

	out, err := runShortcutWithRetry(ctxRun, shortcutName, input, opts, shortcutRunOptions(shortcutKind))
	if err != nil {
		_ = appendTrace(cmdOpts.trace, traceEntry{Wrapper: shortcutName, Input: input, Error: err.Error()})
		return exitError(ExitCodeActionFailed, err)
	}
	_ = appendTrace(cmdOpts.trace, traceEntry{Wrapper: shortcutName, Input: input, Output: out})
	if opts.isJSON() || opts.isPlain() {
		var payload any
		if err := json.Unmarshal(out, &payload); err != nil {
			if shortcutKind == shortcutWrapper {
				return exitError(ExitCodeActionFailed, fmt.Errorf("wrapper output is not valid JSON: %w", err))
			}
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
	_, err = fmt.Fprint(os.Stdout, string(out))
	return err
}

func shortcutRunOptions(kind shortcutKind) shortcuts.RunOptions {
	if kind == shortcutWrapper {
		return shortcuts.RunOptions{OutputType: "public.json"}
	}
	return shortcuts.RunOptions{}
}

func resolveActionShortcut(ctx context.Context, def discovery.ActionDef, cmdOpts *actionCmdOptions, cfg config.Config) (string, shortcutKind, error) {
	if cmdOpts.shortcut != "" {
		return cmdOpts.shortcut, shortcutExplicit, nil
	}

	taskForShortcut := cmdOpts.task
	if taskForShortcut == "" {
		if task := taskFromInput(cmdOpts.input); task != "" {
			taskForShortcut = task
		}
	}
	name, candidates, err := resolveDirectShortcut(ctx, def, taskForShortcut)
	if err != nil {
		return "", shortcutDirect, err
	}
	if name != "" {
		return name, shortcutDirect, nil
	}

	wrapperName := cfg.Wrappers[def.ID]
	if wrapperName == "" {
		wrapperName = config.WrapperName(cfg.WrapperPrefix, def.ID)
	}
	wrapperExists, err := shortcutExists(ctx, wrapperName)
	if err != nil {
		return "", shortcutWrapper, err
	}
	if !wrapperExists && cfg.WrapperPrefix == config.DefaultWrapperPrefix {
		legacy := config.WrapperName(config.LegacyWrapperPrefix, def.ID)
		legacyExists, legacyErr := shortcutExists(ctx, legacy)
		if legacyErr == nil && legacyExists {
			wrapperName = legacy
			wrapperExists = true
		}
	}
	if wrapperExists {
		return wrapperName, shortcutWrapper, nil
	}

	if len(candidates) > 0 {
		return "", shortcutWrapper, exitError(ExitCodeWrappersMissing, fmt.Errorf("no matching Streaks shortcut found; expected one of: %s", strings.Join(candidates, ", ")))
	}
	return "", shortcutWrapper, exitError(ExitCodeWrappersMissing, fmt.Errorf("wrapper shortcut not found: %s", wrapperName))
}

func resolveDirectShortcut(ctx context.Context, def discovery.ActionDef, task string) (string, []string, error) {
	disc, err := discover(ctx)
	if err != nil {
		return "", nil, nil
	}
	available, err := listShortcuts(ctx)
	if err != nil {
		return "", nil, err
	}
	candidates := discovery.ActionShortcutCandidates(def, disc.App, disc.AppIntentKeys, disc.AppShortcutPhrases, task)
	if len(candidates) == 0 {
		return "", nil, nil
	}
	name := matchShortcutName(available, candidates)
	return name, candidates, nil
}

func matchShortcutName(shortcuts []shortcuts.Shortcut, candidates []string) string {
	if len(shortcuts) == 0 || len(candidates) == 0 {
		return ""
	}
	exact := make(map[string]string, len(shortcuts))
	normalized := make(map[string]string, len(shortcuts))
	for _, sc := range shortcuts {
		exact[sc.Name] = sc.Name
		normalized[strings.ToLower(strings.TrimSpace(sc.Name))] = sc.Name
	}
	for _, cand := range candidates {
		if name, ok := exact[cand]; ok {
			return name
		}
		if name, ok := normalized[strings.ToLower(strings.TrimSpace(cand))]; ok {
			return name
		}
	}
	return ""
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

func taskFromInput(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return ""
	}
	if task, ok := payload["task"].(string); ok {
		return strings.TrimSpace(task)
	}
	return ""
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
