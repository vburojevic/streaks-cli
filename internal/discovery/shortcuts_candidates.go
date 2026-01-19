package discovery

import (
	"sort"
	"strings"
)

func ActionShortcutCandidates(def ActionDef, app AppInfo, intentKeys []AppIntentKey, shortcutPhrases []AppIntentKey, task string) []string {
	if def.Transport != TransportShortcuts {
		return nil
	}
	appName := app.Name
	if appName == "" {
		appName = "Streaks"
	}
	intentValues := make(map[string]string, len(intentKeys))
	for _, key := range intentKeys {
		if key.Value != "" {
			intentValues[key.Key] = key.Value
		}
	}
	intentNames := actionIntentNames(def)

	templates := make([]string, 0, len(def.Keys)+len(shortcutPhrases)+1)
	templates = append(templates, def.Title)
	for _, key := range def.Keys {
		if value, ok := intentValues[key]; ok {
			templates = append(templates, value)
		}
	}
	for _, phrase := range shortcutPhrases {
		for _, intent := range intentNames {
			if strings.Contains(phrase.Key, "AppIntent."+intent+".") {
				if phrase.Value != "" {
					templates = append(templates, phrase.Value)
				}
				break
			}
		}
	}

	seen := make(map[string]struct{}, len(templates))
	candidates := make([]string, 0, len(templates))
	for _, tmpl := range templates {
		if tmpl == "" {
			continue
		}
		out, ok := expandShortcutTemplate(tmpl, task, appName)
		if !ok {
			continue
		}
		if _, exists := seen[out]; exists {
			continue
		}
		seen[out] = struct{}{}
		candidates = append(candidates, out)
	}
	return candidates
}

func actionIntentNames(def ActionDef) []string {
	names := make([]string, 0, len(def.Keys))
	seen := make(map[string]struct{}, len(def.Keys))
	for _, key := range def.Keys {
		if !strings.HasPrefix(key, "AppIntent.") {
			continue
		}
		trimmed := strings.TrimPrefix(key, "AppIntent.")
		parts := strings.SplitN(trimmed, ".", 2)
		if len(parts) == 0 || parts[0] == "" {
			continue
		}
		if _, ok := seen[parts[0]]; ok {
			continue
		}
		seen[parts[0]] = struct{}{}
		names = append(names, parts[0])
	}
	sort.Strings(names)
	return names
}

func expandShortcutTemplate(tmpl, task, appName string) (string, bool) {
	if strings.Contains(tmpl, "${task}") || strings.Contains(tmpl, "%@") {
		if strings.TrimSpace(task) == "" {
			return "", false
		}
	}
	out := strings.ReplaceAll(tmpl, "${task}", task)
	out = strings.ReplaceAll(out, "%@", task)
	out = strings.ReplaceAll(out, "${applicationName}", appName)
	out = strings.TrimSpace(out)
	if out == "" {
		return "", false
	}
	return out, true
}
