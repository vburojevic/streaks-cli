package cli

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestEmitActionEnvelope(t *testing.T) {
	opts := &rootOptions{agent: true}
	result := runResult{Output: []byte(`{"ok":true}`), Attempts: 1, Duration: 15 * time.Millisecond}

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := emitActionOutput("task-list", "All Tasks", []byte(`{"input":true}`), result, opts)
	_ = w.Close()
	os.Stdout = origStdout

	if err != nil {
		t.Fatalf("emitActionOutput: %v", err)
	}
	var payload map[string]any
	dec := json.NewDecoder(r)
	if err := dec.Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	_ = r.Close()

	if ok, _ := payload["ok"].(bool); !ok {
		t.Fatalf("expected ok=true, got %v", payload["ok"])
	}
	action, _ := payload["action"].(map[string]any)
	if action["id"] != "task-list" {
		t.Fatalf("expected action id, got %v", action)
	}
	shortcut, _ := payload["shortcut"].(map[string]any)
	if shortcut["name"] != "All Tasks" {
		t.Fatalf("expected shortcut name, got %v", shortcut)
	}
	if _, ok := payload["result"]; !ok {
		t.Fatalf("expected result")
	}
}
