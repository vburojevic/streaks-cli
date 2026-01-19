# streaks-cli

A macOS command-line interface for the Streaks (Crunchy Bagel) app. The command name is `st`.

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

If you already have Streaks shortcuts in your library, `st <action>` will use
them automatically. If not, create Streaks shortcuts in the Shortcuts app and
re-run `st doctor`.

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

# Plain output for scripts
st --output plain actions list

# Dry-run to see shortcut payload
st task-complete --task "Read" --dry-run
```

## Output modes

```
st --output human   # default
st --output json    # structured output
st --output plain   # line-based output
```

`--json` is equivalent to `--output json`.
`--no-output` suppresses all output (exit code only).

Errors are printed as JSON to stderr in JSON mode, e.g.:

```json
{"error":"message","code":10}
```

## Agent usage

Use `--agent` or `STREAKS_CLI_AGENT=1` for deterministic JSON output:

```
st --agent doctor
st --agent actions list
st --agent actions list
```

## Command reference

- `docs/commands.md` – full command/flag reference.
- `docs/schema.md` – JSON output schema.

## Docs

- `docs/setup.md` – discovery + setup
- `docs/faq.md` – troubleshooting
- `docs/release.md` – release workflow

## License

MIT. See `LICENSE`.
