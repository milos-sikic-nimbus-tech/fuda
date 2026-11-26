# 01: Move fuda's own tickets from .scratch onto a fuda board folder

**What to build:** Turn `.scratch/<feature>/issues/*.md` into task files in a board folder that fuda can open, so fuda shows its own tickets. Features become epic folders, `Status:` lines become statuses, `Blocked by:` becomes the board's dependency field. Update `docs/agents/issue-tracker.md` and `CLAUDE.md` to point at the new folder.

**Blocked by:** upstream-merge 07

**Status:** needs-triage

Why: Nemanja agreed to keep `.scratch/` for now, on condition that it moves onto a fuda board later. Reading fuda's tickets in fuda also tests the product on itself.

Open: where the folder lives. A folder in this repo read with the local source is the simplest; a `fuda-fuda` repo follows idea 1 (a board is its own repo). Decide before starting and record it in `docs/decisions.md`. The glossary move (`CONTEXT.md` to `docs/glossary.md`) and folding `docs/adr/` into `decisions.md` can be decided at the same time.

- [ ] Every ticket in `.scratch/` exists as a task, with status, epic and blockers kept
- [ ] fuda opens the board and shows the tickets with no Problems
- [ ] `.scratch/` removed, and agent docs point at the board
- [ ] `make check` passes
