# 03: Boards without the fuda- prefix work like any other board

**What to build:** A repository named in `FUDA_GITHUB_BOARDS` or `FUDA_AZURE_BOARDS` is a board even without the `fuda-` prefix, and fuda reads and writes its `main`, the same as a `fuda-` board. Works in the server and in the desktop app.

**Blocked by:** 02

**Status:** ready-for-agent

Why: Nemanja's own setup has boards that are not named `fuda-<project>`. Upstream's F-10 lists them; the fork's rule (idea 1) decides which branch is read. This ticket makes the two meet.

- [ ] A listed repo without the prefix shows in the board list (GitHub and Azure)
- [ ] fuda reads and writes its `main`; a listed repo with no `main` shows the fork's "no main branch" error page
- [ ] The desktop app and release build carry the settings, as in upstream's F-10
- [ ] README and Guide (`01-overview`, `02-setup-a-repo`, `04-board-config`) say a board is a `fuda-` repo **or** a listed repo, and stop saying the prefix is required
- [ ] Tests for a listed board without the prefix, on both hosts
