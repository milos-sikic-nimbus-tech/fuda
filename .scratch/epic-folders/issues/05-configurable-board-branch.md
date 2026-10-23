# 05: Configurable board branch

**What to build:** The branch fuda reads (and writes) the board from is a setting, default `develop`.

**Blocked by:** none

**Status:** needs-triage

Why: a separate tickets repo has no reason for a `develop` branch; it would use `main`. Today `develop` is fixed. "Watch main" for In prod also means little in such a repo (see ticket 04).

- Decide where the setting lives: env var, per board in the workspace config, or both. Per board is likely needed since one person can open several boards. Coordinate with `fuda-writes-build` ticket 25 (env variables review).
- If the board branch is `main`, `FUDA_WATCH_MAIN` should not read the same branch twice.

- [ ] The branch is configurable, default `develop`, for reads and writes
- [ ] A missing branch gives a clear error
- [ ] Guide `02-setup-a-repo.md`, `08-self-host.md`, `.env.example` and `docs/architecture.md` updated
