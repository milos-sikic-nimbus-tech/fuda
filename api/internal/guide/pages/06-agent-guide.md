# Agent guide

Copy the block below into `.claude/skills/fuda-tasks/SKILL.md` (or an `AGENTS.md` section), and
add this pointer to your `CLAUDE.md`:

```markdown
## Tasks
Work is tracked as task files in `docs/board/tasks/` (rules: `docs/board/TASKS.md`; columns,
labels and people: `docs/board/stages.md`, `labels.md`, `people.md`). Before starting work, find
or create its task and claim it on `main` of the tickets repository. After the PR that delivers it
merges, set `status: "merged"` and its PR number there. Follow the `fuda-tasks` skill for every task change.
```

## The skill

```markdown
---
name: fuda-tasks
description: Create, claim and move task files in docs/board/tasks so the fuda board stays right. Use before starting any task, when the PR that delivers it merges, and when retesting, blocking, abandoning or archiving one.
---

# Task files for the fuda board

Tasks live in docs/board/tasks/, one file per task, named <ID>-<slug>.md. The frontmatter is the
task's state. Read docs/board/TASKS.md for the rules and docs/board/labels.md, stages.md and
people.md for the allowed label values, statuses and names. Never invent a status or a label
group that is not there; if one is really new, add it to the config file in the same commit.

Before writing code for a task:
1. Pull `main` of the tickets repository and open the task file.
2. If it has an owner, stop and ask the human.
3. Otherwise set owner, status "in progress" and claimed (today), commit only that file as
   "chore(tasks): claim <ID>", show the commit to the human and push to `main` only after they
   approve.

New work found while working:
- Create a task from the template with the next unused id, status "backlog", added today and a
  type label. Show it before pushing.

After the pull request that delivers the task merges:
- In the tickets repository, set status "merged", add the PR number to pr, and write what was checked under "## Evidence".

Retest, blocked, abandon, archive:
- testing (with tester), then validated; on failure back to "in progress" with what failed.
- Set or clear blocked_by.
- Abandon: status "backlog", clear owner and claimed, explain why under "## Evidence".
- Never write status "done" except when moving a file to docs/board/archive/.

Commits that change task files touch only docs/board/ and use
"chore(tasks): <add|claim|testing|validated|blocked|abandon|archive> <ID>".

Keep files valid: an invalid task file disappears from the board into Problems.
```
