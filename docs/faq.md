# FAQ

## Streaks app not found
- Verify Streaks is installed in `/Applications/Streaks.app`.
- Re-run: `st discover`.

## Wrapper shortcut missing
- If you already have a matching Streaks shortcut, `st` will use it automatically.
- Otherwise run: `st wrappers list` to see expected names.
- Create missing shortcuts in Shortcuts.
- Re-run: `st doctor`.

## Shortcuts permission prompts
- The first run may prompt for automation permissions.
- Approve access for Shortcuts and Streaks when requested.

## Discover fails after update
- Re-run discovery and review `Unmapped App Intent Keys` in the output.
- Update wrappers and mapping if Streaks adds new actions.

## Output formats
- Use `--output json` or `--agent` for machine-readable output.
- Use `--output plain` for stable line-based output.
- Use `--no-output` when only exit codes matter.

## Wrapper validation fails
- Run `st wrappers verify --task "Example"`.
- Ensure wrapper shortcuts return a Dictionary/JSON.

## Homebrew install issues
- Ensure the tap is added: `brew tap vburojevic/tap`.
- Update: `brew update` then `brew upgrade streaks-cli`.
