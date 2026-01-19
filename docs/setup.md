# Setup

This CLI only uses officially exposed automation surfaces. For Streaks on macOS,
that means the Shortcuts actions published by the app and the `streaks://` URL
scheme.

## Discovery

Run discovery any time Streaks updates:

```
st discover
```

Discovery reads the app bundle on disk (`Info.plist` + `Localizable.strings`) to
identify supported App Intents and URL schemes. It does **not** read the
Shortcuts database or reverse engineer binaries.

If new App Intent keys appear, re-run `st install` and update wrapper
shortcuts to match the new capabilities.

## Install

```
st install
```

This writes the config file and reports which wrapper shortcuts are missing. The
Shortcuts CLI does not provide a supported way to create shortcuts programmatically,
so wrappers must be created manually.

Optional checklist:

```
st install --checklist wrappers-checklist.txt
```

## Wrapper shortcuts (manual)

Wrapper shortcuts must:

1. Accept input from stdin (JSON).
2. Invoke the corresponding Streaks action.
3. Output a Dictionary (so `shortcuts run --output-type public.json` returns JSON).

Shortcut template:

1. **Get Dictionary from Input** (input is JSON).
2. **Streaks action** that matches the capability.
3. **Dictionary** output with keys like `ok`, `action`, `task`, `timestamp`.

For task-based actions, read `task` from the input dictionary. For pause, read
`status` if supplied.

### Wrapper names

Wrapper names are deterministic and listed in `~/.config/streaks-cli/config.json`.
Default names follow this pattern:

```
st <action-id>
```

To list wrappers and sample inputs:

```
st wrappers list
st wrappers sample task-complete
```

To validate wrappers:

```
st wrappers verify --task "Example Task"
```

### Action map

Match each wrapper to the Streaks action shown in Shortcuts:

- `task-complete` → “Mark [task] as complete” (input: `task`)
- `task-miss` → “Mark [task] as missed” (input: `task`)
- `task-status` → “Status of [task]” (input: `task`)
- `task-reminder` → “Reminder for [task]” (input: `task`)
- `task-list` → “All tasks” (no input required)
- `timer-start` → “Start [task] timer” (input: `task`)
- `timer-stop` → “Stop [task] timer” (input: `task`)
- `pause` → “Pause” (optional input: `status` = `All` or `NotPaused`)
- `resume` → “Resume” (no input required)
- `export-all` → “Export all data” (no input required)
- `export-task` → “Export [task] data” (input: `task`)

After creating wrappers, verify:

```
st doctor
```

See `docs/wrappers.md` for more detail.
