package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIntegrationDiscover(t *testing.T) {
	if os.Getenv("STREAKS_CLI_INTEGRATION") == "" {
		t.Skip("set STREAKS_CLI_INTEGRATION=1 to run integration tests")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(cwd, "../.."))
	cmd := exec.Command("go", "run", "./cmd/streaks-cli", "discover")
	cmd.Dir = root
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if _, ok := payload["app"]; !ok {
		t.Fatalf("missing app field")
	}
}
