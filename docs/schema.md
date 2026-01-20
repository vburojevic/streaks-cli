# JSON Output Schema

Use `--agent` (or `STREAKS_CLI_AGENT=1`) to enable NDJSON output.

## Exit codes

- `0` success
- `2` invalid usage
- `10` Streaks app not found
- `11` Shortcuts CLI missing or failed
- `12` Streaks shortcut missing
- `13` action execution failed

NDJSON outputs are UTF-8 JSON objects printed one per line to stdout. Errors are printed to stderr as:

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
  "app_intent_keys": [{"key":"...","value":"...","locale":"en"}],
  "app_shortcut_keys": ["..."],
  "app_shortcut_phrases": [{"key":"...","value":"...","locale":"en"}],
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
  "note": "The CLI uses existing Streaks shortcuts...",
  "import_dir": "shortcuts",
  "imported": ["shortcuts/Streaks API - List.shortcut"],
  "import_errors": [],
  "import_warning": ""
}
```

## `st links`

NDJSON: one mapping per line.

```json
{"path":"/Users/me/.config/streaks-cli/config.json","action":"task-list","shortcut":{"name":"All Tasks"}}
```

## `st link` / `st unlink`

```json
{
  "path": "/Users/me/.config/streaks-cli/config.json",
  "action": "task-list",
  "shortcut": {"name":"All Tasks"}
}
```

## `st actions list`

NDJSON: one action per line.

```json
{"id":"task-complete","title":"...","transport":"shortcuts","requires_task":true,"parameters":{}}
```

## `st actions describe`

```json
{"action":{"id":"...","title":"...","transport":"shortcuts","requires_task":true,"parameters":{}},"sample_input":{"task":"<task>"},"shortcut_candidates":["Complete ${task} in Streaks"],"mapped_shortcut":{"name":"All Tasks"}}
```

## Action commands (`st <action>`)

Shortcut output is passed through as-is. If a shortcut returns non-JSON,
agent mode wraps it as:

```json
{"raw":"...","format":"text","shortcut":"All Tasks"}
```

With `--agent`, action output is wrapped in a stable envelope:

```json
{
  "ok": true,
  "timestamp": "RFC3339Nano",
  "action": {"id":"task-list"},
  "shortcut": {"name":"All Tasks"},
  "attempts": 1,
  "duration_ms": 12,
  "input": {"task":"Example"},
  "result": {"raw":"...","format":"text","shortcut":"All Tasks"}
}
```

For `--dry-run`, output is:

```json
{"dry_run":true,"shortcut":"Complete Example in Streaks","input":{"task":"Example"}}
```

## Task field schema (custom shortcut output)

If you build shortcuts that output task data as a dictionary, the following
keys are available from the Streaks task schema (values are strings):

```json
{
  "completed_segments": "<string>",
  "number_of_segments": "<string>",
  "timer_finish_time": "<string>",
  "subprogress_percent": "<string>",
  "best_streak": "<string>",
  "will_miss_today": "<string>",
  "progress_percent": "<string>",
  "title": "<string>",
  "is_paused": "<string>",
  "current_streak": "<string>",
  "total_duration": "<string>",
  "is_complete": "<string>",
  "duration_per_segment": "<string>",
  "is_missed": "<string>",
  "remaining_duration": "<string>",
  "timer_is_timing": "<string>",
  "segment_remaining_duration": "<string>",
  "today_status": "<string>"
}
```
