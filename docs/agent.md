# Agent Guide

This CLI is designed for both humans and automation. Agents should always run
with `--agent` (or set `STREAKS_CLI_AGENT=1`) to get NDJSON output and stable
error payloads.

## Recommended pattern

1. **Discover once** on startup and cache the result:
   - `st --agent discover`
2. **Use `st help` when uncertain**:
   - `st --agent help`
   - `st --agent help task-list`
3. **Run actions**, parse the last NDJSON line as the result:
   - `st --agent task-list`
4. **Check errors** on stderr:
   - `{"error":"...","code":12,"error_code":"shortcut_missing"}`

## Output contract (agent mode)

- **stdout**: NDJSON objects (one per line).
- **stderr**: NDJSON error objects (single line).
- **Exit codes**: use numeric `code` plus string `error_code`.

### Error fields

- `code`: numeric exit code.
- `error_code`: stable string for automation.
- `hint`: present for usage errors (points to `st help`).

## Shortcuts best practices

- Prefer Shortcuts that **return a Dictionary** for structured output.
- If a shortcut returns multiple output files, `st` aggregates them into a JSON
  array.
- Use `--shortcuts-output public.json` when you need JSON payloads from
  Shortcuts (defaults to `public.plain-text`).

## Safe defaults for agents

```
st --agent --timeout 30s --retries 1 task-list
```

## Example integration (bash)

```bash
set -euo pipefail

out=$(st --agent task-list)
echo "$out" | tail -n 1 | jq -r '.[].title'
```

## Example integration (JSONL consumer)

```bash
st --agent task-list | while read -r line; do
  echo "$line" | jq -e '.result' >/dev/null
done
```
