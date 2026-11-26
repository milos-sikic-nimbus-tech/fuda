# 05: Boards read main

**What to build:** Remote boards (GitHub, Azure) read and write the branch named `main`. The fixed `develop` branch goes away.

**Blocked by:** none

**Status:** resolved

Why: `develop` came from the time when tasks lived in the code repo and were delivered in the same PR as the code. A board is now its own tickets repo (for example `fuda-eliot`), and `main` is its only branch.

- `workBranch = "develop"` in `api/internal/workspace/boards.go` becomes `main`, for reads, writes and sync. The branch is literally `main`, not the repo's default branch.
- Local boards are unchanged: they read the folder's files on disk, whatever is checked out.
- Repos that keep tasks on `develop` are no longer supported. fuda reads `main` and does not warn about a stale `develop`.
- **No `main`:** a repo whose name starts with `fuda-` and has no `main` stays in the board picker; opening it shows an error page that says the repo has no `main` branch and how to create it. Repos without the prefix are ticket 10.
- **In prod:** turned off. `FUDA_WATCH_MAIN` and the develop/main comparison no longer apply. Add a TBD entry to `docs/decisions.md`: In prod needs a new meaning, likely from code-repo PRs (tickets 04 and 06).
- **Protected `main`:** when a move's commit is refused, the move fails with a clear message (main is protected; allow fuda to push to it directly).

- [x] Remote boards read, sync and write `main`; tests updated
- [x] A `fuda-` repo without `main` opens with the error page
- [x] In prod is off, and the TBD entry is in `docs/decisions.md`
- [x] A refused commit on a protected `main` shows a clear error
- [ ] Guide and README changes are in ticket 09, which lands right after this one
