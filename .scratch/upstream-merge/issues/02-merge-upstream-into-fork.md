# 02: Merge upstream main into the fork

**What to build:** One merge of `upstream/main` (`7ec791d`, last change 29 October) into the fork's `main`, on a branch. Resolve every conflict the way ticket 01 agreed, so the tree builds and `make check` passes. Record the five decisions as ADRs in the same branch.

**Blocked by:** 01

**Status:** ready-for-agent

Why: the agreed ideas are all the fork's, so the fork is the base and upstream's features come in on top. One merge keeps upstream's history, so the later PR back to upstream (07) has no conflicts.

A trial merge (`git merge-tree main upstream/main`) shows 17 files in conflict:

- Board core: `api/internal/board/board.go`, `board_test.go`, `sync.go`, `api/internal/config/config.go`, `api/internal/httpapi/api.go`
- Sources: `api/internal/source/azure/azure_test.go`, `api/internal/source/github/github_test.go`
- Client: `TaskSheet.tsx`, `useFilterGroups.tsx`, `lib/api.ts`, `lib/filters.ts` and its test, `lib/insights.test.ts`, `lib/moves.test.ts`
- Docs: `guide/pages/04-board-config.md`, `docs/architecture.md`, `docs/decisions.md`

How to resolve:

- Where an idea is at stake, take the fork's rule: boards read `main` (idea 1), `done` is a valid board status (idea 2), owners are person ids (idea 3), epic folders (idea 4).
- Take upstream's features as they are and leave their joining to the fork's rules to tickets 03 to 06: board list setting (F-10), Done column (F-6), Activity timeline and details (F-5, F-8), No owner filter (F-2), owner line check.
- Keep `.scratch/`, `docs/adr/`, `docs/agents/` and `CONTEXT.md`: undo upstream's deletions of them (idea 5).
- Drop upstream's code for reading two branches: the `InProd` card field, `FUDA_WATCH_MAIN`, `mainStatuses` and the `develop`/`main` split in `sync.go`. The fork reads only `main`; "in prod" comes back later as a label the user sets.
- In review keeps the fork's meaning: an open PR, or a Status set by a Move (`27a0734`). Upstream's glossary line that only a PR puts a card there is dropped.
- Take upstream's `.github/CODEOWNERS` and `.github/pull_request_template.md`.
- Keep `CONTEXT.md` as the glossary for now; upstream's `docs/glossary.md` merges into it (Activity entry comes in; upstream's Done column is replaced by the Done Stage already in `CONTEXT.md`). Moving the glossary is part of the own-board follow-up.

Decision records:

- Amend ADR 0006 (archive is a place) with the Done Stage rules from ticket 04, and ADR 0007 (owners are person ids) with the strict switch from ticket 05. Upstream never accepted them, so they change in place.
- New ADRs from 0008: a board is its own repo read from `main`, with listed boards without the `fuda-` prefix (idea 1); epic folders with the Activity switch (idea 4); how fuda's own work is tracked (idea 5); one version for fuda and the Board format (tickets 08, 09).
- `docs/decisions.md` keeps upstream's content and gets a TBD entry for "in prod" as a label the user sets.

- [ ] `upstream/main` merged into a branch off the fork's `main`, with upstream's commits kept
- [ ] All 17 conflicts resolved per the list above; no conflict markers left
- [ ] `.scratch/`, `docs/adr/`, `docs/agents/` still present
- [ ] No `InProd`, `FUDA_WATCH_MAIN` or `develop` reading left in the code or docs
- [ ] ADRs 0006 and 0007 amended; new ADRs 0008 and up for ideas 1, 4, 5 and versioning
- [ ] `docs/decisions.md` has a TBD entry for "in prod" as a label the user sets
- [ ] `make check` passes
- [ ] The app runs against a fixture board (test-fuda skill): board, archive, owner filter and task panel load
