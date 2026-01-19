# Release Process

This project uses GoReleaser to build macOS binaries and update the Homebrew tap.

## Prerequisites

- Set `HOMEBREW_TAP_GITHUB_TOKEN` in GitHub Actions with write access to `vburojevic/homebrew-tap`.

## Steps

1. Ensure tests pass:
   ```
   go test ./...
   ```
2. Tag a release:
   ```
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```
3. GitHub Actions runs the `release` workflow and publishes:
   - GitHub release artifacts (darwin amd64/arm64)
   - Updated Homebrew formula in the tap

## Local Dry Run

```
goreleaser release --snapshot --clean
```
