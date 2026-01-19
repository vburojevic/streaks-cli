package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
)

var version = "dev"

type rootOptions struct {
	json   bool
	pretty bool
	config string
	agent  bool
}

const envDisableDiscovery = "STREAKS_CLI_DISABLE_DISCOVERY"
const envAgentMode = "STREAKS_CLI_AGENT"
const envJSONOutput = "STREAKS_CLI_JSON"

func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:           "st",
		Short:         "CLI for Streaks (Crunchy Bagel)",
		Long:          "CLI for Streaks (Crunchy Bagel).\n\nFor automation/agents, use --agent or --json for structured output.",
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.config != "" {
				if err := os.Setenv(config.EnvConfigPath, opts.config); err != nil {
					return err
				}
			}
			if opts.agent || isTruthy(os.Getenv(envAgentMode)) {
				opts.json = true
				opts.pretty = false
			}
			if opts.json {
				_ = os.Setenv(envJSONOutput, "1")
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.json, "json", false, "Output JSON when supported")
	cmd.PersistentFlags().BoolVar(&opts.pretty, "pretty", isTTY(os.Stdout), "Pretty-print JSON output")
	cmd.PersistentFlags().BoolVar(&opts.agent, "agent", false, "Agent-friendly mode (implies --json, disables pretty JSON)")
	cmd.PersistentFlags().StringVar(&opts.config, "config", "", "Path to config file (overrides STREAKS_CLI_CONFIG)")

	cmd.AddCommand(newDiscoverCmd(opts))
	cmd.AddCommand(newDoctorCmd(opts))
	cmd.AddCommand(newInstallCmd(opts))
	cmd.AddCommand(newOpenCmd(opts))
	cmd.AddCommand(newWrappersCmd(opts))

	addActionCommands(cmd, availableActionDefs(), opts)

	return cmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		if code, inner := exitCodeFromError(err); code != 0 {
			if os.Getenv(envJSONOutput) == "1" {
				_ = output.PrintJSON(os.Stderr, map[string]any{"error": inner.Error(), "code": code}, false)
			} else {
				fmt.Fprintln(os.Stderr, inner.Error())
			}
			os.Exit(code)
		}
		if os.Getenv(envJSONOutput) == "1" {
			_ = output.PrintJSON(os.Stderr, map[string]any{"error": err.Error(), "code": 1}, false)
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}

func filterDefs(defs []discovery.ActionDef, present map[string]discovery.Action) []discovery.ActionDef {
	out := make([]discovery.ActionDef, 0, len(defs))
	for _, def := range defs {
		if def.Transport == discovery.TransportURLScheme {
			continue
		}
		if _, ok := present[def.ID]; ok {
			out = append(out, def)
		}
	}
	return out
}

func isTruthy(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
