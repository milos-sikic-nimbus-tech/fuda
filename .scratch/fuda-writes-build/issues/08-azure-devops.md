# 08: Azure DevOps login, Boards and writes

**What to build:** A person logs in with a Microsoft Entra ID app (confidential client on the web). fuda lists the user's organizations and their `fuda-` repositories, reads them with the user's token, and writes Move and Assign through the pushes API as compare-and-swap commits. URLs are `/azure/<org>/<project>/<repo>/`. First verify the open research items: whether users can consent to `vso.code_write` without an admin.

**Blocked by:** 04: Assign on GitHub

**Status:** ready-for-agent

- [ ] Login, Board list, Move and Assign work against a real Azure DevOps test repo (checked by hand)
- [ ] Azure source satisfies the same board interfaces as GitHub
- [ ] Guide and README updated

Spec: `.scratch/fuda-writes/spec.md`.
