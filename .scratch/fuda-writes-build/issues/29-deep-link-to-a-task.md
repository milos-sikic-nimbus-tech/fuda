# 29: Deep link that opens a task in the desktop app

**What to build:** fuda is not hosted, so there is no https link to send a colleague. Add a custom URL scheme, for example `fuda://github/<owner>/<repo>/task/<id>`, that the desktop app registers with the OS. Clicking the link opens fuda (or brings the running window to the front), switches to that board and opens the task sheet, the same as the web app's `?task=<id>` search param. The task sheet gets a "Copy link" action that puts this link on the clipboard.

**Blocked by:** none

**Status:** needs-triage

Open questions:

- **Scheme and shape:** `fuda://` with host, repo and task id. Azure needs org/project/repo; decide the path for each host.
- **Board not open or not allowed:** what the app shows when the link names a board the user has not picked, cannot read, or is not logged in to.
- **Local boards:** a local folder has no shared address, so probably no link for local boards (hide "Copy link").
- **Single instance:** a link while the app runs must go to the running window, not start a second app.

Notes from the code: the app is Wails v3 (`api/cmd/fuda-desktop`). macOS needs `CFBundleURLTypes` in `api/cmd/fuda-desktop/Info.plist`; Windows needs the scheme in the registry, which the bare `.exe` does not do today. The board is not part of the client URL yet, only `task` is (`client/src/lib/filters.ts`).

- [ ] Clicking a `fuda://` link on macOS opens the app on that board with the task sheet open (checked by hand)
- [ ] Same on Windows (checked by hand)
- [ ] A link while the app runs reuses the open window
- [ ] The task sheet has "Copy link" for GitHub and Azure boards
- [ ] A bad or unreachable link shows a clear message, not a broken board
- [ ] Unit tests for parsing the link into board and task
- [ ] README, Guide and `docs/architecture.md` describe the links
- [ ] `make check` passes

Spec: `.scratch/fuda-writes/spec.md`.
