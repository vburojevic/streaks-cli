package cli

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
)

type wrapperEntry struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Wrapper      string              `json:"wrapper"`
	RequiresTask bool                `json:"requires_task"`
	Parameters   map[string][]string `json:"parameters,omitempty"`
}

func newWrappersCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wrappers",
		Short: "Manage wrapper shortcuts",
	}
	cmd.AddCommand(newWrappersListCmd(opts))
	cmd.AddCommand(newWrappersSampleCmd(opts))
	return cmd
}

func newWrappersListCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List expected wrapper shortcuts",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := listWrappers()
			if err != nil {
				return err
			}
			if opts.json {
				return output.PrintJSON(os.Stdout, entries, opts.pretty)
			}
			for _, entry := range entries {
				requires := ""
				if entry.RequiresTask {
					requires = " (task required)"
				}
				fmt.Printf("%s\t%s%s\n", entry.ID, entry.Wrapper, requires)
			}
			return nil
		},
	}
	return cmd
}

func newWrappersSampleCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sample <action-id>",
		Short: "Print a JSON input template for a wrapper",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def, err := findActionDef(args[0])
			if err != nil {
				return exitError(ExitCodeUsage, err)
			}
			payload := samplePayload(def)
			return output.PrintJSON(os.Stdout, payload, opts.pretty)
		},
	}
	return cmd
}

func listWrappers() ([]wrapperEntry, error) {
	defs := availableActionDefs()
	cfg, _, err := config.Load(discovery.DefaultActionDefinitions())
	if err != nil {
		return nil, err
	}
	entries := make([]wrapperEntry, 0, len(defs))
	for _, def := range defs {
		if def.Transport != discovery.TransportShortcuts {
			continue
		}
		wrapper := cfg.Wrappers[def.ID]
		if wrapper == "" {
			wrapper = config.WrapperName(cfg.WrapperPrefix, def.ID)
		}
		entries = append(entries, wrapperEntry{
			ID:           def.ID,
			Title:        def.Title,
			Wrapper:      wrapper,
			RequiresTask: def.RequiresTask,
			Parameters:   def.ParamOptions,
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	return entries, nil
}

func findActionDef(id string) (discovery.ActionDef, error) {
	for _, def := range availableActionDefs() {
		if def.ID == id {
			return def, nil
		}
	}
	return discovery.ActionDef{}, fmt.Errorf("unknown action: %s", id)
}

func samplePayload(def discovery.ActionDef) map[string]any {
	payload := map[string]any{}
	if def.RequiresTask {
		payload["task"] = "<task>"
	}
	if len(def.ParamOptions) > 0 {
		for key, values := range def.ParamOptions {
			if len(values) > 0 {
				payload[key] = values[0]
			} else {
				payload[key] = ""
			}
		}
	}
	return payload
}

func formatChecklist(entries []wrapperEntry) string {
	var b strings.Builder
	b.WriteString("Wrapper Checklist\n\n")
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("- [%s] %s (action: %s)\n", " ", entry.Wrapper, entry.ID))
		if entry.RequiresTask {
			b.WriteString("  input: task\n")
		}
		if len(entry.Parameters) > 0 {
			for key, values := range entry.Parameters {
				b.WriteString(fmt.Sprintf("  input: %s", key))
				if len(values) > 0 {
					b.WriteString(fmt.Sprintf(" (e.g. %s)", strings.Join(values, ", ")))
				}
				b.WriteString("\n")
			}
		}
	}
	return b.String()
}

func writeChecklist(path string, entries []wrapperEntry) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("checklist path is empty")
	}
	return os.WriteFile(path, []byte(formatChecklist(entries)), 0o644)
}
