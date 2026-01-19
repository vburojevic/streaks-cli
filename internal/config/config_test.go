package config

import (
	"os"
	"path/filepath"
	"testing"

	"streaks-cli/internal/discovery"
)

func TestConfigPathOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cfg.json")
	os.Setenv(EnvConfigPath, path)
	defer os.Unsetenv(EnvConfigPath)

	got, err := Path()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	if got != path {
		t.Fatalf("expected override path, got %s", got)
	}
}

func TestConfigReadWrite(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)
	defer os.Unsetenv("HOME")

	defs := discovery.DefaultActionDefinitions()
	cfg := DefaultConfig(defs)
	path, err := Write(cfg, true)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	loaded, present, err := Load(defs)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !present {
		t.Fatalf("expected config present")
	}
	if loaded.WrapperPrefix != DefaultWrapperPrefix {
		t.Fatalf("wrapper prefix mismatch")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config not written: %v", err)
	}
}
