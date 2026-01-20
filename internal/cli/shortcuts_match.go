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
	idExact := make(map[string]string, len(shortcuts))
	idNormalized := make(map[string]string, len(shortcuts))
	for _, sc := range shortcuts {
		exact[sc.Name] = sc.Name
		normalized[strings.ToLower(strings.TrimSpace(sc.Name))] = sc.Name
		if sc.ID != "" {
			idExact[sc.ID] = sc.ID
			idNormalized[strings.ToLower(strings.TrimSpace(sc.ID))] = sc.ID
		}
	}
	for _, cand := range candidates {
		if name, ok := exact[cand]; ok {
			return name
		}
		if name, ok := normalized[strings.ToLower(strings.TrimSpace(cand))]; ok {
			return name
		}
		if id, ok := idExact[cand]; ok {
			return id
		}
		if id, ok := idNormalized[strings.ToLower(strings.TrimSpace(cand))]; ok {
			return id
		}
	}
	return ""
}
