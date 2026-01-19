# Command Reference

This is a concise reference for `st`. Run `st --help` or `st <command> --help` for full flag details.

## Global flags

- `--output human|json|plain` – output mode (default: human).
- `--json` – alias for `--output json`.
- `--plain` – alias for `--output plain`.
- `--agent` – agent mode (JSON output, no pretty formatting).
- `--quiet` / `--verbose` – reduce or increase output.
- `--timeout` – Shortcuts run timeout (default: 30s).
- `--retries` / `--retry-delay` – retry Shortcuts runs on failure.
- `--config` – override config path.

## Core commands

- `st discover` – print discovered capabilities (JSON by default).
- `st discover --markdown` – Markdown report (human output only).
- `st doctor` – verify Streaks install + wrapper readiness.
- `st install` – write config and report missing wrappers.
- `st open` – open Streaks via URL scheme.

## Wrappers

- `st wrappers list` – expected wrapper shortcuts.
- `st wrappers sample <action>` – JSON input template.
- `st wrappers verify --task "Example"` – run wrappers and validate JSON.
- `st wrappers doctor --task "Example"` – report readiness + optional validation.

## Actions

- `st actions list` – list all actions.
- `st actions describe <action>` – show details and sample input.
- `st <action>` – run a Streaks action (e.g., `st task-complete --task "Read"`).

## Action flags

- `--task` – task name for task-based actions.
- `--stdin` – force JSON input from stdin.
- `--input` – raw JSON input string.
- `--dry-run` – print wrapper + payload only.
- `--trace <file>` – append JSON trace records (JSONL).
