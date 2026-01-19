package shortcuts

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Shortcut struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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
	inputPath := ""
	outputPath := ""
	var err error
	if input == nil {
		input = []byte("{}")
	}
	if inputPath, err = writeTempFile("streaks-cli-input-*.json", input); err != nil {
		return nil, err
	}
	defer os.Remove(inputPath)

	if outputPath, err = tempFilePath("streaks-cli-output-*.json"); err != nil {
		return nil, err
	}
	defer os.Remove(outputPath)

	args := []string{"run", name, "--input-path", inputPath, "--output-path", outputPath, "--output-type", "public.json"}
	cmd := exec.CommandContext(ctx, "/usr/bin/shortcuts", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("shortcuts run failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return os.ReadFile(outputPath)
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

func tempFilePath(pattern string) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	path := f.Name()
	if err := f.Close(); err != nil {
		return "", err
	}
	return path, nil
}
