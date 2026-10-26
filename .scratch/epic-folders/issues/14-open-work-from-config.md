# 14: Open work follows the board's Statuses

**What to build:** Insights decides which Tasks are open from the board's configuration instead of a hardcoded list.

**Blocked by:** 13

**Status:** needs-triage

Why: `isOpen` (`client/src/lib/insights.ts:114`) treats every Status except `merged`, `testing`, `validated` and `done` as open. Boards that use other words count closed work as open: a `wontfix` Task from the Matt tracker doc (ticket 08) shows as open work.

Questions to settle:

- Where a board marks its closed Statuses: a flag on a Stage in `stages.md`, a separate list, or a convention.
- Whether the default stays the current list when `stages.md` says nothing.

- [ ] A `wontfix` Task on a board that marks it closed is not counted as open
- [ ] Boards without the new setting count open work as today
