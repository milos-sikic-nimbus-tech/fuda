# 06: Can fuda commit the status when a code PR merges?

**What to build:** An investigation, no code. Decide whether fuda can write a commit to the tickets repo (`status: merged`, the PR link) when a linked pull request in a code repo is merged, or whether the merged state should stay derived, shown only on the board like In review.

**Blocked by:** none

**Status:** needs-triage

Why: ticket 04 shows merged state from code-repo PRs. Committing it keeps the task file in step with the board; deriving it keeps the repo as the only record without extra commits. Ticket 04 waits for this answer.

Questions to answer:

- **Who commits, and when?** fuda writes with the logged-in person's token, only while someone uses the app. A merge with no one online: does the commit happen at the next sync by whoever opens the board? Is that acceptable, and whose name goes on it?
- **Triggers.** Polling at sync versus webhooks. fuda has no webhook endpoint and the desktop app can't receive one. What would each cost?
- **Duplicates and races.** Two people open the board after the merge: how do we avoid two commits, or a conflict with a manual edit? See `fuda-writes` tickets 07 (concurrent edits) and 08 (who writes the commit).
- **Access.** The person needs write access to the tickets repo and read access to the code repo. What happens when they have only one?
- **Alternatives.** A GitHub Action in the code repo that commits to the tickets repo; or derived only. Compare.
- **Local boards.** Does this apply to a local source at all?

- [ ] Findings and a recommendation are written under `## Answer` in this file
- [ ] The decision is recorded in `docs/decisions.md` (or as TBD with a proposal)
- [ ] Ticket 04 is updated to match
