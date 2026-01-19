package shortcuts

import "testing"

func TestParseList(t *testing.T) {
	input := []byte("Example One (AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE)\nPlain Shortcut\n")
	shortcuts := parseList(input)
	if len(shortcuts) != 2 {
		t.Fatalf("expected 2 shortcuts, got %d", len(shortcuts))
	}
	if shortcuts[0].Name != "Example One" || shortcuts[0].ID == "" {
		t.Fatalf("unexpected shortcut[0]: %+v", shortcuts[0])
	}
	if shortcuts[1].Name != "Plain Shortcut" {
		t.Fatalf("unexpected shortcut[1]: %+v", shortcuts[1])
	}
}
