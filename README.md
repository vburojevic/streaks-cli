# streaks-cli

A macOS command-line interface for the Streaks (Crunchy Bagel) app.

This CLI integrates with Streaks using **official automation surfaces only**:
Shortcuts actions (App Intents) and the `streaks://` URL scheme.

## Requirements

- macOS with Streaks installed
- Shortcuts app (`/usr/bin/shortcuts`)

## Build

```
go build -o bin/streaks-cli ./cmd/streaks-cli
```

## Quick start

```
./bin/streaks-cli discover
./bin/streaks-cli install
./bin/streaks-cli doctor
./bin/streaks-cli open
```

## Commands

- `discover` – print discovered capabilities as JSON
- `doctor` – verify Streaks installation and wrapper setup
- `install` – write config and report missing wrappers
- `wrappers` – list wrappers and sample inputs
- `open` – open Streaks via URL scheme
- `<action>` – run a Streaks action via wrapper shortcut

Run `streaks-cli --help` to see all action commands.

## Config

Config is stored at `~/.config/streaks-cli/config.json`.
Override via `STREAKS_CLI_CONFIG` or `--config`.

## Docs

See `docs/setup.md` for discovery details and wrapper shortcut setup.
See `docs/wrappers.md` and `docs/faq.md` for troubleshooting.

## License

MIT. See `LICENSE`.
