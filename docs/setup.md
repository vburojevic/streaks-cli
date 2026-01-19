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

If new App Intent keys appear, re-run `st discover` and check the action list.

## Direct shortcuts (required)

`st` runs **existing Streaks shortcuts** in your Shortcuts library. Create those
shortcuts via the Shortcuts app (e.g., from Streaks “Add Shortcut” buttons).

To see candidates for a specific action:

```
st actions describe task-complete --task "Example Task"
```

Run a specific shortcut explicitly:

```
st task-list --shortcut "All Tasks"
```

## Install

```
st install
```

This verifies whether non-task actions already have matching shortcuts. Task-based
actions require a task name and are checked when you run them.

After creating shortcuts, verify:

```
st doctor
```

See `docs/commands.md` for more detail.
