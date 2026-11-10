# 30: Write a comment on a task

**What to build:** A person writes a comment in the task sheet and fuda appends it to the task file, under `## Comments` at the bottom, starting with the date and the author (the convention in Guide page `06-agent-guide.md`). If the file has no `## Comments` heading, fuda adds it at the end. The write goes through the same path as Move and Assign (`api/internal/board`): one commit per comment, the same compare-and-swap, saving state and rollback in the UI. Works on every host that already supports writes (GitHub, Azure, local in the desktop app).

**Blocked by:** none

**Status:** needs-triage

Open questions:

- **Author:** the short name from `people.md` for the logged-in account, the host login, or both. What to write when the account has no `people.md` entry.
- **Line format:** exact shape, e.g. `- 2026-11-10 milos: text`, and how a multi-line comment is written so it stays inside its entry.
- **Conflict:** another commit changed the file since it was read. A comment only appends, so retry on the new file instead of failing like a Move conflict?
- **Edit and delete:** only new comments, or also editing or deleting your own.
- **Read-only boards:** hide the comment box, the same as Move and Assign.

- [ ] Pure edit tests: append under an existing `## Comments`, add the heading when missing, multi-line text
- [ ] Board service tests: comment happy path and conflict
- [ ] Comment box in the task sheet with saving and rollback like Move
- [ ] The new comment shows in the task sheet after the write
- [ ] No comment box on read-only boards
- [ ] Guide pages and `docs/architecture.md` updated
- [ ] `make check` passes

Spec: `.scratch/fuda-writes/spec.md`.
