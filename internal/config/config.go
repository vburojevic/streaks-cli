package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	DefaultConfigDirName  = "streaks-cli"
	DefaultConfigFileName = "config.json"
	EnvConfigPath         = "STREAKS_CLI_CONFIG"
)

type ShortcutRef struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type Config struct {
	Mappings map[string]ShortcutRef `json:"mappings,omitempty"`
	Prefer   string                 `json:"prefer,omitempty"` // "shim" or "auto"
}

func DefaultConfig() Config {
	return Config{
		Mappings: make(map[string]ShortcutRef),
	}
}

func Path() (string, error) {
	if override := os.Getenv(EnvConfigPath); override != "" {
		return override, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", DefaultConfigDirName, DefaultConfigFileName), nil
}

func Load() (Config, bool, error) {
	path, err := Path()
	if err != nil {
		return Config{}, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), false, nil
		}
		return Config{}, false, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, false, err
	}
	if cfg.Mappings == nil {
		cfg.Mappings = make(map[string]ShortcutRef)
	}
	return cfg, true, nil
}

func Write(cfg Config) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
