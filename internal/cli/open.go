package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func newOpenCmd(_ *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open Streaks via URL scheme",
		RunE: func(_ *cobra.Command, _ []string) error {
			return openURL("streaks://")
		},
	}
	return cmd
}

func openURL(url string) error {
	cmd := exec.Command("/usr/bin/open", url)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open failed: %w", err)
	}
	return nil
}
