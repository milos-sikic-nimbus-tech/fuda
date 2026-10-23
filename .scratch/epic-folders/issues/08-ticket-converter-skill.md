# 08: Converter skill for Matt-style tickets

**What to build:** A skill that converts Matt-style (wayfinder / `.scratch`) tickets into fuda task files, so a repo like eliot can move its tickets into a fuda tickets repo (`fuda-eliot`). This ticket lives here for now and can move out of fuda later.

**Blocked by:** 01, 02

**Status:** needs-triage

Why: eliot's tickets are in two Matt-style formats that fuda can't read. fuda stays strict about one format, so the conversion happens once, in a skill, not in fuda.

Input, in eliot (`/Users/milos192/projects/eliot`):

- `.scratch/wayfinder/<epic>/map.md` plus `tickets/NN-*.md`: frontmatter `type`, `status: open|closed`, `assignee`, `blocked-by`; no `id` or `title`.
- `.scratch/<feature>/issues/NN-*.md`: no frontmatter; `**Status:** ready-for-agent` and `**Blocked by:**` lines in the body.

Output, in `fuda-eliot`:

- `docs/board/tasks/<epic>/EL-N-<slug>.md`, one level of epic folders (ticket 01). The folder is the epic; no `epic:` or `theme:` labels (deprecated).
- One id sequence for the whole repo, `EL-N`, no per-epic prefixes. Renumber every ticket and fix map links and cross-references (`blocked-by`, links between tickets).
- `map.md` is kept per epic, optional, in `tasks/<epic>/map.md` (tickets 02, 03). fuda only reads it as a map if it starts with a YAML frontmatter block holding `label: wayfinder:map`:

  ```markdown
  ---
  label: wayfinder:map
  ---
  # Map: ...
  ```

  Wayfinder maps carry the label as a body line (`Label: wayfinder:map`) and no frontmatter, so fuda ignores them as they are. The skill adds the frontmatter block at the top and removes the body line. A map with broken frontmatter becomes a Problem.
- Missing fuda fields (`id`, `title`, `added`, …) are added; existing fields are kept. Body-line fields (`**Status:**`, `**Blocked by:**`) become frontmatter.
- Statuses keep the Matt-style words, listed in `stages.md` (a column can hold several statuses).
- **Closed tickets:** `closed` covers both done and won't-do. Done goes to `archive/` with `status: done`. Won't-do or abandoned goes back to backlog (or is dropped) and is never archived, because fuda counts archived tasks as finished work in epic progress (ticket 03).

Questions to settle while building:

- How to tell done from won't-do in a closed ticket (an `## Answer`, a `wontfix` label, asking the user).
- The status words and how they map to `stages.md` columns.
- Whether the skill is re-runnable or one-shot.

- [ ] The skill converts both input formats; a dry run lists every rename and field change before writing
- [ ] Ids, map links and cross-references are consistent after the renumbering
- [ ] Every converted `map.md` starts with a `label: wayfinder:map` frontmatter block (no body-line label left), and fuda shows it as a map
- [ ] Won't-do tickets are never archived
- [ ] fuda loads the converted repo with no Problems except intended ones
