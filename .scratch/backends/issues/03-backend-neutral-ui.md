# 03: The payload and the UI do not assume git

**What to build:** The Board payload and the client work for any Backend. A markdown Board looks the same as today; a Board whose Backend has no Docs, Archive or In review simply does not show them.

**Blocked by:** 02

**Status:** needs-triage

Update from ticket 01's "UI" answer, then set `ready-for-agent`.

Known git leaks:

- **Capabilities.** The payload lists what the Backend has. The client hides the Docs tab, the Archive view, the In review Stage, the Problems drawer and epic maps when absent (`AppHeader.tsx`, `CommandPalette.tsx`, `ProblemsDrawer.tsx`).
- **Freshness.** `sync.branch.{branch,sha}` becomes `syncedAt`, `lastError` and an optional short revision label and source name that the Backend fills in (markdown: branch and short SHA, or "working tree"). `SyncIndicator.tsx`, `WorkspacePicker.tsx` and `routes/__root.tsx` show what is present.
- **Errors.** The `no-main-branch` case in `routes/__root.tsx` becomes a Backend-supplied message, not a code the client knows.
- **Board path.** The server sends each Board's path and login URL. `lib/boardPath.ts` stops counting segments per host and guessing the login from the path.
- **People.** A person carries an avatar URL; `person.github` leaves the client (`lib/people.ts`, `OwnerPicker.tsx`, `lib/api.ts`).
- **Wording.** "repo" and "repository" texts become neutral or come from the Backend: `routes/index.tsx` (read-only note), `routes/docs/$.tsx`, `ProblemsDrawer.tsx`, `WorkspacePicker.tsx` ("Open repo"), `BoardList.tsx` (the `fuda-` hint), the "Reading the repo…" toast in `lib/queries.ts`.
- **Origin.** `Origin.repo` and `lib/hosts.ts` grouping: Boards group by Host and show the Backend where one Host has two.
- **Stays as is.** The GitHub App install banner belongs to the GitHub Host. Pull request chips on cards stay; they show only with the In review capability.

- [ ] Payload carries capabilities, generic freshness, the Board path and avatar URLs; Go tests cover the markdown values
- [ ] Client DOM tests: a Board without Docs, Archive and In review shows none of them and no git wording
- [ ] A markdown Board looks and behaves as before (checked by hand, light and dark, screenshots)
- [ ] No `github`, `azure`, `branch` or `sha` in client code outside Host login, the Host icon and the install banner
- [ ] `docs/architecture.md` (API table) and the Guide pages are updated
- [ ] `make check` passes
