package shortcuts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const envShortcutDir = "STREAKS_CLI_SHORTCUT_DIR"

type ImportResult struct {
	Dir     string   `json:"dir"`
	Files   []string `json:"files,omitempty"`
	Opened  []string `json:"opened,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	Warning string   `json:"warning,omitempty"`
}

func DefaultImportDir() string {
	if override := os.Getenv(envShortcutDir); override != "" {
		return override
	}
	if existsDir("shortcuts") {
		return "shortcuts"
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates := []string{
			filepath.Join(base, "..", "share", "streaks-cli", "shortcuts"),
			"/usr/local/share/streaks-cli/shortcuts",
			"/opt/homebrew/share/streaks-cli/shortcuts",
		}
		for _, candidate := range candidates {
			if existsDir(candidate) {
				return candidate
			}
		}
	}
	return ""
}

func FindShortcutFiles(dir string) ([]string, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, fmt.Errorf("shortcut directory not specified")
	}
	if !existsDir(dir) {
		return nil, fmt.Errorf("shortcut directory not found: %s", dir)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.shortcut"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	return matches, nil
}

func ImportShortcutFiles(ctx context.Context, dir string) (ImportResult, error) {
	files, err := FindShortcutFiles(dir)
	if err != nil {
		return ImportResult{Dir: dir, Warning: err.Error()}, err
	}
	result := ImportResult{Dir: dir, Files: files}
	for _, file := range files {
		cmd := exec.CommandContext(ctx, "/usr/bin/open", "-a", "Shortcuts", file)
		if err := cmd.Run(); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", filepath.Base(file), err))
			continue
		}
		result.Opened = append(result.Opened, file)
	}
	return result, nil
}

func existsDir(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}
