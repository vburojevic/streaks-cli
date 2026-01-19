# streaks-cli

A macOS command-line interface for the Streaks (Crunchy Bagel) app. The command name is `st`.

This CLI integrates with Streaks using **official automation surfaces only**:
Shortcuts actions (App Intents) and the `streaks://` URL scheme.

## Requirements

- macOS with Streaks installed
- Shortcuts app (`/usr/bin/shortcuts`)

## Build

```
go build -o bin/st ./cmd/streaks-cli
```

## Quick start

```
./bin/st discover
./bin/st install
./bin/st doctor
./bin/st open
```

## Commands

- `discover` – print discovered capabilities as JSON
- `doctor` – verify Streaks installation and wrapper setup
- `install` – write config and report missing wrappers
- `wrappers` – list wrappers and sample inputs
- `wrappers verify` – validate wrappers return JSON
- `wrappers doctor` – full wrapper readiness report
- `actions` – list and describe available actions
- `open` – open Streaks via URL scheme
- `<action>` – run a Streaks action via wrapper shortcut

Run `st --help` to see all action commands.

## Config

Config is stored at `~/.config/streaks-cli/config.json`.
Override via `STREAKS_CLI_CONFIG` or `--config`.

## AI Agent Usage

Use structured output for agents:

```
st --agent discover
st --agent doctor
st --agent wrappers list
st --agent actions list
```

`--agent` implies `--json` and disables pretty formatting for stable parsing.
On failures, the CLI returns JSON errors to stderr (e.g. `{"error":"...","code":10}`).
You can also set `STREAKS_CLI_AGENT=1`.

## Output Modes

```
st --output human   # default
st --output json    # structured output
st --output plain   # line-based output
```

`--json` is equivalent to `--output json`.

## Action Inputs

- Use `--stdin` to force stdin JSON input.
- Use `--dry-run` to print the wrapper and payload without running.
- Use `--trace <file>` to append JSON trace records of input/output.
- Use `--timeout`, `--retries`, `--retry-delay` to control Shortcuts execution.

## Docs

See `docs/setup.md` for discovery details and wrapper shortcut setup.
See `docs/wrappers.md` and `docs/faq.md` for troubleshooting.
See `docs/schema.md` for JSON output schema.

## License

MIT. See `LICENSE`.
