# 03: Map-aware epic view

**What to build:** When an epic folder has a `map.md` with frontmatter `label: wayfinder:map`, fuda shows it in a modal for that epic: the map text and the epic's progress.

**Blocked by:** 01, 02

**Status:** done

Why: the map holds an epic's Destination, Notes, Decisions so far, Not yet specified and Out of scope. Seeing them next to the cards shows where the epic is going and how far it is.

- A `map.md` without the label is skipped (ticket 02) and not shown. The label is a frontmatter field, `label: wayfinder:map`; only `tasks/<epic>/map.md` counts, not `tasks/map.md` or `archive/<epic>/map.md`.
- A broken map (bad frontmatter) becomes a Problem and never breaks the board.

UI: a **map button** in the toolbar, right of the active-filter chips and the Clear button (it moves next to Filter while no filter is active).

- There is no separate epic picker: the Filter menu's Epic group already filters by epic (multi-select).
- The button is hidden unless at least one epic on the board has a map. It is enabled only when exactly one epic with a map is selected in the Filter menu.
- It opens a modal that shows the map's body as raw text, without its YAML frontmatter. A nicer map component comes in a later ticket.
- **Progress** sits at the top of the modal: a count of the epic's tasks per board column, as the board shows them (so In review included), plus an Archived count. No single percentage, because fuda has no "done" column and teams order `stages.md` differently. Archived tasks are finished work (archive means shipped, not deleted).

- [x] The API returns the map per epic
- [x] The map button is hidden without maps, enabled only for exactly one selected epic with a map, and opens a modal with the map body (no frontmatter) and per-column counts plus Archived
- [x] Tests for parsing, the label check and the per-column counts
- [x] README/Guide and `docs/architecture.md` updated
