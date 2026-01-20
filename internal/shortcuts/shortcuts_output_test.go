package shortcuts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadOutputDirEmpty(t *testing.T) {
	dir := t.TempDir()
	got, err := readOutputDir(dir)
	if err != nil {
		t.Fatalf("readOutputDir returned error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %q", string(got))
	}
}

func TestReadOutputDirSingleFile(t *testing.T) {
	dir := t.TempDir()
	payload := []byte(`{"ok":true}`)
	if err := os.WriteFile(filepath.Join(dir, "Dictionary.json"), payload, 0600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	got, err := readOutputDir(dir)
	if err != nil {
		t.Fatalf("readOutputDir returned error: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("expected %q, got %q", string(payload), string(got))
	}
}

func TestReadOutputDirMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.json"), []byte(`{"a":1}`), 0600); err != nil {
		t.Fatalf("write a.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("hello"), 0600); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "c.json"), []byte(`{"b":2}`), 0600); err != nil {
		t.Fatalf("write c.json: %v", err)
	}
	got, err := readOutputDir(dir)
	if err != nil {
		t.Fatalf("readOutputDir returned error: %v", err)
	}

	var items []any
	if err := json.Unmarshal(got, &items); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected first item: %#v", items[0])
	}
	if v, ok := first["a"].(float64); !ok || v != 1 {
		t.Fatalf("unexpected first item: %#v", items[0])
	}
	second, ok := items[1].(string)
	if !ok || second != "hello" {
		t.Fatalf("unexpected second item: %#v", items[1])
	}
	third, ok := items[2].(map[string]any)
	if !ok {
		t.Fatalf("unexpected third item: %#v", items[2])
	}
	if v, ok := third["b"].(float64); !ok || v != 2 {
		t.Fatalf("unexpected third item: %#v", items[2])
	}
}
