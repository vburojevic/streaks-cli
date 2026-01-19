package discovery

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestReadAppInfo(t *testing.T) {
	root := t.TempDir()
	appPath := filepath.Join(root, "Streaks.app")
	infoPath := filepath.Join(appPath, "Contents", "Info.plist")

	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>com.streaksapp.streak</string>
	<key>CFBundleShortVersionString</key>
	<string>1.2.3</string>
	<key>CFBundleVersion</key>
	<string>123</string>
	<key>CFBundleURLTypes</key>
	<array>
		<dict>
			<key>CFBundleURLSchemes</key>
			<array>
				<string>streaks</string>
			</array>
		</dict>
	</array>
	<key>NSUserActivityTypes</key>
	<array>
		<string>INAddTasksIntent</string>
	</array>
</dict>
</plist>`

	writeFile(t, infoPath, plist)
	info, schemes, activities, err := ReadAppInfo(context.Background(), appPath)
	if err != nil {
		t.Fatalf("ReadAppInfo: %v", err)
	}
	if info.BundleID != "com.streaksapp.streak" {
		t.Fatalf("bundle id mismatch: %s", info.BundleID)
	}
	if len(schemes) != 1 || schemes[0] != "streaks" {
		t.Fatalf("schemes mismatch: %v", schemes)
	}
	if len(activities) != 1 || activities[0] != "INAddTasksIntent" {
		t.Fatalf("activities mismatch: %v", activities)
	}
}

func TestReadAppIntentKeys(t *testing.T) {
	root := t.TempDir()
	res := filepath.Join(root, "Resources")
	path := filepath.Join(res, "en.lproj", "Localizable.strings")

	stringsPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIntent.TaskList.AllTasks</key>
	<string>All Tasks</string>
	<key>SomethingElse</key>
	<string>Nope</string>
</dict>
</plist>`

	writeFile(t, path, stringsPlist)
	keys, err := ReadAppIntentKeys(context.Background(), res)
	if err != nil {
		t.Fatalf("ReadAppIntentKeys: %v", err)
	}
	if len(keys) != 1 || keys[0].Key != "AppIntent.TaskList.AllTasks" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestReadAppShortcutKeys(t *testing.T) {
	root := t.TempDir()
	res := filepath.Join(root, "Resources")
	path := filepath.Join(res, "fr.lproj", "AppShortcuts.strings")

	stringsPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>#!SET#!_AppIntent.TaskList.ListOf${applicationName}[0]</key>
	<string>Liste des taches</string>
</dict>
</plist>`

	writeFile(t, path, stringsPlist)
	keys, err := ReadAppShortcutKeys(context.Background(), res)
	if err != nil {
		t.Fatalf("ReadAppShortcutKeys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestReadAppShortcutPhrases(t *testing.T) {
	root := t.TempDir()
	res := filepath.Join(root, "Resources")
	path := filepath.Join(res, "fr.lproj", "AppShortcuts.strings")

	stringsPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>#!SET#!_AppIntent.TaskList.ListOf${applicationName}[0]</key>
	<string>Liste des taches</string>
</dict>
</plist>`

	writeFile(t, path, stringsPlist)
	t.Setenv("LANG", "fr_FR.UTF-8")
	keys, err := ReadAppShortcutPhrases(context.Background(), res)
	if err != nil {
		t.Fatalf("ReadAppShortcutPhrases: %v", err)
	}
	if len(keys) != 1 || keys[0].Value == "" {
		t.Fatalf("unexpected phrases: %v", keys)
	}
}

func TestActionShortcutCandidates(t *testing.T) {
	def := ActionDef{
		ID:        "task-list",
		Title:     "List tasks",
		Transport: TransportShortcuts,
		Keys:      []string{"AppIntent.TaskList.AllTasks"},
	}
	app := AppInfo{Name: "Streaks"}
	intentKeys := []AppIntentKey{{Key: "AppIntent.TaskList.AllTasks", Value: "All Tasks"}}
	phrases := []AppIntentKey{{Key: "#!SET#!_AppIntent.TaskList.ListOf${applicationName}[0]", Value: "List of ${applicationName}"}}

	candidates := ActionShortcutCandidates(def, app, intentKeys, phrases, "")
	joined := strings.Join(candidates, "|")
	if !strings.Contains(joined, "All Tasks") || !strings.Contains(joined, "List of Streaks") {
		t.Fatalf("unexpected candidates: %v", candidates)
	}
}
