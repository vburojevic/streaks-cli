package discovery

import "strings"

type ActionDef struct {
	ID           string
	Title        string
	Transport    string
	RequiresTask bool
	Keys         []string
	ParamOptions map[string][]string
}

const (
	TransportShortcuts = "shortcuts"
	TransportURLScheme = "url-scheme"
)

func DefaultActionDefinitions() []ActionDef {
	return []ActionDef{
		{
			ID:           "open",
			Title:        "Open Streaks",
			Transport:    TransportURLScheme,
			RequiresTask: false,
			Keys:         nil,
		},
		{
			ID:           "task-complete",
			Title:        "Mark task complete",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.TaskComplete.Mark${task}AsComplete",
				"AppIntent.TaskComplete.Mark%@AsComplete",
			},
		},
		{
			ID:           "task-miss",
			Title:        "Mark task missed",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.TaskMiss.Mark${task}AsMissed",
			},
		},
		{
			ID:           "task-list",
			Title:        "List tasks",
			Transport:    TransportShortcuts,
			RequiresTask: false,
			Keys: []string{
				"AppIntent.TaskList.AllTasks",
			},
		},
		{
			ID:           "task-status",
			Title:        "Task status",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.Status.StatusOf${task}",
			},
		},
		{
			ID:           "task-reminder",
			Title:        "Task reminder",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.Notification.ReminderFor${task}",
			},
		},
		{
			ID:           "timer-start",
			Title:        "Start task timer",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.StartTimer.Start${task}Timer",
				"AppIntent.StartTimer.Start%@Timer",
			},
		},
		{
			ID:           "timer-stop",
			Title:        "Stop task timer",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.StopTimer.Stop${task}Timer",
			},
		},
		{
			ID:           "pause",
			Title:        "Pause tasks",
			Transport:    TransportShortcuts,
			RequiresTask: false,
			Keys: []string{
				"AppIntent.Pause.Title",
			},
			ParamOptions: map[string][]string{
				"status": {"All", "NotPaused"},
			},
		},
		{
			ID:           "export-all",
			Title:        "Export all data",
			Transport:    TransportShortcuts,
			RequiresTask: false,
			Keys: []string{
				"AppIntent.DataExport.ExportAllData",
			},
		},
		{
			ID:           "export-task",
			Title:        "Export task data",
			Transport:    TransportShortcuts,
			RequiresTask: true,
			Keys: []string{
				"AppIntent.DataExport.Export${task}Data",
			},
		},
	}
}

func DetectActions(keys []string) ([]Action, []string) {
	keySet := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keySet[key] = struct{}{}
	}

	defs := DefaultActionDefinitions()
	actions := make([]Action, 0, len(defs))
	mappedKeys := make(map[string]struct{})

	for _, def := range defs {
		if def.Transport == TransportURLScheme {
			actions = append(actions, Action{
				ID:           def.ID,
				Title:        def.Title,
				Transport:    def.Transport,
				RequiresTask: def.RequiresTask,
				Parameters:   def.ParamOptions,
				SourceKeys:   def.Keys,
			})
			continue
		}

		present := false
		for _, k := range def.Keys {
			if _, ok := keySet[k]; ok {
				present = true
				mappedKeys[k] = struct{}{}
			}
		}
		if !present {
			continue
		}
		actions = append(actions, Action{
			ID:           def.ID,
			Title:        def.Title,
			Transport:    def.Transport,
			RequiresTask: def.RequiresTask,
			Parameters:   def.ParamOptions,
			SourceKeys:   def.Keys,
		})
	}

	unmapped := make([]string, 0)
	for _, key := range keys {
		if _, ok := mappedKeys[key]; !ok {
			if strings.HasPrefix(key, "AppIntent.") && !strings.Contains(key, "PauseStatus") {
				unmapped = append(unmapped, key)
			}
		}
	}
	return actions, unmapped
}
