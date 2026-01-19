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
			RunE: func(cmd *cobra.Command, args []string) error {
				return runActionCommand(context.Background(), def, cmdOpts, opts)
			},
		}
		cmd.Flags().StringVar(&cmdOpts.input, "input", "", "Raw input to pass to the wrapper (overrides --task/--status)")
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
	if !exists {
		return exitError(ExitCodeWrappersMissing, fmt.Errorf("wrapper shortcut not found: %s", wrapperName))
	}
	input, err := buildActionInput(def, cmdOpts)
	if err != nil {
		return exitError(ExitCodeUsage, err)
	}
	out, err := runShortcut(ctx, wrapperName, input)
	if err != nil {
		return exitError(ExitCodeActionFailed, err)
	}
	if opts.json {
		var payload any
		if err := json.Unmarshal(out, &payload); err != nil {
			return fmt.Errorf("wrapper output is not valid JSON: %w", err)
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

	if !isTTY(os.Stdin) {
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
