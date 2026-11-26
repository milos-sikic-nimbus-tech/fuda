# 04: One Done Stage, built from status

**What to build:** Join the fork's Done Stage and upstream's Done column into one Done Stage (see `CONTEXT.md`). It shows every task with status `done` still on the board, plus archived tasks with status `done` whose `done` date is within the last `days` days (default 14), newest first. Only archived cards in it are locked against Move, drop and Assign.

**Blocked by:** 02

**Status:** ready-for-agent

Why: we agreed Done is a status and the archive is a place (idea 2). Upstream's column read only archived tasks, so with both, a board would show done work in two places.

Rules:

- In `stages.md` the Done Stage is the Stage that lists the status `done`. If it lists more statuses (for example `[done, won't do]`), only archived `done` tasks come back through the window.
- `days` is a key on the Done Stage entry. `days` on any other Stage, or a value that is not a positive whole number, is a Problem; the window then uses 14.
- `when: done` is gone. A `stages.md` that still has it is invalid under the existing rule: fuda uses the default Stages, and the Problem says to replace it with `status: [done]` and links the migration page (ticket 09).
- The default Stages keep their Done Stage. A board whose own `stages.md` has no Done Stage shows none, and shows a quiet line at the end of the board: "Missing the Done column? Ask your agent to add it to stages.md", linking the Guide. It is not a Problem.
- `done` on the board is not a Problem. A Done task on the board is a normal card and can be moved back.

- [ ] Done tasks on the board always show; archived `done` tasks show within the window, once each
- [ ] An archived won't-do task does not come back, even when the Done Stage lists it
- [ ] Archived cards in the Done Stage are locked; Done cards on the board are not
- [ ] `days` on the Done Stage sets the window; elsewhere or invalid, a Problem
- [ ] `when: done` gives the Problem with the fix and link
- [ ] The "Missing the Done column?" line shows only when the board's own `stages.md` has no Done Stage
- [ ] Insights and the Archive view still count and show status as in the fork
- [ ] Guide `04-board-config`, `docs/architecture.md` and ADR 0006 describe the rules
