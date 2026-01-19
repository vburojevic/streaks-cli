package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
	"streaks-cli/internal/shortcuts"
)

type installResult struct {
	ShortcutActionsAvailable []string `json:"shortcut_actions_available,omitempty"`
	ShortcutActionsMissing   []string `json:"shortcut_actions_missing,omitempty"`
	Note                     string   `json:"note,omitempty"`
}

func newInstallCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Verify Streaks shortcuts are ready to use",
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := runInstall(context.Background())
			if err != nil {
				return err
			}
			if opts.isJSON() {
				if err := output.PrintJSON(os.Stdout, result, opts.pretty); err != nil {
					return err
				}
				return installExitError(result)
			}
			if opts.noOutput {
				return installExitError(result)
			}
			if opts.isPlain() {
				for _, action := range result.ShortcutActionsAvailable {
					fmt.Printf("action\tavailable\t%s\n", action)
				}
				for _, action := range result.ShortcutActionsMissing {
					fmt.Printf("action\tmissing\t%s\n", action)
				}
				if result.Note != "" {
					fmt.Printf("note\t%s\n", result.Note)
				}
				return installExitError(result)
			}
			fmt.Println("Streaks shortcut readiness")
			if len(result.ShortcutActionsMissing) == 0 {
				fmt.Println("All non-task actions have matching shortcuts.")
			} else {
				fmt.Printf("Missing %d shortcuts for non-task actions:\n", len(result.ShortcutActionsMissing))
				for _, action := range result.ShortcutActionsMissing {
					fmt.Printf("  - %s\n", action)
				}
				fmt.Println("Create matching shortcuts in the Shortcuts app.")
			}
			if result.Note != "" {
				fmt.Printf("Note: %s\n", result.Note)
			}
			return installExitError(result)
		},
	}
	return cmd
}

func runInstall(ctx context.Context) (installResult, error) {
	note := "The CLI uses existing Streaks shortcuts. Create shortcuts for the actions you need or pass --shortcut."
	disc, err := discovery.Discover(ctx)
	if err != nil {
		return installResult{}, exitError(ExitCodeAppMissing, err)
	}
	if _, err := os.Stat(disc.ShortcutsCLIPath); err != nil {
		return installResult{}, exitError(ExitCodeShortcutsMissing, errors.New("shortcuts CLI not available"))
	}
	list, err := shortcuts.List(ctx)
	if err != nil {
		return installResult{}, exitError(ExitCodeShortcutsMissing, err)
	}
	available, missing := shortcutCoverage(discovery.DefaultActionDefinitions(), disc, list)
	return installResult{
		ShortcutActionsAvailable: available,
		ShortcutActionsMissing:   missing,
		Note:                     note,
	}, nil
}

func installExitError(result installResult) error {
	if len(result.ShortcutActionsMissing) > 0 {
		return exitError(ExitCodeShortcutMissing, fmt.Errorf("missing %d Streaks shortcuts", len(result.ShortcutActionsMissing)))
	}
	return nil
}
