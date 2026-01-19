# JSON Output Schema

Use `--output json` or `--agent` to enable JSON output.

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
  "config_path": "~/.config/streaks-cli/config.json",
  "config_present": true,
  "wrapper_shortcuts": ["st task-complete"],
  "missing_wrappers": [],
  "url_schemes": ["streaks"],
  "warnings": []
}
```

## `st install`

```json
{
  "config_path": "~/.config/streaks-cli/config.json",
  "checklist_path": "wrappers-checklist.txt",
  "missing_wrappers": ["st task-complete"]
}
```

## `st wrappers list`

```json
[{"id":"task-complete","title":"...","wrapper":"st task-complete","requires_task":true,"parameters":{}}]
```

## `st wrappers verify`

```json
[{"id":"task-complete","wrapper":"st task-complete","exists":true,"output_valid":true,"skipped":false,"error":""}]
```

## `st wrappers doctor`

```json
{"config_path":"...","missing_wrappers":[],"verify_results":[...]}
```

## `st actions list`

```json
[{"id":"task-complete","title":"...","transport":"shortcuts","requires_task":true,"parameters":{}}]
```

## `st actions describe`

```json
{"action":{"id":"...","title":"...","transport":"shortcuts","requires_task":true,"parameters":{}},"wrapper":"st task-complete","sample_input":{"task":"<task>"}}
```

## Action commands (`st <action>`)

Wrapper output is passed through as JSON. For `--dry-run`, output is:

```json
{"dry_run":true,"wrapper":"st task-complete","input":{"task":"Example"}}
```
