# 01: Agree the merge path between upstream and the fork

**What to build:** A decision, no code. Agree with the upstream maintainer (Nemanja) on five ideas that stop the fork (`milos-sikic-nimbus-tech/fuda`) and upstream (`nemanja-mudrinic-nimbus-tech/fuda`) from merging. Record each agreed idea in `docs/decisions.md`, then split the merge into reviewed PRs, one per idea.

The comparison and the options for each idea are in the review doc: https://claude.ai/artifact/4K8tYgq7UqphxN98Q3wMCc

**Blocked by:** none

**Status:** needs-triage

Why: both copies moved about 30 commits apart since 13 October 2026, and they now disagree on what a board is, not only on features. Each week the two drift further, so the merge gets harder. The code is easy to join once the ideas are agreed.

The five ideas:

1. **A board is its own repository, read from `main`.** Upstream reads `docs/board` on `develop` from the code repo, with an "in prod" badge from `main`. The fork uses a `fuda-<project>` repo whose `main` is the board, with no submodule in the code repo. Open: whether "in prod" comes back from merged code PRs, and whether upstream's "any repo name" setting stays.
2. **Done is a status; the archive is a place.** Upstream: `done` means archived, with an optional Done column for the last 14 days. The fork: Done is a step after Validated, and archiving keeps the status. Open: whether the Done column is on by default.
3. **Owners are person ids.** Upstream stores names and matches them loosely. The fork stores ids from `people.md`, linked to GitHub logins and marked human or agent. Open: strict switch now, or a grace period that accepts names with a warning.
4. **Activity log and epic folders.** Upstream: agents log each step in `## Activity`, shown as a timeline. The fork: epics are folders, with optional maps. Open: whether epic folders are the direction, and whether Activity logging is a rule for everyone or only for the pipeline.
5. **How fuda's own work is tracked.** Upstream removed `.scratch/` tickets, ADRs and retros after PR review, and keeps decisions in `docs/decisions.md`. The fork kept them. Open: whether to adopt upstream's rules and move the tickets onto a fuda board.

- [ ] Upstream maintainer has seen the review doc
- [ ] Each of the five ideas has a recorded decision in `docs/decisions.md` (or a TBD entry with a proposal)
- [ ] The merge is split into PRs, one per agreed idea, with the smaller one-side changes placed with the idea they depend on
- [ ] Open questions for each idea are answered or moved to TBD
