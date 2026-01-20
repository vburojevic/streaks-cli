package cli

var wrapperShortcutNames = map[string][]string{
	"task-complete": {"Complete Task"},
	"task-miss":     {"Mark Task Missed"},
	"task-list":     {"Task List"},
	"task-status":   {"Get Task"},
	"task-reminder": {"Send Task Reminder"},
	"timer-start":   {"Start Task Timer"},
	"timer-stop":    {"Stop Task Timer"},
	"pause":         {"Pause Tasks"},
	"export-all":    {"Export Data"},
}

func addWrapperCandidates(actionID string, candidates []string) []string {
	wrappers, ok := wrapperShortcutNames[actionID]
	if !ok || len(wrappers) == 0 {
		return candidates
	}
	seen := make(map[string]struct{}, len(candidates))
	for _, cand := range candidates {
		seen[cand] = struct{}{}
	}
	for _, wrapper := range wrappers {
		if _, ok := seen[wrapper]; ok {
			continue
		}
		candidates = append(candidates, wrapper)
	}
	return candidates
}
