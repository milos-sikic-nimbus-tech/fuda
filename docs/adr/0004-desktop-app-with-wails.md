---
status: proposed
---

# A desktop app with Wails, next to the web server

The web shell cannot be hosted at $0 (no free host has a disk and no sleep), so fuda adds a desktop app that costs nothing to run. It is built with Wails v3 and mounts the same Go HTTP handler and React UI as the web server: one Go core, two shells. The web shell stays supported, but fuda picks no host; whoever wants it hosts it and pays.

## Consequences

- Unsigned builds warn on macOS and Windows. Avoiding the warning costs $99 a year on macOS. Windows signing is free through SignPath for open source.
- Updates come from the Wails updater reading free GitHub Releases. Wails v3 is still beta.
