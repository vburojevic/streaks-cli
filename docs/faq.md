# FAQ

## Streaks app not found
- Verify Streaks is installed in `/Applications/Streaks.app`.
- Re-run: `streaks-cli discover`.

## Wrapper shortcut missing
- Run: `streaks-cli wrappers list` to see expected names.
- Create missing shortcuts in Shortcuts.
- Re-run: `streaks-cli doctor`.

## Shortcuts permission prompts
- The first run may prompt for automation permissions.
- Approve access for Shortcuts and Streaks when requested.

## Discover fails after update
- Re-run discovery and review `Unmapped App Intent Keys` in the output.
- Update wrappers and mapping if Streaks adds new actions.

## Homebrew install issues
- Ensure the tap is added: `brew tap vburojevic/tap`.
- Update: `brew update` then `brew upgrade streaks-cli`.

## Agent output parsing
- Use `--agent` or `--json` for machine-readable output.
- Error responses are JSON when `--agent` is enabled.
