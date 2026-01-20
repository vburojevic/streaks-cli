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
	keys := make([]AppIntentKey, 0)
	for locale, path := range paths {
		data, err := readPlistAsJSON(ctx, path)
		if err != nil {
			return nil, err
		}
		var values map[string]string
		if err := json.Unmarshal(data, &values); err != nil {
			return nil, err
		}
		for key, value := range values {
			if strings.Contains(key, "AppIntent.") {
				keys = append(keys, AppIntentKey{Key: key, Value: value, Locale: locale})
			}
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Key == keys[j].Key {
			return keys[i].Locale < keys[j].Locale
		}
		return keys[i].Key < keys[j].Key
	})
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
