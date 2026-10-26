# 13: Archive is a place, not a status

**What to build:** The Archive becomes only a place for Tasks that are off the board, with no tie to Statuses. `status: "done"` is valid in `tasks/`, an archived Task can have any Status, and done is a workflow step on the board. Record it as ADR 0006.

**Blocked by:** 08

**Status:** ready-for-agent

Why: today `status: done` outside `archive/` is a Problem (`api/internal/taskfiles/task.go:74`), so done and archived mean one thing. No commit, spec or doc records a reason; the rule came with the first parser commit (`3ce19af`). It forces the Matt tracker doc (ticket 08) into workarounds: a resolved wayfinder Task must leave the board at once, and a won't-do Task must never be archived.

- **Rule removed:** no Problem for `done` outside the Archive; an archived Task keeps whatever Status it has.
- **Default Stages** (no `stages.md`): a Done Stage with `[done]` after Validated (`api/internal/board/columns.go`). A board with its own `stages.md` adds it itself.
- **Workflow:** validated → done is a step on the board, with `done` = the date. Archiving is later cleanup: a file move that changes nothing in the file, done or won't-do, fixing links to the moved file in the same commit.
- **Move to Done** writes only the Status line (ADR 0001); the `done` date is set in git.
- **Insights:** the "Archived as done" tile becomes "Done": Tasks with Status `done` and a `done` date in the period, on the board or archived (`client/src/routes/insights.tsx:49`).
- **Archive view:** a Status badge on each Task, so done and won't-do are told apart.
- **Matt tracker doc** (Agent guide): Resolve = `## Answer`, `status: "done"`, `done` today, the Task stays in `tasks/<epic>/`, the map's Decisions-so-far line links it there. Blocking: a blocker is closed when its Status is `done` or `wontfix`, wherever it is. Closing: won't-do Tasks may be archived; drop the "fuda counts archived as done" reason. Archiving fixes every link to the moved file.
- **`fuda-tasks` skill** (Agent guide): drop "Never write status done except when moving to archive"; add done as a step; archive keeps the Status.
- **ADR 0006** "Archive is a place, not a status": the old rule (done only in the Archive, no recorded reason), the problems it caused (ticket 08's workarounds, won't-do Tasks either cluttering the board or counted as finished), the options (keep it; rename `done` to `archived` and add `done` later, which needs every board's archived files rewritten and keeps the coupling), and the choice.
- **Docs:** Guide Task format (status table, "When a file is wrong"), Workflow (table and diagram), Overview if it mentions the rule, `docs/architecture.md` (Validation, sequence diagram, Epic maps), `CONTEXT.md` (Archive, and a Done entry once built), README if it mentions the rule.
- **eliot prompts:** update `.scratch/epic-folders/eliot-conversion-prompt.md` (resolved Tasks stay in `tasks/` as `done`) for future projects, and write `.scratch/epic-folders/eliot-archive-update-prompt.md` for `fuda-eliot`, which is already converted:
  - copy the new tracker doc, `fuda-tasks` skill and `TASKS.md`; add a Done Stage after Validated to its `stages.md`;
  - `git mv` the 13 resolved Tasks from `archive/document-engine/` to `tasks/document-engine/` in a moves-only commit, then fix every link to them, the map's Decisions so far included;
  - check with fuda for no Problems, and stop for approval before pushing to `main`.

Not here: Insights' open-work count (`isOpen`, a hardcoded list where `wontfix` counts as open) is ticket 14.

- [ ] A `done` Task in `tasks/` loads with no Problem and shows in the default Done Stage
- [ ] An archived Task with any Status loads, and the Archive view shows its Status
- [ ] The Done tile counts `done` Tasks by `done` date in both places, and no archived `wontfix`
- [ ] ADR 0006 written; no doc says `done` belongs only in the Archive
- [ ] Both eliot prompts written; `fuda-eliot` shows the 13 resolved Tasks in its Done column with no Problems
