---
status: accepted
---

# A board is its own repository, read from `main`, whatever its name

A board is a tasks repository, and fuda reads and writes only its `main`. The code repository keeps no copy of the Tasks, not even a submodule. fuda lists the `fuda-` repositories a person can read, plus the repositories named in `FUDA_GITHUB_BOARDS` and `FUDA_AZURE_BOARDS`, which need no prefix. Agreed with the upstream maintainer on 25 November 2026. This narrows ADR 0002, which still mounted the board as a submodule.

## Why

- One branch means one Status per Task. Upstream read `docs/board` on `develop`, with `main` only adding an "in prod" badge, so a Task had two Statuses to reason about.
- The submodule kept showing as modified in the code repository and was empty in fresh clones.
- Some real boards are not named `fuda-<project>`; the upstream maintainer needs those to work, so the listed repositories stay.

## Considered Options

- `develop` as the board branch with an "in prod" badge from `main`: rejected, two Statuses per Task.
- Require the `fuda-` prefix: rejected, it breaks existing boards with other names.

## Consequences

- The "in prod" badge and `FUDA_WATCH_MAIN` are gone. "In prod" may come back later only as a label the user sets (TBD in `decisions.md`).
- A listed repository without `main` shows the same error page as a `fuda-` one.
- Moving an old `develop` board is the first step of the 0.1 to 0.2 migration (ADR 0011).
