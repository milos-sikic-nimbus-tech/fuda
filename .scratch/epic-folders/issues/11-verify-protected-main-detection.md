# 11: Verify protected-main detection against real hosts

**What to build:** Confirm that a move refused by a protected `main` shows "main is protected; allow fuda to push to it directly" on GitHub and on Azure DevOps, and fix the detection where the real response differs.

**Blocked by:** 05

**Status:** ready-for-human

Why: ticket 05 detects the refusal from guessed response shapes, tested only against faked responses.

- GitHub (`api/internal/source/github/write.go`, `protectedBranch`): a 403, 409 or 422 whose body contains "protected branch", "rule violations" or "through a pull request".
- Azure DevOps (`api/internal/source/azure/write.go`): a 403 push whose body contains `TF402455`.
- Cases to try on GitHub: classic branch protection, a ruleset requiring a pull request, and a ruleset with a bypass for the fuda GitHub App. A stale-file 409 must still come back as "changed since read", not as protected.

- [ ] A move on a protected `main` shows the protected message on GitHub (classic protection and a ruleset)
- [ ] The same on Azure DevOps with a branch policy
- [ ] The detection and its tests match the real response bodies
