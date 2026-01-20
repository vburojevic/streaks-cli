package cli

import "testing"

func TestErrorCodeLabel(t *testing.T) {
	cases := map[int]string{
		ExitCodeUsage:            "usage",
		ExitCodeAppMissing:       "app_missing",
		ExitCodeShortcutsMissing: "shortcuts_missing",
		ExitCodeShortcutMissing:  "shortcut_missing",
		ExitCodeActionFailed:     "action_failed",
		0:                        "",
		999:                      "",
	}
	for code, want := range cases {
		if got := errorCodeLabel(code); got != want {
			t.Fatalf("errorCodeLabel(%d) = %q, want %q", code, got, want)
		}
	}
}
