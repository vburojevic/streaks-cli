# streaks-cli

A macOS command-line interface for the [Streaks](https://streaksapp.com) app. The command name is `st`.

Fast, automation-first workflows for Shortcuts-backed Streaks actions, with a clean NDJSON agent mode.

This CLI integrates with Streaks using **official automation surfaces only**:
Shortcuts actions (App Intents) and the `streaks://` URL scheme.

## Install

### Homebrew

```
brew tap vburojevic/tap
brew install streaks-cli
```

### From source

```
make build
./bin/st --help
```

## Quick start

```
st discover
st doctor
st open
```

Need help?

```
st help
st help task-list
```

Agents: use `st --agent help` to get NDJSON help payloads.

If you already have Streaks shortcuts in your library, `st <action>` will use
them automatically. If not, create Streaks shortcuts in the Shortcuts app and
re-run `st doctor`.

If your shortcut names differ from the defaults, map them once:

```
st link task-list --shortcut "My Tasks"
st links
```

## Testing

```
make test
make smoke   # requires Streaks installed
```

## Common usage

```
# Run a task action
st task-complete --task "Read 20 pages"

# Run a specific shortcut directly
st task-list --shortcut "All Tasks"

# Agent-friendly JSON
st --agent discover

# Agent NDJSON
st --agent actions list

# Dry-run to see shortcut payload
st task-complete --task "Read" --dry-run
```

## Wrapper shortcut import (optional)

If you have exported wrapper `.shortcut` files locally (or use the bundled ones
in `./shortcuts`), import them in one step:

```
st install --import
```

Use `--from-dir` to point at a custom directory of `.shortcut` files.

Wrapper shortcut names are used exactly as file names. Helper shortcuts like
"Get Task Object" and "Get Task Details" are dependencies and not mapped to
CLI actions.

## Output modes

- Default: human-readable output for meta commands, raw shortcut output for actions.
- Agent mode: NDJSON (`--agent` or `STREAKS_CLI_AGENT=1`).
- `--no-output` suppresses all output (exit code only).

Default Shortcuts output is plain text. If you need JSON, set:

```
st --shortcuts-output public.json task-list
```

## Config

Mappings live at `~/.config/streaks-cli/config.json` by default. Override with:

```
st --config /path/to/config.json links
```

## Agent quick start

```
st --agent help
st --agent discover
st --agent actions list
st --agent task-list
```

If an agent is unsure how to call a command, it should run:

```
st --agent help <command>
```

Usage errors in agent mode include a `hint` field pointing to `st help`.

Errors are printed as JSON to stderr in agent mode, e.g.:

```json
{"error":"message","code":10}
```

## Agent usage

Use `--agent` or `STREAKS_CLI_AGENT=1` for NDJSON output:

```
st --agent doctor
st --agent actions list
st --agent task-list
```

Action commands emit a stable JSON envelope in agent mode.

## Command reference

- `docs/commands.md` – full command/flag reference.
- `docs/schema.md` – JSON output schema.

## Docs

- `docs/setup.md` – discovery + setup
- `docs/faq.md` – troubleshooting
- `docs/release.md` – release workflow

## License

MIT. See `LICENSE`.
