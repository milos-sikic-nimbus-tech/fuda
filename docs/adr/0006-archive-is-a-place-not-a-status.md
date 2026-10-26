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
- Moving a Task to Done in fuda writes only the Status line (ADR 0001); the `done` date is set in git.
