package cli

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
)

type actionInfo struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Transport    string              `json:"transport"`
	RequiresTask bool                `json:"requires_task"`
	Parameters   map[string][]string `json:"parameters,omitempty"`
}

type actionDetail struct {
	Action             actionInfo          `json:"action"`
	Sample             map[string]any      `json:"sample_input"`
	ShortcutCandidates []string            `json:"shortcut_candidates,omitempty"`
	MappedShortcut     *config.ShortcutRef `json:"mapped_shortcut,omitempty"`
}

func newActionsCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Inspect available actions",
	}
	cmd.AddCommand(newActionsListCmd(opts))
	cmd.AddCommand(newActionsDescribeCmd(opts))
	return cmd
}

func newActionsListCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available actions",
		RunE: func(_ *cobra.Command, _ []string) error {
			infos := make([]actionInfo, 0)
			for _, def := range availableActionDefs() {
				if def.Transport != discovery.TransportShortcuts {
					continue
				}
				infos = append(infos, actionInfo{
					ID:           def.ID,
					Title:        def.Title,
					Transport:    def.Transport,
					RequiresTask: def.RequiresTask,
					Parameters:   def.ParamOptions,
				})
			}
			sort.Slice(infos, func(i, j int) bool { return infos[i].ID < infos[j].ID })
			if opts.noOutput {
				return nil
			}
			if opts.isAgent() {
				for _, info := range infos {
					if err := output.PrintJSON(os.Stdout, info, false); err != nil {
						return err
					}
				}
				return nil
			}
			for _, info := range infos {
				requires := ""
				if info.RequiresTask {
					requires = " (task required)"
				}
				fmt.Printf("%s\t%s%s\n", info.ID, info.Title, requires)
			}
			return nil
		},
	}
	return cmd
}

func newActionsDescribeCmd(opts *rootOptions) *cobra.Command {
	var task string
	cmd := &cobra.Command{
		Use:   "describe <action-id>",
		Short: "Describe an action and its input",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			def, err := findActionDef(args[0])
			if err != nil {
				return exitError(ExitCodeUsage, err)
			}
			var shortcutCandidates []string
			if disc, err := discover(context.Background()); err == nil {
				shortcutCandidates = actionCandidatesFromDiscovery(def, disc, task)
			}
			var mapped *config.ShortcutRef
			if cfg, _, err := config.Load(); err == nil {
				if ref, ok := cfg.Mappings[def.ID]; ok {
					copy := ref
					mapped = &copy
				}
			}
			detail := actionDetail{
				Action: actionInfo{
					ID:           def.ID,
					Title:        def.Title,
					Transport:    def.Transport,
					RequiresTask: def.RequiresTask,
					Parameters:   def.ParamOptions,
				},
				Sample:             samplePayload(def),
				ShortcutCandidates: shortcutCandidates,
				MappedShortcut:     mapped,
			}
			if opts.noOutput {
				return nil
			}
			if opts.isAgent() {
				return output.PrintJSON(os.Stdout, detail, false)
			}
			fmt.Printf("ID: %s\nTitle: %s\n", detail.Action.ID, detail.Action.Title)
			if len(detail.Sample) > 0 {
				fmt.Printf("Sample input: %v\n", detail.Sample)
			}
			if len(detail.ShortcutCandidates) > 0 {
				fmt.Printf("Shortcut candidates: %s\n", strings.Join(detail.ShortcutCandidates, ", "))
			}
			if detail.MappedShortcut != nil {
				fmt.Printf("Mapped shortcut: %s\n", shortcutLabel(*detail.MappedShortcut))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&task, "task", "", "Task name to expand shortcut templates")
	return cmd
}
