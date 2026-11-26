# 09: Docs for separate tickets repos

**What to build:** The README, Guide pages and `docs/` describe a board as its own tickets repo on `main`.

**Blocked by:** 05

**Status:** resolved

Why: ticket 05 changes how boards work, and the docs still describe tasks on `develop` in the code repo. Land this right after 05, since `docs/` must state only what is true in the code.

- One tickets repo per board, named `fuda-<project>` by convention, with its tasks on `main`.
- Remove same-PR delivery and the develop/main In prod from the workflow (`05-workflow.md`), setup (`02-setup-a-repo.md`), overview, run-locally and self-host pages. `FUDA_WATCH_MAIN` goes from `.env.example` and the self-host table if 05 removes it.
- `main` must allow direct pushes from fuda; say what happens when it is protected.
- In review still comes from code repos listed in `repos.md`.
- A short migration section: moving an old `develop` board into a `fuda-` repo.
- `docs/architecture.md` and `docs/decisions.md` match.

- [x] Every Guide page, the README, `.env.example`, `CONTRIBUTING.md` and `docs/` match the code after 05
- [x] The migration section exists
