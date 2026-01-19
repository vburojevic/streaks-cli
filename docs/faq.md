# FAQ

## Streaks app not found
- Verify Streaks is installed in `/Applications/Streaks.app`.
- Re-run: `st discover`.

## Streaks shortcut missing
- Create the corresponding shortcut in the Shortcuts app.
- Re-run: `st doctor` to verify.
- You can also run a specific shortcut by name: `st <action> --shortcut "Name"`.

## Shortcuts permission prompts
- The first run may prompt for automation permissions.
- Approve access for Shortcuts and Streaks when requested.

## Discover fails after update
- Re-run discovery and review `Unmapped App Intent Keys` in the output.
- Update action matching if Streaks adds new actions.

## Output formats
- Use `--output json` or `--agent` for machine-readable output.
- Use `--output plain` for stable line-based output.
- Use `--no-output` when only exit codes matter.

## JSON output is not structured
- Direct shortcuts may return plain text.
- Use `--output json` to wrap text in a JSON object, or build shortcuts that return a Dictionary.

## Homebrew install issues
- Ensure the tap is added: `brew tap vburojevic/tap`.
- Update: `brew update` then `brew upgrade streaks-cli`.
