package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"streaks-cli/internal/discovery"
	"streaks-cli/internal/output"
)

func newDiscoverCmd(opts *rootOptions) *cobra.Command {
	var markdown bool
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Print discovered automation capabilities as JSON",
		RunE: func(_ *cobra.Command, _ []string) error {
			disc, err := discovery.Discover(context.Background())
			if err != nil {
				return exitError(ExitCodeAppMissing, err)
			}
			if markdown {
				if opts.isJSON() {
					return exitError(ExitCodeUsage, fmt.Errorf("--markdown is incompatible with JSON output"))
				}
				_, err := os.Stdout.WriteString(formatDiscoverMarkdown(disc))
				return err
			}
			if opts.isPlain() {
				return output.PrintJSON(os.Stdout, disc, false)
			}
			return output.PrintJSON(os.Stdout, disc, opts.pretty)
		},
	}
	cmd.Flags().BoolVar(&markdown, "markdown", false, "Output discovery report as Markdown")
	return cmd
}
