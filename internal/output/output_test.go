package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	payload := map[string]string{"ok": "true"}
	if err := PrintJSON(buf, payload, true); err != nil {
		t.Fatalf("PrintJSON: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"ok\"") {
		t.Fatalf("unexpected output: %s", out)
	}
}
