# 04: GitHub Issues Backend

**What to build:** A person opens a GitHub repository's Issues as a Board: issues are Tasks, a Move and an Assign change the issue on GitHub, and the issues are the Board's only source of truth. Built properly in `internal/backend/githubissues`, from what the spike (ticket 01) learned; no spike code is reused.

**Blocked by:** 02, 03

**Status:** needs-triage

Open questions, each answered by ticket 01 before this is ready:

- **Stages:** where they and their order come from.
- **Discovery and Board id:** which repositories are Issues Boards and how they are listed next to markdown Boards. Proposed id: `/github-issues/<owner>/<repo>/`.
- **Status:** the `status:` label, or something else the spike found better. What a closed issue and an issue with no Status are.
- **Archive:** whether closed issues older than some age are the Archive, or there is none.
- **Size:** what is on the Board when a repository has thousands of issues.
- **Permissions:** what the fuda GitHub App needs, and what a person sees when it lacks them.

Fixed already:

- Login is the existing GitHub Host login; every call uses the caller's token.
- Implements `Backend` from ticket 02 and passes its contract test. No Docs capability.
- Move and Assign refuse when the value changed since the person saw it.
- No mirroring to markdown files and no syncing between Backends.

- [ ] An Issues Board lists, opens and filters like a markdown Board
- [ ] Move and Assign write to GitHub; a stale one goes back with a message
- [ ] A read-only account gets a Read-only Board
- [ ] A failed read keeps the last Board; the Backend passes the contract test
- [ ] Unit tests with recorded API responses in `testdata/`
- [ ] Checked against a real repository in the running app (by hand, screenshots)
- [ ] README, Guide pages, `CONTEXT.md` (Task is no longer "one markdown file"), `docs/architecture.md` and `docs/decisions.md` are updated
- [ ] `make check` passes
