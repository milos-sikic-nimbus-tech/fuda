Run this once, in a session opened on the `fuda-eliot` checkout, after `.scratch/` was imported there
with `git filter-repo` (history kept). Paste everything below the line.

---

Convert eliot's Matt-style tickets in this repository into a fuda board. This is a one-shot
migration: do it once, carefully, and stop for my approval at the two points marked **Stop**.

## Read first

- `/Users/milos192/projects/fuda/api/internal/guide/pages/03-task-format.md`: the task file format.
- `/Users/milos192/projects/fuda/api/internal/guide/pages/04-board-config.md`: `stages.md`, `labels.md`.
- `/Users/milos192/projects/fuda/api/internal/guide/pages/06-agent-guide.md`: the `fuda-tasks` skill,
  and the "Matt Pocock's skills" section (Board setup, the tracker doc).

## Input

Under `.scratch/` at the repository root (if it is elsewhere, find it and tell me):

1. `.scratch/wayfinder/document-engine/map.md` and `tickets/NN-*.md`. Frontmatter `type`,
   `status: open|closed`, `assignee`, `blocked-by: [NN-slug, …]`; no `id` or `title` (the title is the
   `# ` heading). Resolved tickets hold their answer under `## Resolution`. The map already starts with
   `label: wayfinder:map` frontmatter and links tickets as `tickets/NN-slug.md`.
2. `.scratch/engine-fixes/issues/NN-*.md`. No frontmatter; `# NN: Title` heading, `**Status:**` and
   `**Blocked by:**` lines in the body (`None (can start immediately)`, or numbers with or without
   titles, e.g. `05, 06, 07, 08`, `04 Bronze in markdown mode`).

## Output

```
docs/board/
  stages.md, labels.md        from the guide's Board setup (with a Done Stage after Validated)
  TASKS.md                    the guide's Workflow page
  tasks/document-engine/      map.md + open wayfinder tickets
  tasks/engine-fixes/         engine-fixes tickets
docs/agents/issue-tracker.md  the guide's tracker doc, <PREFIX> = EL
.claude/skills/fuda-tasks/SKILL.md   the guide's fuda-tasks skill
```

Rules:

- **Ids**: one sequence, `EL-N`. document-engine tickets first in their number order (01 → EL-1 …),
  then engine-fixes in theirs. Filename `EL-N-<old slug without the NN- prefix>.md`. The map gets no id.
- **Frontmatter** for every ticket: `id`, `title` (from the heading, without an `NN:` prefix),
  `status`, `blocked_by`, `owner`, `labels`, `added`, `claimed`, plus `done` for resolved ones. Fields
  already there are carried over (`assignee` → `owner`, `blocked-by` → `blocked_by`, `type: X` →
  `labels: ["wayfinder:X"]`); any other existing field is kept as it is. `added` is the date of the
  file's first commit (`git log --diff-filter=A --format=%as -- <old path> | tail -1`). Body-line
  `**Status:**` and `**Blocked by:**` move into frontmatter and leave the body.
- **Heading**: `# EL-N: Title`.
- **blocked_by**: ids only, `"EL-3, EL-9"`; "None…" becomes `""`. Resolve number references within the
  ticket's own folder.
- **Status**:
  - wayfinder `open` with no assignee → `backlog`; `open` with an assignee → `in progress`, `claimed`
    = the date of the commit that set the assignee.
  - wayfinder `closed` with a `## Resolution` → rename that heading to `## Answer`, stays in
    `tasks/document-engine/`, `status: "done"`, `done` = the date of the commit that closed it.
    Archiving is later cleanup, not part of this conversion.
  - wayfinder `closed` that was ruled out of scope or abandoned (no real answer, or listed under the
    map's Out of scope) → stays in `tasks/`, `status: "wontfix"`, reason under `## Answer`. When you cannot tell, ask me.
  - engine-fixes: the `**Status:**` word as is (`ready-for-agent` etc.). A ticket whose boxes are all
    ticked and whose work is in eliot's `main` → ask me before changing it.
- **Links**: rewrite every link and path between tickets and in the map to the new relative paths:
  the map's `tickets/NN-slug.md` links, the Notes line naming `.scratch/engine-fixes/issues/`, and
  links inside ticket bodies. A bare `NN` or old filename mentioned in prose becomes the `EL-N` id.
  Afterwards no `.scratch/` path and no `tickets/` link remains (`rg -n '\.scratch/|tickets/' docs`).

## Steps

1. Build the full plan without writing anything: a table of old path → new path, id, status, and every
   field and link change per file, plus the list of `wontfix` decisions and anything unclear.
   **Stop**: show me the dry run and wait.
2. Commit 1, moves only: `git mv` every file to its new path and change nothing inside. Git stores no
   renames; it guesses them later by comparing content, so a commit that only moves files keeps
   `git log --follow` working for each ticket. Message: `chore(tasks): move eliot tickets to docs/board`.
3. Commit 2, contents: frontmatter, headings, links, statuses, plus `stages.md`, `labels.md`,
   `TASKS.md`, the tracker doc and the `fuda-tasks` skill. Message:
   `chore(tasks): convert eliot tickets to fuda task files`. Remove `.scratch/` if anything is left in it.
4. Check with fuda on this checkout:
   ```bash
   FUDA_HOSTS=local FUDA_LOCAL_PATH="$PWD" FUDA_ADDR=:8099 /Users/milos192/projects/fuda/bin/fuda
   ```
   then `curl -s localhost:8099/api/boards` for the board's path, and its `/board` response's
   `problems` must be empty. Open `http://localhost:8099` and confirm the document-engine epic shows
   its Map and its resolved tickets in the Done column. Fix and amend commit 2 until clean, then stop the server.
5. **Stop**: show me `git log --stat -2` and the problems check. Push to `main` only after I approve.
