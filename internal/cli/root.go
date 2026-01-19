package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "streaks-cli",
		Short:   "CLI for Streaks (Crunchy Bagel)",
		Version: version,
	}
	return cmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
