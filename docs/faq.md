# FAQ

## Streaks app not found
- Verify Streaks is installed in `/Applications/Streaks.app`.
- Re-run: `st discover`.

## Streaks shortcut missing
- Create the corresponding shortcut in the Shortcuts app.
- `st` will try known Streaks shortcut names even if they’re not listed, but
  Shortcuts still requires the shortcut to exist in your library.
- Re-run: `st doctor` to verify.
- You can also run a specific shortcut by name: `st <action> --shortcut "Name"`.
- If your shortcut names differ from the defaults, map them once:
  `st link task-list --shortcut "My Tasks"`.

## Shortcuts permission prompts
- The first run may prompt for automation permissions.
- Approve access for Shortcuts and Streaks when requested.

## Discover fails after update
- Re-run discovery and review `Unmapped App Intent Keys` in the output.
- Update action matching if Streaks adds new actions.

## Output formats
- Use `--agent` for machine-readable NDJSON output.
- Default output is human-friendly for meta commands and raw shortcut output for actions.
- Use `--no-output` when only exit codes matter.

## JSON output is not structured
- Direct shortcuts may return plain text.
- Build shortcuts that return a Dictionary for structured output.
- Use `--agent` to get a stable JSON wrapper for action outputs.
- If you need JSON output from Shortcuts, use `--shortcuts-output public.json`.
- If a shortcut outputs multiple files, `st` aggregates them into a JSON array.

## Homebrew install issues
- Ensure the tap is added: `brew tap vburojevic/tap`.
- Update: `brew update` then `brew upgrade streaks-cli`.
