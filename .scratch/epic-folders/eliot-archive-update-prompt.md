Run this once, in a session opened on the `fuda-eliot` checkout, which is already converted (see
`eliot-conversion-prompt.md`). Paste everything below the line.

---

Bring this board up to fuda's rule that the Archive is a place, not a status (ADR 0006). Resolved
tickets stay in `tasks/` as `done`. Stop for my approval before pushing to `main`.

## Read first

- `/Users/milos192/projects/fuda/api/internal/guide/pages/06-agent-guide.md`: the `fuda-tasks` skill,
  Board setup, and the tracker doc.
- `/Users/milos192/projects/fuda/api/internal/guide/pages/05-workflow.md`: the workflow (TASKS.md).
- `/Users/milos192/projects/fuda/docs/adr/0006-archive-is-a-place-not-a-status.md`.

## Steps

1. Copy the new tracker doc to `docs/agents/issue-tracker.md` (`<PREFIX>` = EL), the `fuda-tasks`
   skill to `.claude/skills/fuda-tasks/SKILL.md`, and the workflow page to `docs/board/TASKS.md`.
   Add a Done Stage after Validated to `docs/board/stages.md`:
   ```yaml
   - name: Done
     status: [done]
   ```
2. Commit 1, moves only: `git mv` the 13 resolved Tasks from `docs/board/archive/document-engine/`
   to `docs/board/tasks/document-engine/` and change nothing inside. Message:
   `chore(tasks): move resolved document-engine tickets back to tasks`.
3. Commit 2, contents: step 1's files, then fix every link to the moved Tasks, the map's Decisions so
   far included (`rg -n 'archive/document-engine' docs`; none may remain). Message:
   `chore(tasks): done is a status on the board, not an archive`.
4. Check with fuda on this checkout:
   ```bash
   FUDA_HOSTS=local FUDA_LOCAL_PATH="$PWD" FUDA_ADDR=:8099 /Users/milos192/projects/fuda/bin/fuda
   ```
   then `curl -s localhost:8099/api/boards` for the board's path; its `/board` response's `problems`
   must be empty, and the 13 Tasks must show in the Done column. Stop the server.
5. **Stop**: show me `git log --stat -2` and the check. Push to `main` only after I approve.
