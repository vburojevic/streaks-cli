# Wrapper Shortcuts Guide

Wrapper shortcuts connect `st` to Streaks’ Shortcuts actions.

## List required wrappers

```
st wrappers list
```

## Generate input templates

```
st wrappers sample task-complete
```

## Validate wrapper output

```
st wrappers verify --task "Example Task"
```

## Wrapper doctor

```
st wrappers doctor --task "Example Task"
```

## Checklist file

```
st install --checklist wrappers-checklist.txt
```

## Manual creation steps

1. Create a new shortcut with the name shown in `wrappers list`.
2. Add **Get Dictionary from Input** (JSON).
3. Add the corresponding **Streaks** action.
4. Read `task` (and `status` for pause) from the input dictionary.
5. Return a **Dictionary** as output (e.g. `ok`, `action`, `task`, `timestamp`).

## Action mapping

- `task-complete` → “Mark [task] as complete” (input: `task`)
- `task-miss` → “Mark [task] as missed” (input: `task`)
- `task-status` → “Status of [task]” (input: `task`)
- `task-reminder` → “Reminder for [task]” (input: `task`)
- `task-list` → “All tasks” (no input)
- `timer-start` → “Start [task] timer” (input: `task`)
- `timer-stop` → “Stop [task] timer” (input: `task`)
- `pause` → “Pause” (optional input: `status` = `All` or `NotPaused`)
- `resume` → “Resume” (no input)
- `export-all` → “Export all data” (no input)
- `export-task` → “Export [task] data” (input: `task`)
