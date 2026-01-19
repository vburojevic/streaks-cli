package cli

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/shortcuts"
)

func newInstallCmd(opts *rootOptions) *cobra.Command {
	var force bool
	var checklistPath string
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install or verify wrapper shortcuts and config",
		RunE: func(cmd *cobra.Command, args []string) error {
			defs := discovery.DefaultActionDefinitions()
			cfg := config.DefaultConfig(defs)
			path, err := config.Write(cfg, force)
			if err != nil {
				return err
			}
			fmt.Printf("Config written: %s\n", path)

			if checklistPath != "" {
				entries, err := listWrappers()
				if err != nil {
					return err
				}
				if err := writeChecklist(checklistPath, entries); err != nil {
					return err
				}
				fmt.Printf("Checklist written: %s\n", checklistPath)
			}

			missing, err := missingWrappers(context.Background(), cfg.Wrappers)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Shortcuts check failed: %v\n", err)
				fmt.Println("See docs/setup.md for manual wrapper setup.")
				return nil
			}
			if len(missing) > 0 {
				fmt.Printf("Missing %d wrapper shortcuts.\n", len(missing))
				for _, name := range missing {
					fmt.Printf("  - %s\n", name)
				}
				fmt.Println("See docs/setup.md for manual wrapper setup.")
				return nil
			}
			fmt.Println("Wrapper shortcuts: OK")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing config")
	cmd.Flags().StringVar(&checklistPath, "checklist", "", "Write a wrapper checklist to a file")
	return cmd
}

func missingWrappers(ctx context.Context, wrappers map[string]string) ([]string, error) {
	list, err := shortcuts.List(ctx)
	if err != nil {
		return nil, err
	}
	installed := make(map[string]struct{}, len(list))
	for _, sc := range list {
		installed[sc.Name] = struct{}{}
	}
	missing := make([]string, 0)
	for _, name := range wrappers {
		if _, ok := installed[name]; !ok {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	return missing, nil
}
