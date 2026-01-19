package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
)

var version = "dev"

type rootOptions struct {
	json   bool
	pretty bool
	config string
}

const envDisableDiscovery = "STREAKS_CLI_DISABLE_DISCOVERY"

func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:     "streaks-cli",
		Short:   "CLI for Streaks (Crunchy Bagel)",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.config != "" {
				if err := os.Setenv(config.EnvConfigPath, opts.config); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.json, "json", false, "Output JSON when supported")
	cmd.PersistentFlags().BoolVar(&opts.pretty, "pretty", isTTY(os.Stdout), "Pretty-print JSON output")
	cmd.PersistentFlags().StringVar(&opts.config, "config", "", "Path to config file (overrides STREAKS_CLI_CONFIG)")

	cmd.AddCommand(newDiscoverCmd(opts))
	cmd.AddCommand(newDoctorCmd(opts))
	cmd.AddCommand(newInstallCmd(opts))
	cmd.AddCommand(newOpenCmd(opts))

	defs := discovery.DefaultActionDefinitions()
	if os.Getenv(envDisableDiscovery) == "" {
		ctx := context.Background()
		disc, err := discovery.Discover(ctx)
		if err == nil && len(disc.Actions) > 0 {
			present := make(map[string]discovery.Action, len(disc.Actions))
			for _, action := range disc.Actions {
				present[action.ID] = action
			}
			defs = filterDefs(defs, present)
		}
	}
	addActionCommands(cmd, defs, opts)

	return cmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
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
