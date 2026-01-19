package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ReadAppShortcutPhrases(ctx context.Context, resourcesPath string) ([]AppIntentKey, error) {
	paths, err := findAppShortcutsFiles(resourcesPath)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no AppShortcuts.strings found")
	}
	selected := selectAppShortcutsPath(paths, preferredLocales())
	data, err := readPlistAsJSON(ctx, selected)
	if err != nil {
		return nil, err
	}
	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	keys := make([]AppIntentKey, 0, len(values))
	for key, value := range values {
		if strings.Contains(key, "AppIntent.") {
			keys = append(keys, AppIntentKey{Key: key, Value: value})
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Key < keys[j].Key })
	return keys, nil
}

func findAppShortcutsFiles(resourcesPath string) (map[string]string, error) {
	paths := make(map[string]string)
	err := filepath.WalkDir(resourcesPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) != "AppShortcuts.strings" {
			return nil
		}
		dir := filepath.Base(filepath.Dir(path))
		locale := strings.TrimSuffix(dir, ".lproj")
		if locale == "" {
			return nil
		}
		if _, ok := paths[locale]; !ok {
			paths[locale] = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func selectAppShortcutsPath(paths map[string]string, locales []string) string {
	for _, locale := range locales {
		if path, ok := paths[locale]; ok {
			return path
		}
		parts := strings.Split(locale, "-")
		if len(parts) > 1 {
			if path, ok := paths[parts[0]]; ok {
				return path
			}
		}
	}
	if path, ok := paths["en"]; ok {
		return path
	}
	localesAvailable := make([]string, 0, len(paths))
	for locale := range paths {
		localesAvailable = append(localesAvailable, locale)
	}
	sort.Strings(localesAvailable)
	return paths[localesAvailable[0]]
}

func preferredLocales() []string {
	envs := []string{"LC_ALL", "LC_MESSAGES", "LANG"}
	locales := make([]string, 0, len(envs))
	for _, env := range envs {
		if value := os.Getenv(env); value != "" {
			locale := strings.Split(value, ".")[0]
			locale = strings.ReplaceAll(locale, "_", "-")
			if locale != "" {
				locales = append(locales, locale)
			}
		}
	}
	return locales
}
