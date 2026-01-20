package cli

import (
	"os"
	"testing"
)

func TestNewRootCmdIncludesCommands(t *testing.T) {
	os.Setenv(envDisableDiscovery, "1")
	defer os.Unsetenv(envDisableDiscovery)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--agent", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("help failed: %v", err)
	}
	seen := map[string]bool{}
	for _, sub := range cmd.Commands() {
		seen[sub.Name()] = true
	}

	mustHave := []string{"discover", "doctor", "install", "link", "unlink", "links", "open", "actions", "task-complete", "task-list"}
	for _, name := range mustHave {
		if !seen[name] {
			t.Fatalf("missing command: %s", name)
		}
	}
}
