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
st install
st doctor
st open
```

If wrappers are missing, follow `docs/wrappers.md` to create them in Shortcuts.

## Testing

```
make test
make smoke   # requires Streaks installed
```

## Common usage

```
# Run a task action
st task-complete --task "Read 20 pages"

# Agent-friendly JSON
st --agent discover

# Plain output for scripts
st --output plain wrappers list

# Dry-run to see wrapper payload
st task-complete --task "Read" --dry-run
```

## Output modes

```
st --output human   # default
st --output json    # structured output
st --output plain   # line-based output
```

`--json` is equivalent to `--output json`.

Errors are printed as JSON to stderr in JSON mode, e.g.:

```json
{"error":"message","code":10}
```

## Agent usage

Use `--agent` or `STREAKS_CLI_AGENT=1` for deterministic JSON output:

```
st --agent doctor
st --agent wrappers list
st --agent actions list
```

## Config

Config is stored at `~/.config/streaks-cli/config.json`.
Override via `STREAKS_CLI_CONFIG` or `--config`.

## Command reference

- `docs/commands.md` – full command/flag reference.
- `docs/schema.md` – JSON output schema.

## Docs

- `docs/setup.md` – discovery + setup
- `docs/wrappers.md` – wrapper creation, verify/doctor
- `docs/faq.md` – troubleshooting
- `docs/release.md` – release workflow

## License

MIT. See `LICENSE`.
