package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"streaks-cli/internal/discovery"
)

const (
	DefaultConfigDirName  = "streaks-cli"
	DefaultConfigFileName = "config.json"
	EnvConfigPath         = "STREAKS_CLI_CONFIG"
	DefaultWrapperPrefix  = "st"
	LegacyWrapperPrefix   = "streaks-cli"
)

type Config struct {
	WrapperPrefix string            `json:"wrapper_prefix"`
	Wrappers      map[string]string `json:"wrappers"`
}

func DefaultConfig(actions []discovery.ActionDef) Config {
	wrappers := make(map[string]string)
	for _, action := range actions {
		if action.Transport != discovery.TransportShortcuts {
			continue
		}
		wrappers[action.ID] = WrapperName(DefaultWrapperPrefix, action.ID)
	}
	return Config{
		WrapperPrefix: DefaultWrapperPrefix,
		Wrappers:      wrappers,
	}
}

func WrapperName(prefix, actionID string) string {
	return fmt.Sprintf("%s %s", prefix, actionID)
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

func Load(actions []discovery.ActionDef) (Config, bool, error) {
	path, err := Path()
	if err != nil {
		return Config{}, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(actions), false, nil
		}
		return Config{}, false, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, false, err
	}
	if cfg.WrapperPrefix == "" {
		cfg.WrapperPrefix = DefaultWrapperPrefix
	}
	if cfg.Wrappers == nil {
		cfg.Wrappers = DefaultConfig(actions).Wrappers
	}
	return cfg, true, nil
}

func Write(cfg Config, force bool) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	if !force {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
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
