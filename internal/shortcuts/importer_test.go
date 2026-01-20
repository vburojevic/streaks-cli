package shortcuts

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFindShortcutFiles(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "Test.shortcut")
	if err := os.WriteFile(file, []byte("placeholder"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	files, err := FindShortcutFiles(dir)
	if err != nil {
		t.Fatalf("FindShortcutFiles: %v", err)
	}
	if len(files) != 1 || files[0] != file {
		t.Fatalf("unexpected files: %v", files)
	}
}

func TestImportShortcutFilesMissingDir(t *testing.T) {
	_, err := ImportShortcutFiles(context.Background(), filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatalf("expected error")
	}
}
