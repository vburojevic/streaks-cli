package cli

import (
	"os"
	"testing"
)

func TestNewRootCmdIncludesCommands(t *testing.T) {
	os.Setenv(envDisableDiscovery, "1")
	defer os.Unsetenv(envDisableDiscovery)

	cmd := newRootCmd()
	seen := map[string]bool{}
	for _, sub := range cmd.Commands() {
		seen[sub.Name()] = true
	}

	mustHave := []string{"discover", "doctor", "install", "open", "wrappers", "task-complete", "task-list"}
	for _, name := range mustHave {
		if !seen[name] {
			t.Fatalf("missing command: %s", name)
		}
	}
}
