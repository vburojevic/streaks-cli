package cli

import (
	"context"
	"encoding/json"
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

type wrapperVerifyResult struct {
	ID          string `json:"id"`
	Wrapper     string `json:"wrapper"`
	Exists      bool   `json:"exists"`
	OutputValid bool   `json:"output_valid"`
	Skipped     bool   `json:"skipped"`
	Error       string `json:"error,omitempty"`
}

type wrappersDoctorReport struct {
	ConfigPath string                `json:"config_path"`
	Missing    []string              `json:"missing_wrappers"`
	Results    []wrapperVerifyResult `json:"verify_results,omitempty"`
}

func newWrappersCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wrappers",
		Short: "Manage wrapper shortcuts",
	}
	cmd.AddCommand(newWrappersListCmd(opts))
	cmd.AddCommand(newWrappersSampleCmd(opts))
	cmd.AddCommand(newWrappersVerifyCmd(opts))
	cmd.AddCommand(newWrappersDoctorCmd(opts))
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
			if opts.isJSON() {
				return output.PrintJSON(os.Stdout, entries, opts.pretty)
			}
			if opts.isPlain() {
				for _, entry := range entries {
					fmt.Printf("%s\t%s\t%t\n", entry.ID, entry.Wrapper, entry.RequiresTask)
				}
				return nil
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
			if opts.isPlain() {
				return output.PrintJSON(os.Stdout, payload, false)
			}
			return output.PrintJSON(os.Stdout, payload, opts.pretty)
		},
	}
	return cmd
}

func newWrappersVerifyCmd(opts *rootOptions) *cobra.Command {
	var task string
	var status string
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify wrapper shortcuts return JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := verifyWrappers(task, status, opts)
			if err != nil {
				return err
			}
			verifyErr := verifyExitError(results)
			if opts.isJSON() {
				if err := output.PrintJSON(os.Stdout, results, opts.pretty); err != nil {
					return err
				}
				return verifyErr
			}
			if opts.noOutput {
				return verifyErr
			}
			if opts.isPlain() {
				for _, res := range results {
					fmt.Printf("%s\t%s\t%t\t%t\t%t\t%s\n", res.ID, res.Wrapper, res.Exists, res.OutputValid, res.Skipped, res.Error)
				}
				return verifyErr
			}
			for _, res := range results {
				status := "ok"
				if !res.Exists {
					status = "missing"
				} else if res.Skipped {
					status = "skipped"
				} else if !res.OutputValid {
					status = "invalid"
				}
				if res.Error != "" {
					fmt.Printf("%s\t%s\t%s (%s)\n", res.ID, res.Wrapper, status, res.Error)
				} else {
					fmt.Printf("%s\t%s\t%s\n", res.ID, res.Wrapper, status)
				}
			}
			return verifyErr
		},
	}
	cmd.Flags().StringVar(&task, "task", "", "Task name for task-based actions")
	cmd.Flags().StringVar(&status, "status", "", "Status for pause action")
	return cmd
}

