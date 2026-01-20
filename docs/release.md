# Release Process

This project uses GoReleaser to build macOS binaries and update the Homebrew tap.

## Prerequisites

- Set `HOMEBREW_TAP_GITHUB_TOKEN` in GitHub Actions with write access to `vburojevic/homebrew-tap`.

## Steps

1. Ensure tests pass and the working tree is clean:
   ```
   git status --short
   go test ./...
   ```
2. Commit all changes for the release.
3. Tag a release:
   ```
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```
4. Push `main` (if needed):
   ```
   git push origin main
   ```
5. GitHub Actions runs the `release` workflow and publishes:
   - GitHub release artifacts (darwin amd64/arm64)
   - Updated Homebrew formula in the tap

6. Verify the release:
   ```
   gh release view vX.Y.Z
   gh run list --limit 5
   ```

7. Verify the Homebrew tap points at the new version. If it didn’t update,
   update it manually:
   ```
   curl -L -o /tmp/streaks-cli-vX.Y.Z.tar.gz \
     https://github.com/vburojevic/streaks-cli/archive/refs/tags/vX.Y.Z.tar.gz
   shasum -a 256 /tmp/streaks-cli-vX.Y.Z.tar.gz
   # Update Formula/streaks-cli.rb in vburojevic/homebrew-tap with the new url+sha.
   ```

## Local Dry Run

```
goreleaser release --snapshot --clean
```
