package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"streaks-cli/internal/config"
	"streaks-cli/internal/discovery"
)

func TestBuildActionInputFromFlags(t *testing.T) {
	def := discovery.ActionDef{ID: "task-complete", RequiresTask: true}
	opts := &actionCmdOptions{task: "Read", status: "All"}
	data, err := buildActionInput(def, opts)
	if err != nil {
		t.Fatalf("buildActionInput: %v", err)
	}
	var payload map[string]string
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["task"] != "Read" || payload["status"] != "All" {
		t.Fatalf("unexpected payload: %v", payload)
	}
}

func TestBuildActionInputFromStdin(t *testing.T) {
	def := discovery.ActionDef{ID: "task-list", RequiresTask: false}
	opts := &actionCmdOptions{}
	input := []byte(`{"ok":true}`)

	orig := os.Stdin
	r, w, _ := os.Pipe()
	_, _ = w.Write(input)
	_ = w.Close()
	os.Stdin = r
	defer func() {
		os.Stdin = orig
		_ = r.Close()
	}()

	data, err := buildActionInput(def, opts)
	if err != nil {
		t.Fatalf("buildActionInput: %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(data), input) {
		t.Fatalf("unexpected stdin payload: %s", string(data))
	}
}

func TestRunActionCommandUsesWrapper(t *testing.T) {
	origRun := runShortcut
	origLoad := loadConfig
	origExists := shortcutExists
	defer func() {
		runShortcut = origRun
		loadConfig = origLoad
		shortcutExists = origExists
	}()

	called := struct {
		name  string
		input []byte
	}{}

	runShortcut = func(_ context.Context, name string, input []byte) ([]byte, error) {
		called.name = name
		called.input = input
		return []byte(`{"ok":true}`), nil
	}
	shortcutExists = func(_ context.Context, name string) (bool, error) {
		return true, nil
	}
	loadConfig = func(_ []discovery.ActionDef) (config.Config, bool, error) {
		return config.Config{WrapperPrefix: "streaks-cli", Wrappers: map[string]string{"task-list": "wrapper"}}, true, nil
	}

	opts := &rootOptions{json: true, pretty: false}
	def := discovery.ActionDef{ID: "task-list", Title: "List tasks", Transport: discovery.TransportShortcuts}

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runActionCommand(context.Background(), def, &actionCmdOptions{}, opts)
	_ = w.Close()
	os.Stdout = origStdout

	out, _ := io.ReadAll(r)
	_ = r.Close()

	if err != nil {
		t.Fatalf("runActionCommand: %v", err)
	}
	if called.name != "wrapper" {
		t.Fatalf("expected wrapper name, got %s", called.name)
	}
	if !bytes.Contains(out, []byte("ok")) {
		t.Fatalf("unexpected output: %s", string(out))
	}
}
