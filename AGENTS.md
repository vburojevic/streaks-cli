# Repository Guidelines

## Project Structure & Module Organization
- `cmd/streaks-cli/`: CLI entrypoint (`main.go`).
- `internal/cli/`: Cobra commands and CLI wiring.
- `internal/discovery/`: Local capability discovery from the Streaks app bundle.
- `internal/shortcuts/`: Shortcuts CLI integration.
- `internal/config/`: Config read/write (`~/.config/streaks-cli/config.json`).
- `internal/output/`: JSON output helpers.
- `docs/`: Contributor and setup docs (`docs/setup.md`).
- `bin/`: Local build output (ignored by git).

## Build, Test, and Development Commands
- `go build -o bin/streaks-cli ./cmd/streaks-cli` — build the CLI binary.
- `go test ./...` — run all unit tests.
- `goreleaser release --clean` — build and publish releases (CI only).
- `bin/streaks-cli discover` — print discovered capabilities.
- `bin/streaks-cli doctor` — verify Streaks + wrapper setup.
- `bin/streaks-cli install` — write config and report missing wrappers.

## Coding Style & Naming Conventions
- Go formatting: run `gofmt -w` on modified `.go` files.
- Package names are lowercase, short, and descriptive (e.g., `discovery`).
- Command names are kebab-case and map to discovered actions (e.g., `task-complete`).
- Use explicit error messages; return errors instead of panicking.

## Testing Guidelines
- Framework: Go’s built-in testing (`*_test.go`).
- Run `go test ./...` before committing.
- Tests should avoid calling the real Streaks app or Shortcuts; use stubs/mocks.
- CLI tests may set `STREAKS_CLI_DISABLE_DISCOVERY=1` to avoid discovery I/O.

## Commit & Pull Request Guidelines
- Commit messages follow Conventional Commits (e.g., `feat: ...`, `chore: ...`, `test: ...`).
- PRs should include: summary, test results, and any new discovery mappings.
- If discovery logic changes, update `docs/setup.md` with user-facing steps.

## Security & Configuration Notes
- Do not read sandboxed databases or reverse-engineer binaries.
- Automation is limited to official surfaces (Shortcuts, URL scheme).
- Config path is `~/.config/streaks-cli/config.json` (override with `STREAKS_CLI_CONFIG`).

## Release Workflow
- Tag releases with `vX.Y.Z` (e.g., `v0.2.0`), then push the tag.
- GitHub Actions runs GoReleaser to publish binaries and update the Homebrew tap.
- Requires secret `HOMEBREW_TAP_GITHUB_TOKEN` with write access to `vburojevic/homebrew-tap`.
