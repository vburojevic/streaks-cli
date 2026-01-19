package cli

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
	"streaks-cli/internal/shortcuts"
)

type doctorReport struct {
	AppInstalled     bool     `json:"app_installed"`
	AppPath          string   `json:"app_path,omitempty"`
	BundleID         string   `json:"bundle_id,omitempty"`
	Version          string   `json:"version,omitempty"`
	ShortcutsCLI     bool     `json:"shortcuts_cli"`
	ShortcutsCLIPath string   `json:"shortcuts_cli_path,omitempty"`
	ConfigPath       string   `json:"config_path,omitempty"`
	ConfigPresent    bool     `json:"config_present"`
	WrapperShortcuts []string `json:"wrapper_shortcuts"`
	MissingWrappers  []string `json:"missing_wrappers"`
	URLSchemes       []string `json:"url_schemes,omitempty"`
	Warnings         []string `json:"warnings,omitempty"`
}

func newDoctorCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Verify Streaks installation and automation setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := runDoctor(context.Background())
			if err != nil {
				return err
			}
			if opts.json {
				return output.PrintJSON(os.Stdout, report, opts.pretty)
			}
			printDoctor(report)
			if !report.AppInstalled || !report.ShortcutsCLI {
				return fmt.Errorf("doctor checks failed")
			}
			return nil
		},
	}
	return cmd
}

func runDoctor(ctx context.Context) (doctorReport, error) {
	report := doctorReport{}
	defActions := discovery.DefaultActionDefinitions()
	cfg, present, _ := config.Load(defActions)
	configPath, _ := config.ConfigPath()
	report.ConfigPath = configPath
	report.ConfigPresent = present

	if _, err := os.Stat("/usr/bin/shortcuts"); err == nil {
		report.ShortcutsCLI = true
		report.ShortcutsCLIPath = "/usr/bin/shortcuts"
	}

	disc, err := discovery.Discover(ctx)
	if err == nil {
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
		report.Warnings = append(report.Warnings, err.Error())
	}

	if report.ShortcutsCLI {
		shortcutsList, err := shortcuts.List(ctx)
		if err == nil {
			installed := make(map[string]struct{}, len(shortcutsList))
			for _, sc := range shortcutsList {
				installed[sc.Name] = struct{}{}
			}
			wrappers := make([]string, 0, len(cfg.Wrappers))
			missing := make([]string, 0)
			for _, name := range cfg.Wrappers {
				wrappers = append(wrappers, name)
				if _, ok := installed[name]; !ok {
					missing = append(missing, name)
				}
			}
			sort.Strings(wrappers)
			sort.Strings(missing)
			report.WrapperShortcuts = wrappers
			report.MissingWrappers = missing
		} else {
			report.Warnings = append(report.Warnings, err.Error())
		}
	}

	return report, nil
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
	if report.ConfigPresent {
		fmt.Printf("Config: OK (%s)\n", report.ConfigPath)
	} else {
		fmt.Printf("Config: MISSING (%s)\n", report.ConfigPath)
	}
	if len(report.MissingWrappers) == 0 {
		fmt.Println("Wrapper shortcuts: OK")
	} else {
		fmt.Printf("Wrapper shortcuts: missing %d\n", len(report.MissingWrappers))
		for _, name := range report.MissingWrappers {
			fmt.Printf("  - %s\n", name)
		}
	}
	if len(report.Warnings) > 0 {
		fmt.Println("Warnings:")
		for _, warning := range report.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}
}