func newWrappersDoctorCmd(opts *rootOptions) *cobra.Command {
	var task string
	var status string
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Report wrapper readiness (existence + optional JSON validation)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _, err := config.Load(discovery.DefaultActionDefinitions())
			if err != nil {
				return err
			}
			configPath, _ := config.Path()
			missing, err := missingWrappers(context.Background(), cfg.Wrappers)
			if err != nil {
				return err
			}
			report := wrappersDoctorReport{ConfigPath: configPath, Missing: missing}

			if task != "" || status != "" {
				results, err := verifyWrappers(task, status, opts)
				if err != nil {
					return err
				}
				report.Results = results
			}

			if opts.isJSON() {
				if err := output.PrintJSON(os.Stdout, report, opts.pretty); err != nil {
					return err
				}
				if len(missing) > 0 {
					return exitError(ExitCodeWrappersMissing, fmt.Errorf("missing %d wrapper shortcuts", len(missing)))
				}
				return nil
			}
			if opts.noOutput {
				if len(missing) > 0 {
					return exitError(ExitCodeWrappersMissing, fmt.Errorf("missing %d wrapper shortcuts", len(missing)))
				}
				return nil
			}

			if opts.isPlain() {
				fmt.Printf("config\t%s\n", report.ConfigPath)
				if len(missing) == 0 {
					fmt.Println("wrappers\tok\t0")
				} else {
					fmt.Printf("wrappers\tmissing\t%d\n", len(missing))
					for _, name := range missing {
						fmt.Printf("wrapper-missing\t%s\n", name)
					}
				}
				for _, res := range report.Results {
					fmt.Printf("verify\t%s\t%t\t%t\t%t\t%s\n", res.ID, res.Exists, res.OutputValid, res.Skipped, res.Error)
				}
				if len(missing) > 0 {
					return exitError(ExitCodeWrappersMissing, fmt.Errorf("missing %d wrapper shortcuts", len(missing)))
				}
				return nil
			}

			if !opts.quiet {
				fmt.Printf("Config: %s\n", report.ConfigPath)
				if len(missing) == 0 {
					fmt.Println("Wrapper shortcuts: OK")
				} else {
					fmt.Printf("Missing %d wrapper shortcuts.\n", len(missing))
					for _, name := range missing {
						fmt.Printf("  - %s\n", name)
					}
				}
				if len(report.Results) > 0 {
					fmt.Println("Verification:")
					for _, res := range report.Results {
						if res.Error != "" {
							fmt.Printf("  - %s: %s\n", res.ID, res.Error)
						} else if res.Skipped {
							fmt.Printf("  - %s: skipped\n", res.ID)
						} else if res.OutputValid {
							fmt.Printf("  - %s: ok\n", res.ID)
						} else {
							fmt.Printf("  - %s: invalid output\n", res.ID)
						}
					}
				}
			}
			if len(missing) > 0 {
				return exitError(ExitCodeWrappersMissing, fmt.Errorf("missing %d wrapper shortcuts", len(missing)))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&task, "task", "", "Task name for task-based actions")
	cmd.Flags().StringVar(&status, "status", "", "Status for pause action")
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

func verifyWrappers(task, status string, opts *rootOptions) ([]wrapperVerifyResult, error) {
	entries, err := listWrappers()
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.RequiresTask && strings.TrimSpace(task) == "" {
			return nil, exitError(ExitCodeUsage, errors.New("--task is required to verify task-based wrappers"))
		}
	}

	results := make([]wrapperVerifyResult, 0, len(entries))
	ctx := context.Background()
	for _, entry := range entries {
		exists, err := shortcutExists(ctx, entry.Wrapper)
		if err != nil {
			results = append(results, wrapperVerifyResult{ID: entry.ID, Wrapper: entry.Wrapper, Exists: false, Error: err.Error()})
			continue
		}
		if !exists {
			results = append(results, wrapperVerifyResult{ID: entry.ID, Wrapper: entry.Wrapper, Exists: false})
			continue
		}

		input := map[string]any{}
		if entry.RequiresTask {
			input["task"] = strings.TrimSpace(task)
		}
		if len(entry.Parameters) > 0 {
			if status != "" {
				input["status"] = status
			} else {
				for key, values := range entry.Parameters {
					if len(values) > 0 {
						input[key] = values[0]
					} else {
						input[key] = ""
					}
				}
			}
		}
		payload, err := json.Marshal(input)
		if err != nil {
			results = append(results, wrapperVerifyResult{ID: entry.ID, Wrapper: entry.Wrapper, Exists: true, Error: err.Error()})
			continue
		}

		ctxRun := ctx
		var cancel context.CancelFunc
		if opts != nil && opts.timeout > 0 {
			ctxRun, cancel = context.WithTimeout(ctx, opts.timeout)
		}

		out, err := runShortcutWithRetry(ctxRun, entry.Wrapper, payload, opts)
		if cancel != nil {
			cancel()
		}
		if err != nil {
			results = append(results, wrapperVerifyResult{ID: entry.ID, Wrapper: entry.Wrapper, Exists: true, Error: err.Error()})
			continue
		}
		if !json.Valid(out) {
			results = append(results, wrapperVerifyResult{ID: entry.ID, Wrapper: entry.Wrapper, Exists: true, OutputValid: false, Error: "invalid JSON output"})
			continue
		}
		results = append(results, wrapperVerifyResult{ID: entry.ID, Wrapper: entry.Wrapper, Exists: true, OutputValid: true})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	return results, nil
}

func verifyExitError(results []wrapperVerifyResult) error {
	missing := 0
	failed := 0
	for _, res := range results {
		if !res.Exists {
			missing++
			continue
		}
		if res.Error != "" || (!res.Skipped && !res.OutputValid) {
			failed++
		}
	}
	if missing > 0 {
		return exitError(ExitCodeWrappersMissing, fmt.Errorf("missing %d wrapper shortcuts", missing))
	}
	if failed > 0 {
		return exitError(ExitCodeActionFailed, fmt.Errorf("%d wrapper checks failed", failed))
	}
	return nil
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
