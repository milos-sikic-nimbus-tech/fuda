# 12: Code repos keep no copy of Tasks

**What to build:** Drop the `docs/board` submodule from the docs. A code repo keeps no copy of its board's Tasks; agents there only name the task id. Record it as ADR 0005.

**Blocked by:** 09

**Status:** ready-for-agent

Why: the submodule pointer shows as modified in the code repo, and new clones and worktrees have an empty `docs/board` until someone runs `git submodule update --init`. A code-repo agent does not need the Tasks: In review comes from the task id in a PR's title or branch, and merged will come from code PRs (tickets 04, 06).

- **No copy in code repos.** No submodule and no clone inside the code repo. The board stays its own `fuda-<project>` repo on `main` (ADR 0002). A developer's own checkout of the board, anywhere on disk, is the intended way to work on tickets: its path is the `fuda board:` line in their gitignored `CLAUDE.local.md`, and agents use it through the board's tracker doc (ticket 08).
- **In a code repo**, an agent puts the task id in the branch name and PR title, and nothing else.
- **Claim, create and move** happen in fuda (Owner, Move to), or by an agent in the developer's board checkout, from a session in either repo.
- **Merged:** until 04/06 land, a person moves the card to Merged in fuda after the code PR merges. The Workflow page's Deliver row says so.
- **The `fuda-tasks` skill** lives in `.claude/skills/fuda-tasks/` of the `fuda-` repo. The Agent guide page shows it there.
- **Code-repo rule:** the Agent guide page gives a short block for the code repo's `CLAUDE.md`/`AGENTS.md`: Tasks live in `fuda-<project>`; your checkout's path is the `fuda board:` line in `CLAUDE.local.md`; change tasks only there, following `<board>/docs/agents/issue-tracker.md`; put the task id in the branch and PR title. It replaces the submodule rule in `CONTRIBUTING.md`, and matches the pointer in the Agent guide's Matt Pocock's skills section.
- **Migration** (Guide `02-setup-a-repo.md` and `CONTRIBUTING.md`): `git subtree split --prefix=docs/board`, push as `main` of `fuda-<project>`, `git rm -r docs/board` in the code repo, add the rule. One line: a code repo that already has the submodule removes it with `git rm docs/board` and deletes its `.gitmodules` entry. "Keep the submodule current" goes.
- **ADR 0005** "Code repositories keep no copy of Tasks", with the submodule as a considered and dropped option and why. ADR 0002 gets a line: its submodule part is superseded by 0005.
- `docs/decisions.md` (Boards proposal, migration bullet), the README link to the migration section, and the Workflow and Agent guide pages match. No `submodule` mentions remain outside ADR 0002, 0005 and the one removal line.

- [ ] ADR 0005 written; ADR 0002 points to it
- [ ] Guide setup, workflow and agent guide pages, `CONTRIBUTING.md`, README and `docs/decisions.md` describe no copy in code repos
- [ ] The migration has no submodule step, plus the one removal line
