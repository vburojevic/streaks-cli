package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
)

func newDiscoverCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Print discovered automation capabilities as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			disc, err := discovery.Discover(context.Background())
			if err != nil {
				return err
			}
			return output.PrintJSON(os.Stdout, disc, opts.pretty)
		},
	}
	return cmd
}
