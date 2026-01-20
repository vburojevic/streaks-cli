package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
	"streaks-cli/internal/shortcuts"
)

type installResult struct {
	ShortcutActionsAvailable []string `json:"shortcut_actions_available,omitempty"`
	ShortcutActionsMissing   []string `json:"shortcut_actions_missing,omitempty"`
	Note                     string   `json:"note,omitempty"`
	ImportDir                string   `json:"import_dir,omitempty"`
	Imported                 []string `json:"imported,omitempty"`
	ImportErrors             []string `json:"import_errors,omitempty"`
	ImportWarning            string   `json:"import_warning,omitempty"`
}

func newInstallCmd(opts *rootOptions) *cobra.Command {
	installOpts := &installOptions{}
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Verify Streaks shortcuts are ready to use",
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := runInstall(context.Background(), installOpts)
			if err != nil {
				return err
			}
			if opts.isAgent() {
				if err := output.PrintJSON(os.Stdout, result, false); err != nil {
					return err
				}
				return installExitError(result)
			}
			if opts.noOutput {
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
			if result.ImportDir != "" {
				fmt.Printf("Import directory: %s\n", result.ImportDir)
			}
			if len(result.Imported) > 0 {
				fmt.Printf("Opened %d shortcut files for import.\n", len(result.Imported))
			}
			if len(result.ImportErrors) > 0 {
				fmt.Println("Import errors:")
				for _, err := range result.ImportErrors {
					fmt.Printf("  - %s\n", err)
				}
			}
			if result.ImportWarning != "" {
				fmt.Printf("Import warning: %s\n", result.ImportWarning)
			}
			if result.Note != "" {
				fmt.Printf("Note: %s\n", result.Note)
			}
			return installExitError(result)
		},
	}
	cmd.Flags().BoolVar(&installOpts.importShortcuts, "import", false, "Open bundled .shortcut wrapper files for import")
	cmd.Flags().StringVar(&installOpts.importDir, "from-dir", "", "Directory containing .shortcut files to import (default: bundled location)")
	return cmd
}

type installOptions struct {
	importShortcuts bool
	importDir       string
}

func runInstall(ctx context.Context, installOpts *installOptions) (installResult, error) {
	note := "The CLI uses existing Streaks shortcuts. Create shortcuts or map them with st link."
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
	cfg, _, cfgErr := config.Load()
	if cfgErr != nil {
		return installResult{}, cfgErr
	}
	available, missing := shortcutCoverage(discovery.DefaultActionDefinitions(), disc, list, cfg.Mappings)
	result := installResult{
		ShortcutActionsAvailable: available,
		ShortcutActionsMissing:   missing,
		Note:                     note,
	}
	if installOpts != nil && installOpts.importShortcuts {
		dir := strings.TrimSpace(installOpts.importDir)
		if dir == "" {
			dir = shortcuts.DefaultImportDir()
		}
		if dir == "" {
			result.ImportWarning = "no default shortcut directory found; pass --from-dir"
		} else {
			imp, impErr := shortcuts.ImportShortcutFiles(ctx, dir)
			result.ImportDir = imp.Dir
			result.Imported = imp.Opened
			result.ImportErrors = imp.Errors
			if imp.Warning != "" && result.ImportWarning == "" {
				result.ImportWarning = imp.Warning
			}
			if impErr != nil && result.ImportWarning == "" {
				result.ImportWarning = impErr.Error()
			}
		}
	}
	return result, nil
}

func installExitError(result installResult) error {
	if len(result.ShortcutActionsMissing) > 0 {
		return exitError(ExitCodeShortcutMissing, fmt.Errorf("missing %d Streaks shortcuts", len(result.ShortcutActionsMissing)))
	}
	return nil
}
