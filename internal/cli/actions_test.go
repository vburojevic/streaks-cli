package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
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

func TestRunActionCommandUsesExplicitShortcut(t *testing.T) {
	origRun := runShortcut
	origDiscover := discover
	origList := listShortcuts
	defer func() {
		runShortcut = origRun
		discover = origDiscover
		listShortcuts = origList
	}()

	called := ""

	runShortcut = func(_ context.Context, name string, input []byte, _ shortcuts.RunOptions) ([]byte, error) {
		called = name
		return []byte(`{"ok":true}`), nil
	}
	discover = func(_ context.Context) (discovery.Discovery, error) {
		return discovery.Discovery{}, nil
	}
	listShortcuts = func(_ context.Context) ([]shortcuts.Shortcut, error) {
		return nil, nil
	}

	opts := &rootOptions{agent: true}
	def := discovery.ActionDef{ID: "task-list", Title: "List tasks", Transport: discovery.TransportShortcuts}

	t.Setenv("STREAKS_CLI_CONFIG", filepath.Join(t.TempDir(), "config.json"))
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runActionCommand(context.Background(), def, &actionCmdOptions{shortcut: "My Shortcut"}, opts)
	_ = w.Close()
	os.Stdout = origStdout

	out, _ := io.ReadAll(r)
	_ = r.Close()

	if err != nil {
		t.Fatalf("runActionCommand: %v", err)
	}
	if called != "My Shortcut" {
		t.Fatalf("expected shortcut name, got %s", called)
	}
	if !bytes.Contains(out, []byte("ok")) {
		t.Fatalf("unexpected output: %s", string(out))
	}
}

func TestRunActionCommandUsesDirectShortcut(t *testing.T) {
	origRun := runShortcut
	origDiscover := discover
	origList := listShortcuts
	defer func() {
		runShortcut = origRun
		discover = origDiscover
		listShortcuts = origList
	}()

	runShortcut = func(_ context.Context, name string, _ []byte, _ shortcuts.RunOptions) ([]byte, error) {
		if name != "All Tasks" {
			t.Fatalf("expected direct shortcut name, got %s", name)
		}
		return []byte(`{"ok":true}`), nil
	}
	discover = func(_ context.Context) (discovery.Discovery, error) {
		return discovery.Discovery{
			App: discovery.AppInfo{Name: "Streaks"},
			AppIntentKeys: []discovery.AppIntentKey{
				{Key: "AppIntent.TaskList.AllTasks", Value: "All Tasks", Locale: "en"},
			},
		}, nil
	}
	listShortcuts = func(_ context.Context) ([]shortcuts.Shortcut, error) {
		return []shortcuts.Shortcut{{Name: "All Tasks"}}, nil
	}

	opts := &rootOptions{agent: true}
	def := discovery.ActionDef{ID: "task-list", Title: "List tasks", Transport: discovery.TransportShortcuts, Keys: []string{"AppIntent.TaskList.AllTasks"}}

	t.Setenv("STREAKS_CLI_CONFIG", filepath.Join(t.TempDir(), "config.json"))
	if err := runActionCommand(context.Background(), def, &actionCmdOptions{}, opts); err != nil {
		t.Fatalf("runActionCommand: %v", err)
	}
}

func TestRunActionCommandUsesMapping(t *testing.T) {
	origRun := runShortcut
	origDiscover := discover
	origList := listShortcuts
	defer func() {
		runShortcut = origRun
		discover = origDiscover
		listShortcuts = origList
	}()

	called := ""
	runShortcut = func(_ context.Context, name string, _ []byte, _ shortcuts.RunOptions) ([]byte, error) {
		called = name
		return []byte(`{"ok":true}`), nil
	}
	discover = func(_ context.Context) (discovery.Discovery, error) {
		return discovery.Discovery{
			App: discovery.AppInfo{Name: "Streaks"},
			AppIntentKeys: []discovery.AppIntentKey{
				{Key: "AppIntent.TaskList.AllTasks", Value: "All Tasks", Locale: "en"},
			},
		}, nil
	}
	listShortcuts = func(_ context.Context) ([]shortcuts.Shortcut, error) {
		return []shortcuts.Shortcut{{Name: "All Tasks"}}, nil
	}

	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("STREAKS_CLI_CONFIG", path)
	_, err := config.Write(config.Config{
		Mappings: map[string]config.ShortcutRef{
			"task-list": {Name: "Mapped Shortcut"},
		},
	})
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	opts := &rootOptions{agent: true}
	def := discovery.ActionDef{ID: "task-list", Title: "List tasks", Transport: discovery.TransportShortcuts, Keys: []string{"AppIntent.TaskList.AllTasks"}}

	if err := runActionCommand(context.Background(), def, &actionCmdOptions{}, opts); err != nil {
		t.Fatalf("runActionCommand: %v", err)
	}
	if called != "Mapped Shortcut" {
		t.Fatalf("expected mapped shortcut, got %s", called)
	}
}
