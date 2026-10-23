# 04: Merged state from code-repo pull requests

**What to build:** A task can get its merged state from a pull request in a code repository, not only from its own file. Two ways link a PR to a task, and both are supported:

1. **Id matching** (exists for In review): the task id is in the PR's title or branch name in a repo from `repos.md`. Extend it to merged PRs.
2. **Explicit links:** `pr: ["nimbus-tech-io/eliot#123"]` in the task's frontmatter. fuda reads that PR's state.

**Blocked by:** 06

**Status:** needs-triage

Why: when tickets live in a separate repo, the task file can't change in the same PR as the code, so "merged" never shows. In review already works across repos through `repos.md`.

- Plain numbers in `pr:` keep their current meaning (a PR in the board repo).
- Open PR → In review (as today). Merged PR → merged. Closed unmerged → no effect.
- Whether a merged code PR only changes the shown column or also commits `status` to the task file follows the answer in ticket 06.
- "In prod": decide whether a merged code PR whose merge commit reaches the code repo's `main` counts. Write it down even if it is TBD.
- Uses the logged-in person's token. If they can't see the repo, the board still loads without it.
- Watch API cost: merged PRs are many more than open ones. Limit by time or only query PRs a task names.

- [ ] Merged PRs found by id matching mark tasks merged
- [ ] `owner/repo#N` in `pr:` is parsed and its state used; bad values become Problems
- [ ] Tests for both paths and for missing access
- [ ] Guide `04-board-config.md`, `05-workflow.md`, `03-task-format.md` and `docs/architecture.md` updated
