---
status: accepted
---

# Archive is a place, not a status

The Archive is only a place for Tasks that are off the board. It has no tie to Statuses: `status: "done"` is valid in `tasks/`, an archived Task keeps whatever Status it has, and done is a step on the board after Validated, with the `done` date. Archiving is later cleanup, a file move that changes nothing in the file.

## What the old rule was

A Task with `status: done` outside `archive/` was a Problem, so "done" and "archived" meant one thing. No commit, spec or doc records a reason; the rule came with the first parser commit.

## Problems it caused

- A resolved wayfinder Task had to leave the board at once, and the map's Decisions-so-far links pointed into the Archive.
- A won't-do Task could not be archived, or it would be counted as finished work. Either it cluttered the board or its Status lied.
- The board had no place to show finished work before it was cleaned up.

## Considered Options

- Keep the rule: rejected. It forces the workarounds above.
- Rename `done` to `archived` and add `done` later: rejected. Every board's archived files would be rewritten, and the coupling between a Status and a folder would stay.

## Consequences

- The default Stages gain a Done Stage after Validated. A board with its own `stages.md` adds it itself.
- Insights counts Tasks with Status `done` and a `done` date in the period, on the board or archived. Archived `wontfix` Tasks are not counted.
- The Archive view shows each Task's Status, so done and won't-do are told apart.
- Moving a Task to Done in fuda writes the Status line, plus `claimed` on a first move out of the first Stage (ADR 0001); the `done` date is set in git.

## One Done Stage (agreed with upstream, 25 November 2026)

Upstream had built a Done column from archived Tasks (`when: done` with `days`). It and the fork's Done Stage become one:

- The Done Stage is the Stage that lists the Status `done`. It shows every Done Task still on the board, and archived Done Tasks whose `done` date is within the last `days` days (default 14), newest first. If the Stage lists more Statuses, only archived `done` Tasks come back through the window.
- `days` is a key on the Done Stage entry; anywhere else it is a Problem.
- Only archived cards in it are locked. A Done Task on the board is a normal card and can be moved back.
- `when: done` is dropped. A `stages.md` that still has it is invalid under the usual rule, and the Problem names the fix and links the migration page (ADR 0011).
- A board whose own `stages.md` has no Done Stage shows none, with a quiet hint at the end of the board to ask an agent to add it.

Rejected: two columns (the Done Stage plus an archive-only Done column), because done work would show twice; a window on every Done Task, because a Done Task still in `tasks/` would vanish from the board while its file stays; keeping `when: done` as a second way to write the same Stage.
