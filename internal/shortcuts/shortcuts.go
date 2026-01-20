package shortcuts

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Shortcut struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type RunOptions struct {
	OutputType string
}

var listLine = regexp.MustCompile(`^(.*) \(([0-9A-Fa-f-]+)\)$`)

func List(ctx context.Context) ([]Shortcut, error) {
	cmd := exec.CommandContext(ctx, "/usr/bin/shortcuts", "list", "--show-identifiers")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("shortcuts list failed: %w", err)
	}
	return parseList(out), nil
}

func parseList(output []byte) []Shortcut {
	var shortcuts []Shortcut
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		match := listLine.FindStringSubmatch(line)
		if len(match) == 3 {
			shortcuts = append(shortcuts, Shortcut{Name: match[1], ID: match[2]})
		} else {
			shortcuts = append(shortcuts, Shortcut{Name: line})
		}
	}
	return shortcuts
}

func Exists(ctx context.Context, name string) (bool, error) {
	shortcuts, err := List(ctx)
	if err != nil {
		return false, err
	}
	for _, sc := range shortcuts {
		if sc.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func Run(ctx context.Context, name string, input []byte) ([]byte, error) {
	return RunWithOptions(ctx, name, input, RunOptions{OutputType: "public.json"})
}

func RunWithOptions(ctx context.Context, name string, input []byte, opts RunOptions) ([]byte, error) {
	inputPath := ""
	outputDir := ""
	var err error
	if input == nil {
		input = []byte("{}")
	}
	if inputPath, err = writeTempFile("streaks-cli-input-*.json", input); err != nil {
		return nil, err
	}
	defer os.Remove(inputPath)

	if outputDir, err = os.MkdirTemp("", "streaks-cli-output-*"); err != nil {
		return nil, err
	}
	defer os.RemoveAll(outputDir)

	args := []string{"run", name, "--input-path", inputPath, "--output-path", outputDir}
	if strings.TrimSpace(opts.OutputType) != "" {
		args = append(args, "--output-type", opts.OutputType)
	}
	cmd := exec.CommandContext(ctx, "/usr/bin/shortcuts", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("shortcuts run failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return readOutputDir(outputDir)
}

func writeTempFile(pattern string, data []byte) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func readOutputDir(dir string) ([]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		files = append(files, name)
	}
	if len(files) == 0 {
		return []byte{}, nil
	}
	sort.Strings(files)
	if len(files) == 1 {
		return os.ReadFile(filepath.Join(dir, files[0]))
	}

	items := make([]any, 0, len(files))
	for _, name := range files {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		var payload any
		if err := json.Unmarshal(data, &payload); err == nil {
			items = append(items, payload)
			continue
		}
		items = append(items, strings.TrimSpace(string(data)))
	}
	return json.Marshal(items)
}
