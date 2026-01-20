# Repository Guidelines

## Project Structure & Module Organization
- `cmd/streaks-cli/`: CLI entrypoint.
- `internal/cli/`: Cobra commands and exit-code handling.
- `internal/discovery/`: App bundle discovery (Info.plist + Localizable.strings).
- `internal/config/`: Config persistence for action-to-shortcut mappings.
- `internal/shortcuts/`: Shortcuts CLI integration.
- `internal/output/`: JSON helpers.
- `internal/xcallback/`: Stub for future x-callback URL support.
- `docs/`: `setup.md`, `release.md`, `faq.md`, `schema.md`.
- `shortcuts/`: Exported `.shortcut` wrapper files (optional).
- `.github/workflows/`: CI workflows.

## Build, Test, and Development Commands
- `make build` — build `bin/st`.
- `make test` — run all unit tests.
- `make lint` — run golangci-lint.
- `make integration` — run integration tests (requires Streaks).
- `goreleaser release --snapshot --clean` — local dry run.

## Coding Style & Naming Conventions
- Use `gofmt -w` on modified `.go` files.
- Package names are lowercase and descriptive.
- Command names are kebab-case (e.g., `task-complete`).
- Prefer explicit error messages over panics.

## Testing Guidelines
- Framework: Go `*_test.go`.
- Unit tests should not call Streaks or Shortcuts.
- Integration tests are gated by `STREAKS_CLI_INTEGRATION=1`.

## Commit & Pull Request Guidelines
- Use Conventional Commits (e.g., `feat:`, `chore:`, `test:`).
- PRs should include summary + test results.
- If discovery mappings change, update `docs/setup.md`.

## Security & Configuration Notes
- Do not read sandboxed databases or reverse engineer binaries.
- Automation uses official surfaces only (Shortcuts, URL scheme).
- Agent mode: `--agent` or `STREAKS_CLI_AGENT=1` for NDJSON output (action envelope).
- Action execution requires existing Streaks shortcuts; use `--shortcut` or `st link` mappings.
- Config path: `~/.config/streaks-cli/config.json` (override with `--config` or `STREAKS_CLI_CONFIG`).

## Release Workflow
- Tag `vX.Y.Z`, push the tag.
- GitHub Actions runs GoReleaser and updates the Homebrew tap.
- Requires `HOMEBREW_TAP_GITHUB_TOKEN` secret with write access to `vburojevic/homebrew-tap`.
