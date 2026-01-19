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
```

`--agent` implies `--json` and disables pretty formatting for stable parsing.
On failures, the CLI returns JSON errors to stderr (e.g. `{"error":"...","code":10}`).
You can also set `STREAKS_CLI_AGENT=1`.

## Docs

See `docs/setup.md` for discovery details and wrapper shortcut setup.
See `docs/wrappers.md` and `docs/faq.md` for troubleshooting.

## License

MIT. See `LICENSE`.
