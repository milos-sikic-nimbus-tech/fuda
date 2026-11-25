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

## Agreed (25 November 2026)

Nemanja agreed to the fork's side on all five ideas, with two changes:

1. Board repo read from `main`, no submodule. **Change:** boards without the `fuda-` prefix must keep working, because Nemanja's own setup needs it. Upstream's F-10 setting (`FUDA_GITHUB_BOARDS`, `FUDA_AZURE_BOARDS`) stays.
2. Done is a status, and the archive is a place.
3. Owners are person ids from `people.md`.
4. Epic folders with maps, plus upstream's Activity timeline.
5. `.scratch/` tickets and ADRs stay for now. **Change:** a follow-up moves `.scratch/` into a fuda board folder so fuda shows its own tickets (`.scratch/own-board/issues/01-scratch-to-board.md`).

Follow-up answers (25 November 2026):

- "In prod" comes back later, only as a label the user sets. Deferred: a TBD entry in `docs/decisions.md`, no ticket yet.
- The Done column is on by default (04).
- Owners switch to person ids strictly now. The README tells agents how to reconcile `people.md` and task owners (05).

- Activity logging is a board-wide switch in the board repo, off by default, shown to users in the task panel and the Guide (06).

The merge plan is tickets 02 to 07 in this folder.

- [x] Upstream maintainer has seen the review doc
- [ ] Each of the five ideas has a recorded decision in `docs/decisions.md` (done in 02)
- [x] The merge plan is split into tickets, with the smaller one-side changes placed with the idea they depend on
- [x] Open questions for each idea are answered or moved to TBD
