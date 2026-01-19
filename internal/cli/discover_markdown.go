package cli

import (
	"fmt"
	"sort"
	"strings"

	"streaks-cli/internal/discovery"
)

func formatDiscoverMarkdown(d discovery.Discovery) string {
	var b strings.Builder
	b.WriteString("# Streaks Discovery\n\n")
	b.WriteString("## App\n")
	b.WriteString(fmt.Sprintf("- Name: %s\n", d.App.Name))
	b.WriteString(fmt.Sprintf("- Path: %s\n", d.App.Path))
	b.WriteString(fmt.Sprintf("- Bundle ID: %s\n", d.App.BundleID))
	b.WriteString(fmt.Sprintf("- Version: %s (%s)\n\n", d.App.Version, d.App.Build))

	if len(d.URLSchemes) > 0 {
		b.WriteString("## URL Schemes\n")
		for _, scheme := range d.URLSchemes {
			b.WriteString(fmt.Sprintf("- %s\n", scheme))
		}
		b.WriteString("\n")
	}

	if len(d.Actions) > 0 {
		b.WriteString("## Actions\n")
		actions := append([]discovery.Action{}, d.Actions...)
		sort.Slice(actions, func(i, j int) bool { return actions[i].ID < actions[j].ID })
		for _, action := range actions {
			requires := ""
			if action.RequiresTask {
				requires = " (task required)"
			}
			b.WriteString(fmt.Sprintf("- `%s`: %s [%s]%s\n", action.ID, action.Title, action.Transport, requires))
		}
		b.WriteString("\n")
	}

	if len(d.UnmappedKeys) > 0 {
		b.WriteString("## Unmapped App Intent Keys\n")
		for _, key := range d.UnmappedKeys {
			b.WriteString(fmt.Sprintf("- %s\n", key))
		}
		b.WriteString("\n")
	}

	if len(d.Notes) > 0 {
		b.WriteString("## Notes\n")
		for _, note := range d.Notes {
			b.WriteString(fmt.Sprintf("- %s\n", note))
		}
	}

	return b.String()
}
