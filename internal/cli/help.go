package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"streaks-cli/internal/output"
)

func newHelpCmd(root *cobra.Command, opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Show help for a command",
		Args:  cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			target := root
			if len(args) > 0 {
				found, _, err := root.Find(args)
				if err != nil {
					return exitError(ExitCodeUsage, err)
				}
				if found == nil {
					return exitError(ExitCodeUsage, fmt.Errorf("unknown command: %s", strings.Join(args, " ")))
				}
				target = found
			}
			if opts != nil && opts.isAgent() {
				var buf bytes.Buffer
				origOut := target.OutOrStdout()
				origErr := target.ErrOrStderr()
				target.SetOut(&buf)
				target.SetErr(&buf)
				err := target.Help()
				target.SetOut(origOut)
				target.SetErr(origErr)
				if err != nil {
					return err
				}
				name := "st"
				if len(args) > 0 {
					name = "st " + strings.Join(args, " ")
				}
				payload := map[string]any{
					"command": name,
					"help":    strings.TrimSpace(buf.String()),
				}
				return output.PrintJSON(os.Stdout, payload, false)
			}
			return target.Help()
		},
	}
	return cmd
}
