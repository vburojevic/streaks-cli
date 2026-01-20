package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"streaks-cli/internal/config"
	"streaks-cli/internal/output"
)

type linkReport struct {
	Path     string             `json:"path"`
	Action   string             `json:"action"`
	Shortcut config.ShortcutRef `json:"shortcut"`
	Note     string             `json:"note,omitempty"`
}

type linksReport struct {
	Path     string                        `json:"path"`
	Mappings map[string]config.ShortcutRef `json:"mappings,omitempty"`
}

func newLinkCmd(opts *rootOptions) *cobra.Command {
	var shortcut string
	var shortcutName string
	var shortcutID string
	cmd := &cobra.Command{
		Use:   "link <action-id>",
		Short: "Map an action to a specific Shortcuts name or identifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			def, err := findActionDef(args[0])
			if err != nil {
				return exitError(ExitCodeUsage, err)
			}
			ref := config.ShortcutRef{}
			if shortcutID != "" {
				ref.ID = shortcutID
			}
			if shortcutName != "" {
				ref.Name = shortcutName
			}
			if shortcut != "" && ref.Name == "" && ref.ID == "" {
				ref.Name = shortcut
			}
			if ref.Name == "" && ref.ID == "" {
				return exitError(ExitCodeUsage, fmt.Errorf("provide --shortcut, --shortcut-name, or --shortcut-id"))
			}
			cfg, _, err := config.Load()
			if err != nil {
				return err
			}
			cfg.Mappings[def.ID] = ref
			path, err := config.Write(cfg)
			if err != nil {
				return err
			}
			report := linkReport{Path: path, Action: def.ID, Shortcut: ref}
			return printLinkReport(report, opts)
		},
	}
	cmd.Flags().StringVar(&shortcut, "shortcut", "", "Shortcut name or identifier to map to the action")
	cmd.Flags().StringVar(&shortcutName, "shortcut-name", "", "Shortcut name to map to the action")
	cmd.Flags().StringVar(&shortcutID, "shortcut-id", "", "Shortcut identifier to map to the action")
	return cmd
}

func newUnlinkCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlink <action-id>",
		Short: "Remove a shortcut mapping for an action",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			def, err := findActionDef(args[0])
			if err != nil {
				return exitError(ExitCodeUsage, err)
			}
			cfg, _, err := config.Load()
			if err != nil {
				return err
			}
			if _, ok := cfg.Mappings[def.ID]; !ok {
				report := linkReport{
					Path:   mustConfigPath(),
					Action: def.ID,
					Note:   "no mapping found",
				}
				return printLinkReport(report, opts)
			}
			delete(cfg.Mappings, def.ID)
			path, err := config.Write(cfg)
			if err != nil {
				return err
			}
			report := linkReport{
				Path:   path,
				Action: def.ID,
				Note:   "removed",
			}
			return printLinkReport(report, opts)
		},
	}
	return cmd
}

func newLinksCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "List configured action-to-shortcut mappings",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, _, err := config.Load()
			if err != nil {
				return err
			}
			path := mustConfigPath()
			report := linksReport{Path: path, Mappings: cfg.Mappings}
			return printLinksReport(report, opts)
		},
	}
	return cmd
}

func printLinkReport(report linkReport, opts *rootOptions) error {
	if opts.noOutput {
		return nil
	}
	if opts.isAgent() {
		return output.PrintJSON(os.Stdout, report, false)
	}
	if report.Note != "" {
		fmt.Printf("%s\t%s\t%s\n", report.Action, report.Note, report.Path)
		return nil
	}
	fmt.Printf("%s\t%s\n", report.Action, shortcutLabel(report.Shortcut))
	return nil
}

func printLinksReport(report linksReport, opts *rootOptions) error {
	if opts.noOutput {
		return nil
	}
	if opts.isAgent() {
		if len(report.Mappings) == 0 {
			return nil
		}
		ids := make([]string, 0, len(report.Mappings))
		for id := range report.Mappings {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			entry := linkReport{
				Path:     report.Path,
				Action:   id,
				Shortcut: report.Mappings[id],
			}
			if err := output.PrintJSON(os.Stdout, entry, false); err != nil {
				return err
			}
		}
		return nil
	}
	if len(report.Mappings) == 0 {
		fmt.Printf("No mappings configured (%s)\n", report.Path)
		return nil
	}
	ids := make([]string, 0, len(report.Mappings))
	for id := range report.Mappings {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		fmt.Printf("%s\t%s\n", id, shortcutLabel(report.Mappings[id]))
	}
	return nil
}

func shortcutLabel(ref config.ShortcutRef) string {
	if ref.Name != "" {
		return ref.Name
	}
	if ref.ID != "" {
		return ref.ID
	}
	return ""
}

func mustConfigPath() string {
	path, err := config.Path()
	if err != nil {
		return ""
	}
	return path
}
