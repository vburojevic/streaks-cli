# Command Reference

This is a concise reference for `st`. Run `st --help` or `st <command> --help` for full flag details.

## Global flags

- `--output human|json|plain` – output mode (default: human).
- `--json` – alias for `--output json`.
- `--plain` – alias for `--output plain`.
- `--agent` – agent mode (JSON output, no pretty formatting).
- `--quiet` / `--verbose` – reduce or increase output.
- `--no-output` – suppress all output (exit code only).
- `--timeout` – Shortcuts run timeout (default: 30s).
- `--retries` / `--retry-delay` – retry Shortcuts runs on failure.

## Core commands

- `st discover` – print discovered capabilities (JSON by default).
- `st discover --markdown` – Markdown report (human output only).
- `st doctor` – verify Streaks install + shortcut readiness.
- `st install` – verify Streaks shortcuts are ready.
- `st open` – open Streaks via URL scheme.

## Actions

- `st actions list` – list all actions.
- `st actions describe <action>` – show details and sample input.
  - `--task` expands shortcut candidates for task-based actions.
- `st <action>` – run a Streaks action (e.g., `st task-complete --task "Read"`).
  - Uses existing Streaks shortcuts by default.

## Action flags

- `--task` – task name for task-based actions.
- `--stdin` – force JSON input from stdin.
- `--input` – raw JSON input string.
- `--dry-run` – print shortcut + payload only.
- `--trace <file>` – append JSON trace records (JSONL).
- `--shortcut <name-or-id>` – run a specific shortcut by name/identifier.
