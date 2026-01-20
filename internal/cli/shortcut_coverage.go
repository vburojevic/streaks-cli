package cli

import (
	"sort"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
	"streaks-cli/internal/shortcuts"
)

func shortcutCoverage(defs []discovery.ActionDef, disc discovery.Discovery, list []shortcuts.Shortcut, mappings map[string]config.ShortcutRef) ([]string, []string) {
	available := make([]string, 0)
	missing := make([]string, 0)
	for _, def := range defs {
		if def.Transport != discovery.TransportShortcuts {
			continue
		}
		if ref, ok := mappings[def.ID]; ok {
			candidate := shortcutLabel(ref)
			if candidate != "" && matchShortcutName(list, []string{candidate}) != "" {
				available = append(available, def.ID)
				continue
			}
		}
		candidates := actionCandidatesFromDiscovery(def, disc, "")
		if len(candidates) == 0 {
			continue
		}
		if matchShortcutName(list, candidates) != "" {
			available = append(available, def.ID)
		} else {
			missing = append(missing, def.ID)
		}
	}
	sort.Strings(available)
	sort.Strings(missing)
	return available, missing
}
