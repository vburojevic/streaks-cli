package config

import (
	"path/filepath"
	"testing"
)

func TestConfigPathOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "override.json")
	t.Setenv(EnvConfigPath, path)
	got, err := Path()
	if err != nil {
		t.Fatalf("Path: %v", err)
	}
	if got != path {
		t.Fatalf("expected %s, got %s", path, got)
	}
}

func TestConfigLoadMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	t.Setenv(EnvConfigPath, path)
	cfg, ok, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
	if cfg.Mappings == nil || len(cfg.Mappings) != 0 {
		t.Fatalf("expected empty mappings, got %v", cfg.Mappings)
	}
}

func TestConfigWriteAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv(EnvConfigPath, path)
	cfg := Config{
		Mappings: map[string]ShortcutRef{
			"task-list": {Name: "All Tasks"},
		},
	}
	wrote, err := Write(cfg)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if wrote != path {
		t.Fatalf("expected %s, got %s", path, wrote)
	}
	loaded, ok, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if loaded.Mappings["task-list"].Name != "All Tasks" {
		t.Fatalf("unexpected mapping: %v", loaded.Mappings)
	}
}
