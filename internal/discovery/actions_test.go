package discovery

import "testing"

func TestDetectActions(t *testing.T) {
	keys := []string{
		"AppIntent.TaskComplete.Mark${task}AsComplete",
		"AppIntent.TaskMiss.Mark${task}AsMissed",
		"AppIntent.TaskList.AllTasks",
		"AppIntent.Pause.Title",
		"AppIntent.Unknown.NewAction",
	}
	actions, unmapped := DetectActions(keys)

	wantIDs := map[string]bool{
		"task-complete": true,
		"task-miss":     true,
		"task-list":     true,
		"pause":         true,
	}
	for _, action := range actions {
		if action.Transport == TransportURLScheme {
			continue
		}
		if !wantIDs[action.ID] {
			t.Fatalf("unexpected action %s", action.ID)
		}
		delete(wantIDs, action.ID)
	}
	if len(wantIDs) != 0 {
		t.Fatalf("missing actions: %v", wantIDs)
	}

	if len(unmapped) != 1 || unmapped[0] != "AppIntent.Unknown.NewAction" {
		t.Fatalf("unexpected unmapped keys: %v", unmapped)
	}
}
