# 01: File tree of the repo's Docs in the Docs view

**What to build:** The Docs view shows a tree of every Doc, so a reader can browse them instead of only reaching a doc through a link.

**Blocked by:** none

**Status:** ready-for-human

Today the Docs tab opens `board/README.md` (`client/src/components/organisms/AppHeader.tsx`, `CommandPalette.tsx`) and `GET …/docs?path=` returns one doc. There is no list, so other docs are only reachable through links in tasks or docs.

A **Doc** (see `CONTEXT.md`) is any `.md` under `docs/` outside `docs/board/tasks/` (epic folders included) and `docs/board/archive/`. Board-config files (`stages.md`, `labels.md`, …) are Docs and show in the tree.

- API: a list of Doc paths from the snapshot (`internal/board`), no titles. Handler stays thin.
- Client: a fixed pane on the left of `routes/docs/$.tsx`. Entries show the file name; folders first, then A–Z, README not special-cased. Folders collapsible, the open doc highlighted. Clicking a doc navigates to `/docs/<path>`.
- With no `board/README.md`, the Docs tab shows the tree and a "Pick a doc" empty page.
- Update the `Doc` entry's "In code" line in `CONTEXT.md` to name where the list lives.

- [x] `board` returns the Doc list, with unit tests that Tasks, epic files and archived Tasks are left out and config files are kept
- [x] Docs view shows the tree and opens the clicked doc
- [x] The open doc is highlighted and its folders are expanded on a direct link
- [x] Without `board/README.md`, the Docs tab shows the empty page
- [x] `docs/architecture.md` (API table) and the Guide page for Docs are updated
- [ ] Checked in the running app in light and dark mode (by hand)
