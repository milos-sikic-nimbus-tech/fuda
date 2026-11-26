---
status: accepted
---

# Epics are folders; Activity logging is a board switch, off by default

Two features meet on the same task files. Epics are folders (`tasks/<epic>/`), with an optional `map.md`, and planning skills use the board as their tracker. A Task's `## Activity` section is a log of finished steps, one line each (time, actor, event, result, optional usage), shown as a timeline in the task panel and, for a Stage with `activity: true`, as the latest entry on its cards. Agreed with the upstream maintainer on 25 November 2026.

Whether agents write Activity lines is a board-wide switch in the board repository (a key in `stages.md`), off by default. fuda and every agent read the same value, and turning it on is one commit. It is separate from the per-Stage `activity: true`, which only decides what cards show.

## Why

- Upstream's agent pipeline logs every step; boards without that pipeline would get noise or an empty timeline if every agent had to log.
- Keeping the switch in the board repository, not in the UI, fits a board whose repository is its only source of truth, and fuda writes only Moves and Assigns.

## Considered Options

- Activity as a rule for every agent: rejected, it costs every board a line per step whether it wants the log or not.
- Activity only for upstream's pipeline, with no switch: rejected, other boards could not opt in.
- A toggle in the UI: rejected, it would be a new kind of write.

## Consequences

- When the switch is off or a Task has no entries, the task panel says agents log steps there when Activity is on, with a link to the Guide.
- The agent rule block in the Guide logs Activity only when the board has it on.
- Old `epic:` labels are on their way out (TBD in `decisions.md`).
