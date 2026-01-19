package discovery

type AppInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	BundleID  string `json:"bundle_id"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	Resources string `json:"resources_path"`
}

type AppIntentKey struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type Action struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Transport    string              `json:"transport"`
	RequiresTask bool                `json:"requires_task"`
	Parameters   map[string][]string `json:"parameters,omitempty"`
	SourceKeys   []string            `json:"source_keys"`
}

type Discovery struct {
	Timestamp             string         `json:"timestamp"`
	App                   AppInfo        `json:"app"`
	URLSchemes            []string       `json:"url_schemes"`
	NSUserActivityTypes   []string       `json:"ns_user_activity_types"`
	ShortcutsCLIPath      string         `json:"shortcuts_cli_path"`
	ShortcutsCLIAvailable bool           `json:"shortcuts_cli_available"`
	AppIntentKeys         []AppIntentKey `json:"app_intent_keys"`
	AppShortcutKeys       []string       `json:"app_shortcut_keys"`
	Actions               []Action       `json:"actions"`
	UnmappedKeys          []string       `json:"unmapped_keys,omitempty"`
	Notes                 []string       `json:"notes,omitempty"`
}
