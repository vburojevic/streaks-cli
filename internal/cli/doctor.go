package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
	"streaks-cli/internal/shortcuts"
)

type doctorReport struct {
	AppInstalled             bool     `json:"app_installed"`
	AppPath                  string   `json:"app_path,omitempty"`
	BundleID                 string   `json:"bundle_id,omitempty"`
	Version                  string   `json:"version,omitempty"`
	ShortcutsCLI             bool     `json:"shortcuts_cli"`
	ShortcutsCLIPath         string   `json:"shortcuts_cli_path,omitempty"`
	ShortcutCount            int      `json:"shortcut_count,omitempty"`
	ShortcutActionsAvailable []string `json:"shortcut_actions_available,omitempty"`
	ShortcutActionsMissing   []string `json:"shortcut_actions_missing,omitempty"`
	URLSchemes               []string `json:"url_schemes,omitempty"`
	Warnings                 []string `json:"warnings,omitempty"`
}

func newDoctorCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Verify Streaks installation and Shortcuts availability",
		RunE: func(_ *cobra.Command, _ []string) error {
			report, err := runDoctor(context.Background())
			if err != nil {
				return err
			}
			if opts.isAgent() {
				if err := output.PrintJSON(os.Stdout, report, false); err != nil {
					return err
				}
				return doctorExitError(report)
			}
			if opts.noOutput {
				return doctorExitError(report)
			}
			if !opts.quiet || doctorExitError(report) != nil {
				printDoctor(report)
			}
			return doctorExitError(report)
		},
	}
	return cmd
}

func runDoctor(ctx context.Context) (doctorReport, error) {
	report := doctorReport{}

	if _, err := os.Stat("/usr/bin/shortcuts"); err == nil {
		report.ShortcutsCLI = true
		report.ShortcutsCLIPath = "/usr/bin/shortcuts"
	}

	disc, discErr := discovery.Discover(ctx)
	if discErr == nil {
		report.AppInstalled = true
		report.AppPath = disc.App.Path
		report.BundleID = disc.App.BundleID
		report.Version = disc.App.Version
		if disc.ShortcutsCLIAvailable {
			report.ShortcutsCLI = true
			report.ShortcutsCLIPath = disc.ShortcutsCLIPath
		}
		report.URLSchemes = disc.URLSchemes
	} else {
		report.Warnings = append(report.Warnings, discErr.Error())
	}

	if report.ShortcutsCLI {
		list, err := shortcuts.List(ctx)
		if err != nil {
			report.Warnings = append(report.Warnings, err.Error())
			return report, nil
		}
		report.ShortcutCount = len(list)
		if discErr == nil {
			cfg, _, cfgErr := config.Load()
			if cfgErr != nil {
				report.Warnings = append(report.Warnings, cfgErr.Error())
			}
			available, missing := shortcutCoverage(discovery.DefaultActionDefinitions(), disc, list, cfg.Mappings)
			report.ShortcutActionsAvailable = available
			report.ShortcutActionsMissing = missing
		}
	}

	return report, nil
}

func doctorExitError(report doctorReport) error {
	if !report.AppInstalled {
		return exitError(ExitCodeAppMissing, errors.New("streaks app not found"))
	}
	if !report.ShortcutsCLI {
		return exitError(ExitCodeShortcutsMissing, errors.New("shortcuts CLI not available"))
	}
	if len(report.ShortcutActionsMissing) > 0 {
		return exitError(ExitCodeShortcutMissing, fmt.Errorf("missing %d Streaks shortcuts", len(report.ShortcutActionsMissing)))
	}
	return nil
}

func printDoctor(report doctorReport) {
	fmt.Println("Streaks CLI doctor")
	fmt.Println("------------------")
	if report.AppInstalled {
		fmt.Printf("App: OK (%s)\n", report.AppPath)
	} else {
		fmt.Println("App: MISSING")
	}
	if report.ShortcutsCLI {
		fmt.Printf("Shortcuts CLI: OK (%s)\n", report.ShortcutsCLIPath)
	} else {
		fmt.Println("Shortcuts CLI: MISSING")
	}
	if len(report.ShortcutActionsMissing) == 0 {
		fmt.Println("Streaks shortcuts: OK")
	} else {
		fmt.Printf("Streaks shortcuts: missing %d\n", len(report.ShortcutActionsMissing))
		for _, action := range report.ShortcutActionsMissing {
			fmt.Printf("  - %s\n", action)
		}
	}
	if len(report.Warnings) > 0 {
		fmt.Println("Warnings:")
		for _, warning := range report.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}
}

func printDoctorPlain(report doctorReport) {
	status := func(ok bool) string {
		if ok {
			return "ok"
		}
		return "missing"
	}
	fmt.Printf("app\t%s\t%s\n", status(report.AppInstalled), report.AppPath)
	fmt.Printf("shortcuts\t%s\t%s\n", status(report.ShortcutsCLI), report.ShortcutsCLIPath)
	if len(report.ShortcutActionsMissing) == 0 {
		fmt.Printf("actions\tok\t0\n")
	} else {
		fmt.Printf("actions\tmissing\t%d\n", len(report.ShortcutActionsMissing))
		for _, action := range report.ShortcutActionsMissing {
			fmt.Printf("action-missing\t%s\n", action)
		}
	}
	for _, warning := range report.Warnings {
		fmt.Printf("warning\t%s\n", warning)
	}
}
