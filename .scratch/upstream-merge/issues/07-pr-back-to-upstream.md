# 07: Send the merged fork back to upstream as one codebase

**What to build:** Open a PR from the fork's `main` to `nemanja-mudrinic-nimbus-tech/fuda` `main`, with upstream's PR template filled in. Once Nemanja merges it, upstream `main` is the one codebase: the fork resets its `main` to it and works through PRs from then on.

**Blocked by:** 03, 04, 05, 06

**Status:** needs-triage

Why: the goal is one codebase, not two that agree. Because 02 merged upstream in, this PR has no conflicts, and its commits group by idea for review.

- [ ] Fork `main` contains `upstream/main` (no new upstream commits left unmerged; merge again if there are)
- [ ] PR opened against upstream `main`, reviewed by Nemanja (CODEOWNERS)
- [ ] After merge: fork `main` equals upstream `main`, and new work starts from upstream `main`
- [ ] Ticket 01's open points are all answered or under TBD in `docs/decisions.md`
