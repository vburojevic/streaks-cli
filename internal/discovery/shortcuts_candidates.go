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
	intentValues := orderedValues(intentKeys, PreferredLocales())
	intentNames := actionIntentNames(def)

	templates := make([]string, 0, len(def.Keys)+len(shortcutPhrases)+1)
	if def.Title != "" && shouldUseTitleTemplate(def) {
		templates = append(templates, def.Title)
	}
	for _, key := range def.Keys {
		if values, ok := intentValues[key]; ok {
			templates = append(templates, values...)
		}
	}
	phraseValues := orderedValues(shortcutPhrases, PreferredLocales())
	phraseKeys := make([]string, 0, len(phraseValues))
	for key := range phraseValues {
		phraseKeys = append(phraseKeys, key)
	}
	sort.Strings(phraseKeys)
	for _, key := range phraseKeys {
		for _, intent := range intentNames {
			if strings.Contains(key, "AppIntent."+intent+".") {
				templates = append(templates, phraseValues[key]...)
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

func shouldUseTitleTemplate(def ActionDef) bool {
	if !def.RequiresTask {
		return true
	}
	return strings.Contains(def.Title, "${task}") || strings.Contains(def.Title, "%@")
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

func orderedValues(keys []AppIntentKey, preferred []string) map[string][]string {
	byKey := make(map[string]map[string]string)
	for _, entry := range keys {
		if entry.Value == "" {
			continue
		}
		locale := entry.Locale
		if _, ok := byKey[entry.Key]; !ok {
			byKey[entry.Key] = make(map[string]string)
		}
		if _, ok := byKey[entry.Key][locale]; !ok {
			byKey[entry.Key][locale] = entry.Value
		}
	}

	out := make(map[string][]string, len(byKey))
	for key, valuesByLocale := range byKey {
		available := make(map[string]struct{}, len(valuesByLocale))
		for locale := range valuesByLocale {
			available[locale] = struct{}{}
		}
		order := localePriority(preferred, available)
		seen := make(map[string]struct{})
		ordered := make([]string, 0, len(valuesByLocale))
		for _, locale := range order {
			if value, ok := valuesByLocale[locale]; ok {
				if _, dup := seen[value]; dup {
					continue
				}
				seen[value] = struct{}{}
				ordered = append(ordered, value)
			}
		}
		out[key] = ordered
	}
	return out
}

func localePriority(preferred []string, available map[string]struct{}) []string {
	order := make([]string, 0, len(available))
	seen := make(map[string]struct{})
	add := func(locale string) {
		if locale == "" {
			return
		}
		if _, ok := seen[locale]; ok {
			return
		}
		if _, ok := available[locale]; !ok {
			return
		}
		seen[locale] = struct{}{}
		order = append(order, locale)
	}
	for _, locale := range preferred {
		add(locale)
		if parts := strings.Split(locale, "-"); len(parts) > 1 {
			add(parts[0])
		}
	}
	add("en")
	rest := make([]string, 0, len(available))
	for locale := range available {
		if _, ok := seen[locale]; !ok {
			rest = append(rest, locale)
		}
	}
	sort.Strings(rest)
	order = append(order, rest...)
	return order
}
