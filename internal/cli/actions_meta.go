package cli

import (
	"fmt"
	"os"
	"sort"

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
	Action  actionInfo     `json:"action"`
	Wrapper string         `json:"wrapper"`
	Sample  map[string]any `json:"sample_input"`
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
		RunE: func(cmd *cobra.Command, args []string) error {
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
			if opts.isJSON() {
				return output.PrintJSON(os.Stdout, infos, opts.pretty)
			}
			if opts.noOutput {
				return nil
			}
			if opts.isPlain() {
				for _, info := range infos {
					fmt.Printf("%s\t%t\n", info.ID, info.RequiresTask)
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
	cmd := &cobra.Command{
		Use:   "describe <action-id>",
		Short: "Describe an action and its input",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def, err := findActionDef(args[0])
			if err != nil {
				return exitError(ExitCodeUsage, err)
			}
			cfg, _, err := config.Load(discovery.DefaultActionDefinitions())
			if err != nil {
				return err
			}
			wrapper := cfg.Wrappers[def.ID]
			if wrapper == "" {
				wrapper = config.WrapperName(cfg.WrapperPrefix, def.ID)
			}
			detail := actionDetail{
				Action: actionInfo{
					ID:           def.ID,
					Title:        def.Title,
					Transport:    def.Transport,
					RequiresTask: def.RequiresTask,
					Parameters:   def.ParamOptions,
				},
				Wrapper: wrapper,
				Sample:  samplePayload(def),
			}
			if opts.isJSON() {
				return output.PrintJSON(os.Stdout, detail, opts.pretty)
			}
			if opts.noOutput {
				return nil
			}
			if opts.isPlain() {
				return output.PrintJSON(os.Stdout, detail, false)
			}
			fmt.Printf("ID: %s\nTitle: %s\nWrapper: %s\n", detail.Action.ID, detail.Action.Title, detail.Wrapper)
			if len(detail.Sample) > 0 {
				fmt.Printf("Sample input: %v\n", detail.Sample)
			}
			return nil
		},
	}
	return cmd
}
