package cli

import (
	"strings"

	"streaks-cli/internal/shortcuts"
)

func matchShortcutName(shortcuts []shortcuts.Shortcut, candidates []string) string {
	if len(shortcuts) == 0 || len(candidates) == 0 {
		return ""
	}
	exact := make(map[string]string, len(shortcuts))
	normalized := make(map[string]string, len(shortcuts))
	for _, sc := range shortcuts {
		exact[sc.Name] = sc.Name
		normalized[strings.ToLower(strings.TrimSpace(sc.Name))] = sc.Name
	}
	for _, cand := range candidates {
		if name, ok := exact[cand]; ok {
			return name
		}
		if name, ok := normalized[strings.ToLower(strings.TrimSpace(cand))]; ok {
			return name
		}
	}
	return ""
}
