# 04: Build the Done column from status, not from the archive

**What to build:** Keep upstream's `when: done` column and its `days` window, but fill it from tasks whose status is `done` and whose `done` date is in the window, wherever the file is (on the board or archived). The cards stay locked against Move, drop and Assign.

**Blocked by:** 02

**Status:** needs-triage

Why: we agreed Done is a status and the archive is a place (idea 2). Upstream's column reads only archived tasks, so a `done` task still on the board would be missed or shown twice.

Decided: the column is on by default. A board without a `when: done` entry gets a Done column with the 14-day window; `stages.md` can still set `days` or place the column.

- [ ] A board with no `when: done` entry shows the Done column
- [ ] A `done` task on the board and a `done` archived task both show in the column within the window, once each
- [ ] A won't-do archived task does not show
- [ ] `done` on the board is not a Problem
- [ ] Insights and the Archive view still count and show status as in the fork
- [ ] Guide `04-board-config` and `docs/architecture.md` describe the new rule
