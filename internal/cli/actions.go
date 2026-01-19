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
var listShortcuts = shortcuts.List
var discover = discovery.Discover

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
		cmd.Flags().StringVar(&cmdOpts.input, "input", "", "Raw input to pass to the shortcut (overrides --task/--status)")
		cmd.Flags().BoolVar(&cmdOpts.stdin, "stdin", false, "Force reading input from stdin")
		cmd.Flags().BoolVar(&cmdOpts.dryRun, "dry-run", false, "Print shortcut and payload without running")
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
	input, err := buildActionInput(def, cmdOpts)
	if err != nil {
		return exitError(ExitCodeUsage, err)
	}
	if cmdOpts.shortcut != "" {
		return runNamedShortcut(ctx, cmdOpts.shortcut, input, cmdOpts, opts)
	}

	taskForShortcut := cmdOpts.task
	if taskForShortcut == "" {
		if task := taskFromInput(cmdOpts.input); task != "" {
			taskForShortcut = task
		}
	}
	candidates, err := actionCandidates(ctx, def, taskForShortcut)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		return exitError(ExitCodeShortcutMissing, fmt.Errorf("no shortcut candidates found for action %s", def.ID))
	}

	if available, err := listShortcuts(ctx); err == nil {
		if match := matchShortcutName(available, candidates); match != "" {
			if cmdOpts.dryRun {
				return printDryRun(opts, match, input)
			}
			return runNamedShortcut(ctx, match, input, cmdOpts, opts)
		}
	}

	if cmdOpts.dryRun {
		return printDryRun(opts, candidates[0], input)
	}
	return runCandidateShortcuts(ctx, candidates, def.ID, input, cmdOpts, opts)
}

func actionCandidates(ctx context.Context, def discovery.ActionDef, task string) ([]string, error) {
	disc, err := discover(ctx)
	if err != nil {
		return nil, exitError(ExitCodeAppMissing, err)
	}
	candidates := discovery.ActionShortcutCandidates(def, disc.App, disc.AppIntentKeys, disc.AppShortcutPhrases, task)
	return candidates, nil
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

func printDryRun(opts *rootOptions, shortcut string, input []byte) error {
	payload := map[string]any{
		"dry_run":  true,
		"shortcut": shortcut,
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
		fmt.Fprintf(os.Stdout, "Dry run: %s %s\n", shortcut, string(input))
		return nil
	}
	fmt.Fprintf(os.Stdout, "Dry run: %s\n", shortcut)
	return nil
}
