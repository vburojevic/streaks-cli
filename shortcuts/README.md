# Wrapper Shortcuts

This folder is reserved for exported `.shortcut` files that wrap Streaks actions.

The CLI can open these files for import:

```
st install --import
```

To supply your own wrappers:

1. Export the shortcuts from the Shortcuts app (File → Export).
2. Save the `.shortcut` files into this `shortcuts/` directory.
3. Run `st install --import` to open them for import.

The CLI will open each `.shortcut` file in Shortcuts so you can confirm the import.

Note: Helper shortcuts like "Get Task Object" and "Get Task Details" are
dependencies for other wrappers and are not mapped to CLI actions.
