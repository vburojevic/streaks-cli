package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
)

var version = "dev"

type rootOptions struct {
	agent           bool
	quiet           bool
	verbose         bool
	noOutput        bool
	timeout         time.Duration
	retries         int
	retryWait       time.Duration
	configPath      string
	shortcutsOutput string
}

const envDisableDiscovery = "STREAKS_CLI_DISABLE_DISCOVERY"
const envAgentMode = "STREAKS_CLI_AGENT"
const envShortcutsOutput = "STREAKS_CLI_SHORTCUTS_OUTPUT"

func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:           "st",
		Short:         "CLI for Streaks (Crunchy Bagel)",
		Long:          "CLI for Streaks (Crunchy Bagel).\n\nFor automation/agents, use --agent for NDJSON output.",
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if opts.configPath != "" {
				_ = os.Setenv(config.EnvConfigPath, opts.configPath)
			}
			if opts.shortcutsOutput == "" {
				if env := os.Getenv(envShortcutsOutput); env != "" {
					opts.shortcutsOutput = env
				}
			}
			opts.agent = opts.agent || isTruthy(os.Getenv(envAgentMode))
			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.agent, "agent", false, "Agent-friendly mode (NDJSON output)")
	cmd.PersistentFlags().BoolVar(&opts.quiet, "quiet", false, "Suppress non-essential output")
	cmd.PersistentFlags().BoolVar(&opts.verbose, "verbose", false, "Verbose output")
	cmd.PersistentFlags().BoolVar(&opts.noOutput, "no-output", false, "Suppress all output (exit code only)")
	cmd.PersistentFlags().DurationVar(&opts.timeout, "timeout", 30*time.Second, "Timeout for Shortcuts runs")
	cmd.PersistentFlags().IntVar(&opts.retries, "retries", 0, "Retry failed Shortcuts runs")
	cmd.PersistentFlags().DurationVar(&opts.retryWait, "retry-delay", time.Second, "Initial delay between retries")
	cmd.PersistentFlags().StringVar(&opts.shortcutsOutput, "shortcuts-output", "public.plain-text", "Shortcuts output type (UTI), e.g. public.plain-text or public.json")
	cmd.PersistentFlags().StringVar(&opts.configPath, "config", "", "Path to config file (default: ~/.config/streaks-cli/config.json)")

	cmd.AddCommand(newDiscoverCmd(opts))
	cmd.AddCommand(newDoctorCmd(opts))
	cmd.AddCommand(newInstallCmd(opts))
	cmd.AddCommand(newLinkCmd(opts))
	cmd.AddCommand(newUnlinkCmd(opts))
	cmd.AddCommand(newLinksCmd(opts))
	cmd.AddCommand(newHelpCmd(cmd, opts))
	cmd.AddCommand(newOpenCmd(opts))
	cmd.AddCommand(newActionsCmd(opts))

	addActionCommands(cmd, availableActionDefs(), opts)

	return cmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		if code, inner := exitCodeFromError(err); code != 0 {
			printError(inner.Error(), code)
			os.Exit(code)
		}
		printError(err.Error(), 1)
		os.Exit(1)
	}
}

func printError(message string, code int) {
	if isTruthy(os.Getenv(envAgentMode)) {
		payload := map[string]any{"error": message, "code": code}
		if code == ExitCodeUsage {
			payload["hint"] = "Run `st help` or `st help <command>` to see usage."
		}
		_ = output.PrintJSON(os.Stderr, payload, false)
		return
	}
	fmt.Fprintln(os.Stderr, message)
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
