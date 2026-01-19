package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	streaksAppName  = "Streaks"
	streaksBundleID = "com.streaksapp.streak"
)

func Discover(ctx context.Context) (Discovery, error) {
	appPath, err := FindStreaksAppPath(ctx)
	if err != nil {
		return Discovery{}, err
	}

	info, urlSchemes, userActivities, err := ReadAppInfo(ctx, appPath)
	if err != nil {
		return Discovery{}, err
	}

	appIntentKeys, err := ReadAppIntentKeys(ctx, info.Resources)
	if err != nil {
		return Discovery{}, err
	}
	appShortcutKeys, _ := ReadAppShortcutKeys(ctx, info.Resources)

	actions, unmapped := DetectActions(extractKeys(appIntentKeys))

	shortcutsCLIPath := "/usr/bin/shortcuts"
	_, err = os.Stat(shortcutsCLIPath)
	shortcutsAvailable := err == nil

	d := Discovery{
		Timestamp:             time.Now().UTC().Format(time.RFC3339),
		App:                   info,
		URLSchemes:            urlSchemes,
		NSUserActivityTypes:   userActivities,
		ShortcutsCLIPath:      shortcutsCLIPath,
		ShortcutsCLIAvailable: shortcutsAvailable,
		AppIntentKeys:         appIntentKeys,
		AppShortcutKeys:       appShortcutKeys,
		Actions:               actions,
		UnmappedKeys:          unmapped,
	}

	if len(unmapped) > 0 {
		d.Notes = append(d.Notes, "Unmapped AppIntent keys detected; update the CLI mapping to cover new actions.")
	}

	return d, nil
}

func FindStreaksAppPath(ctx context.Context) (string, error) {
	defaultPath := "/Applications/Streaks.app"
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, nil
	}

	cmd := exec.CommandContext(ctx, "/usr/bin/mdfind", fmt.Sprintf("kMDItemCFBundleIdentifier == '%s'", streaksBundleID))
	out, err := cmd.Output()
	if err == nil {
		paths := filterAppPaths(strings.Split(strings.TrimSpace(string(out)), "\n"))
		if len(paths) > 0 {
			return paths[0], nil
		}
	}

	cmd = exec.CommandContext(ctx, "/usr/bin/mdfind", "kMDItemDisplayName == 'Streaks'")
	out, err = cmd.Output()
	if err != nil {
		return "", errors.New("Streaks app not found")
	}
	paths := filterAppPaths(strings.Split(strings.TrimSpace(string(out)), "\n"))
	if len(paths) == 0 {
		return "", errors.New("Streaks app not found")
	}
	return paths[0], nil
}

func filterAppPaths(paths []string) []string {
	filtered := make([]string, 0)
	for _, p := range paths {
		if p == "" {
			continue
		}
		if strings.HasSuffix(p, ".app") {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

type infoPlist struct {
	BundleID            string         `json:"CFBundleIdentifier"`
	ShortVersion        string         `json:"CFBundleShortVersionString"`
	BuildVersion        string         `json:"CFBundleVersion"`
	URLTypes            []urlTypeEntry `json:"CFBundleURLTypes"`
	NSUserActivityTypes []string       `json:"NSUserActivityTypes"`
}

type urlTypeEntry struct {
	Schemes []string `json:"CFBundleURLSchemes"`
}

func ReadAppInfo(ctx context.Context, appPath string) (AppInfo, []string, []string, error) {
	infoPath := filepath.Join(appPath, "Contents", "Info.plist")
	data, err := readPlistAsJSON(ctx, infoPath)
	if err != nil {
		return AppInfo{}, nil, nil, err
	}

	var plist infoPlist
	if err := json.Unmarshal(data, &plist); err != nil {
		return AppInfo{}, nil, nil, err
	}

	urlSchemes := make([]string, 0)
	for _, entry := range plist.URLTypes {
		for _, scheme := range entry.Schemes {
			if scheme != "" {
				urlSchemes = append(urlSchemes, scheme)
			}
		}
	}

	resources := filepath.Join(appPath, "Contents", "Resources")
	info := AppInfo{
		Name:      streaksAppName,
		Path:      appPath,
		BundleID:  plist.BundleID,
		Version:   plist.ShortVersion,
		Build:     plist.BuildVersion,
		Resources: resources,
	}

	return info, urlSchemes, plist.NSUserActivityTypes, nil
}

func ReadAppIntentKeys(ctx context.Context, resourcesPath string) ([]AppIntentKey, error) {
	localizable := filepath.Join(resourcesPath, "en.lproj", "Localizable.strings")
	data, err := readPlistAsJSON(ctx, localizable)
	if err != nil {
		return nil, err
	}

	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}

	keys := make([]AppIntentKey, 0)
	for key, value := range values {
		if strings.HasPrefix(key, "AppIntent.") {
			keys = append(keys, AppIntentKey{Key: key, Value: value})
		}
	}
	return keys, nil
}

func ReadAppShortcutKeys(ctx context.Context, resourcesPath string) ([]string, error) {
	var keys []string
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
		data, err := readPlistAsJSON(ctx, path)
		if err != nil {
			return nil
		}
		var values map[string]string
		if err := json.Unmarshal(data, &values); err != nil {
			return nil
		}
		for key := range values {
			if strings.Contains(key, "AppIntent.") {
				keys = append(keys, key)
			}
		}
		return nil
	})
	if err != nil {
		return keys, err
	}
	return uniqueStrings(keys), nil
}

func readPlistAsJSON(ctx context.Context, path string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "/usr/bin/plutil", "-convert", "json", "-o", "-", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("plutil failed for %s: %w: %s", path, err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func extractKeys(keys []AppIntentKey) []string {
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, k.Key)
	}
	return result
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, v := range input {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
