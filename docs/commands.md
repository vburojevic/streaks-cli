# Command Reference

This is a concise reference for `st`. Run `st --help` or `st <command> --help` for full flag details.

## Global flags

- `--agent` – agent mode (NDJSON output; actions emit a stable envelope).
- `--quiet` / `--verbose` – reduce or increase output.
- `--no-output` – suppress all output (exit code only).
- `--timeout` – Shortcuts run timeout (default: 30s).
- `--retries` / `--retry-delay` – retry Shortcuts runs on failure.
- `--config` – override config path (default: `~/.config/streaks-cli/config.json`).
- `--shortcuts-output` – Shortcuts output UTI (default: `public.plain-text`, use `public.json` for JSON).

## Core commands

- `st discover` – print discovered capabilities (JSON by default).
- `st discover --markdown` – Markdown report (human output only).
- `st doctor` – verify Streaks install + shortcut readiness.
- `st install` – verify Streaks shortcuts are ready.
- `st install --import` – open bundled `.shortcut` wrapper files for import.
- `st link <action-id> --shortcut <name-or-id>` – map an action to a specific shortcut.
- `st unlink <action-id>` – remove action mapping.
- `st links` – list mappings.
- `st help [command]` – show help (agent mode returns NDJSON).
- `st open` – open Streaks via URL scheme.

## Actions

- `st actions list` – list all actions.
- `st actions describe <action>` – show details and sample input.
  - `--task` expands shortcut candidates for task-based actions.
- `st <action>` – run a Streaks action (e.g., `st task-complete --task "Read"`).
  - Uses existing Streaks shortcuts by default.

Action matching uses exact shortcut names. Helper shortcuts (e.g., "Get Task Object",
"Get Task Details") are not exposed as CLI actions.

## Action flags

- `--task` – task name for task-based actions.
- `--stdin` – force JSON input from stdin.
- `--input` – raw JSON input string.
- `--dry-run` – print shortcut + payload only.
- `--trace <file>` – append JSON trace records (JSONL).
- `--shortcut <name-or-id>` – run a specific shortcut by name/identifier.

## Install flags

- `--import` – open bundled `.shortcut` files for import.
- `--from-dir` – directory containing `.shortcut` files to import.

Environment:
- `STREAKS_CLI_SHORTCUT_DIR` – override default wrapper shortcut directory.

## Link flags

- `--shortcut` – shortcut name or identifier (for `st link`).
- `--shortcut-name` – shortcut name.
- `--shortcut-id` – shortcut identifier.
