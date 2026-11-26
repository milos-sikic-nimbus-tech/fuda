# 06: Activity timeline works with epic folders and the agent guide

**What to build:** Upstream's `## Activity` timeline (task panel, latest entry on cards, detail lines) works for tasks inside epic folders. Merge upstream's agent instruction ("one Activity line per finished step") into the fork's agent guide and `fuda-tasks` skill.

**Blocked by:** 02

**Status:** needs-triage

Why: we agreed on both epic folders and the Activity timeline (idea 4). They touch the same task parsing and the same agent guide.

Decided (25 November 2026): Activity logging is a board-wide switch, **off** by default. It lives in the board repo (a board-level key in `stages.md`), so fuda and every agent read the same value and turning it on is one commit. It is separate from the per-column `activity: true`, which only decides whether cards show the latest entry. No toggle in the UI. Record this in the idea 4 ADR (ticket 02).

How users find it:

1. The task panel's Activity area, when logging is off or a task has no entries, says agents log each step there when Activity is on for the board, and links to the Guide.
2. The agent rule block in the Guide starts its Activity paragraph with "If this board has Activity on…", so agents follow the board's setting and nobody edits the rule.
3. Guide `04-board-config` documents the key with an example, next to `when: done` and `days`, and explains how it differs from the per-column `activity: true`.

- [ ] The board-level switch is parsed, off when missing; an invalid value is listed under Problems and treated as off
- [ ] The API tells the client whether Activity is on for the board
- [ ] The task panel shows the empty-state text with a Guide link when Activity is off or a task has no entries
- [ ] Agent rule block and `fuda-tasks` skill log Activity only when the board has it on
- [ ] A task in `tasks/<epic>/` shows its Activity timeline and its latest entry on the card
- [ ] `map.md` files are still skipped and never parsed for Activity
- [ ] Guide `06-agent-guide` has one agent rule set, covering epic folders and Activity
- [ ] Guide `04-board-config` documents the switch and its default
- [ ] Tests for Activity on a task in an epic folder
