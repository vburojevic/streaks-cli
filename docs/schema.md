# JSON Output Schema

Use `--output json` or `--agent` to enable JSON output.

## Exit codes

- `0` success
- `2` invalid usage
- `10` Streaks app not found
- `11` Shortcuts CLI missing or failed
- `12` Streaks shortcut missing
- `13` action execution failed

All JSON outputs are UTF-8 and printed to stdout. Errors are printed to stderr as:

```json
{"error":"message","code":10}
```

## `st discover`

```json
{
  "timestamp": "RFC3339",
  "app": {"name":"Streaks","path":"...","bundle_id":"...","version":"...","build":"...","resources_path":"..."},
  "url_schemes": ["streaks"],
  "ns_user_activity_types": ["..."],
  "shortcuts_cli_path": "/usr/bin/shortcuts",
  "shortcuts_cli_available": true,
  "app_intent_keys": [{"key":"...","value":"..."}],
  "app_shortcut_keys": ["..."],
  "app_shortcut_phrases": [{"key":"...","value":"..."}],
  "xcallback_supported": false,
  "actions": [{"id":"task-complete","title":"...","transport":"shortcuts","requires_task":true}],
  "unmapped_keys": ["..."],
  "notes": ["..."]
}
```

## `st doctor`

```json
{
  "app_installed": true,
  "app_path": "/Applications/Streaks.app",
  "bundle_id": "...",
  "version": "...",
  "shortcuts_cli": true,
  "shortcuts_cli_path": "/usr/bin/shortcuts",
  "shortcut_count": 12,
  "shortcut_actions_available": ["task-list"],
  "shortcut_actions_missing": [],
  "url_schemes": ["streaks"],
  "warnings": []
}
```

## `st install`

```json
{
  "shortcut_actions_available": ["task-list"],
  "shortcut_actions_missing": ["pause"],
  "note": "The CLI uses existing Streaks shortcuts..."
}
```

## `st actions list`

```json
[{"id":"task-complete","title":"...","transport":"shortcuts","requires_task":true,"parameters":{}}]
```

## `st actions describe`

```json
{"action":{"id":"...","title":"...","transport":"shortcuts","requires_task":true,"parameters":{}},"sample_input":{"task":"<task>"},"shortcut_candidates":["Complete ${task} in Streaks"]}
```

## Action commands (`st <action>`)

Shortcut output is passed through as JSON. If a shortcut returns non-JSON,
`--output json` wraps it as:

```json
{"raw":"...","format":"text","shortcut":"All Tasks"}
```

For `--dry-run`, output is:

```json
{"dry_run":true,"shortcut":"Complete Example in Streaks","input":{"task":"Example"}}
```
