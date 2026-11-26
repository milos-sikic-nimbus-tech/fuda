# 09: Board format 0.2 and the migration from 0.1

**What to build:** A Board states its format with `format: 0.2` in `stages.md` (see Board format in `CONTEXT.md`). fuda notices boards on another format and points to one Guide page, "Moving a board from 0.1 to 0.2", that an agent can follow.

**Blocked by:** 03, 04, 05, 08

**Status:** ready-for-agent

Why: the merge changes three things existing boards rely on. Everything before the merge is 0.1. All three land together, so there is one 0.2 migration, not one per change.

Rules:

- `format:` is a top-level key in `stages.md`'s front matter, next to `stages:`, as major.minor.
- A `stages.md` without `format:` is 0.1. A board with no `stages.md` uses fuda's defaults and counts as current.
- fuda reads a board of any format as well as it can. An older format shows a notice: "This board uses the 0.1 format. Ask your agent to follow Moving a board from 0.1 to 0.2", with the link. A newer format shows a notice to update fuda, naming both versions (from `/api/version`, ticket 08). Parts fuda cannot read are Problems as usual.
- The migration page lists the three changes in order, then the marker:
  1. Move the board to its own repository and read it from `main`: no `docs/board` on `develop`, no submodule (from ticket 03 and Guide `02-setup-a-repo`).
  2. Owners to person ids (steps from ticket 05).
  3. `when: done` to `status: [done]`, with `days` on that entry (ticket 04).
  4. Then add `format: 0.2` to `stages.md`.
- The Problems for owners as names and for `when: done` link the page. The README links it too.

- [ ] `format:` parsed; an invalid value is a Problem
- [ ] Notice for an older and for a newer format; none for current or for a board without `stages.md`
- [ ] Guide page "Moving a board from 0.1 to 0.2", written so an agent can follow it, linked from the README and the related Problems
- [ ] Guide `04-board-config` documents `format:`, checked against the constant (ticket 08)
- [ ] Tests for each format case
