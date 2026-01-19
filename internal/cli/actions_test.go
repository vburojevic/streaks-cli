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
	"streaks-cli/internal/shortcuts"
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
	origDiscover := discover
	origList := listShortcuts
	defer func() {
		runShortcut = origRun
		loadConfig = origLoad
		shortcutExists = origExists
		discover = origDiscover
		listShortcuts = origList
	}()

	called := struct {
		name  string
		input []byte
	}{}

	runShortcut = func(_ context.Context, name string, input []byte, _ shortcuts.RunOptions) ([]byte, error) {
		called.name = name
		called.input = input
		return []byte(`{"ok":true}`), nil
	}
	shortcutExists = func(_ context.Context, _ string) (bool, error) {
		return true, nil
	}
	discover = func(_ context.Context) (discovery.Discovery, error) {
		return discovery.Discovery{}, nil
	}
	listShortcuts = func(_ context.Context) ([]shortcuts.Shortcut, error) {
		return nil, nil
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

func TestRunActionCommandUsesDirectShortcutWhenWrapperMissing(t *testing.T) {
	origRun := runShortcut
	origLoad := loadConfig
	origExists := shortcutExists
	origDiscover := discover
	origList := listShortcuts
	defer func() {
		runShortcut = origRun
		loadConfig = origLoad
		shortcutExists = origExists
		discover = origDiscover
		listShortcuts = origList
	}()

	runShortcut = func(_ context.Context, name string, _ []byte, _ shortcuts.RunOptions) ([]byte, error) {
		if name != "All Tasks" {
			t.Fatalf("expected direct shortcut name, got %s", name)
		}
		return []byte(`{"ok":true}`), nil
	}
	shortcutExists = func(_ context.Context, _ string) (bool, error) {
		return false, nil
	}
	loadConfig = func(_ []discovery.ActionDef) (config.Config, bool, error) {
		return config.Config{WrapperPrefix: "st", Wrappers: map[string]string{}}, true, nil
	}
	discover = func(_ context.Context) (discovery.Discovery, error) {
		return discovery.Discovery{
			App: discovery.AppInfo{Name: "Streaks"},
			AppIntentKeys: []discovery.AppIntentKey{
				{Key: "AppIntent.TaskList.AllTasks", Value: "All Tasks"},
			},
		}, nil
	}
	listShortcuts = func(_ context.Context) ([]shortcuts.Shortcut, error) {
		return []shortcuts.Shortcut{{Name: "All Tasks"}}, nil
	}

	opts := &rootOptions{json: true, pretty: false}
	def := discovery.ActionDef{ID: "task-list", Title: "List tasks", Transport: discovery.TransportShortcuts, Keys: []string{"AppIntent.TaskList.AllTasks"}}

	if err := runActionCommand(context.Background(), def, &actionCmdOptions{}, opts); err != nil {
		t.Fatalf("runActionCommand: %v", err)
	}
}
