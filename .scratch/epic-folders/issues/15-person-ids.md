# 15: Owners are person ids, linked to GitHub logins

**What to build:** Each `people.md` entry gets a required, unique `id` that `owner` and `tester` store, plus `github` (login) and `kind` (`human` or `agent`). Record it as ADR 0007.

**Blocked by:** 13

**Status:** ready-for-human

Why: linking a git or GitHub username to a person needs one stable key. Display names change, and alias folding made matching fuzzy ("Nemanja", "nemanja" and "Nemanja Mudrinic" meant one person), so nothing could safely say which person a login or a PR author is. Agents such as `research-agent` own Tasks too and need to be told apart from humans. Strict ids beat an alias fallback: old name-style owners show as Problems, so boards get cleaned up. Breaking `people.md` is accepted (early stage, two boards).

- **Schema:** `id` (required, unique ignoring case), `name` (shown), `aliases` (search only, never used to match an owner), `github` (optional), `kind` (optional, default `human`). A missing or repeated id, or an unknown kind, makes the file a Problem, and the board derives people from the Tasks.
- **Matching:** with `people.md`, `owner` and `tester` must match an id exactly; any other value is a Problem on that Task, and the card shows it as written. Without `people.md`, nothing changes.
- **Assign** writes ids in the file's own style and refuses an unlisted id (422).
- **Me:** the board response includes `me` (GitHub `/user` login → person, cached per token): "Assign to me", a "My tasks" filter, "(you)" in the Owner filter.
- **Cards:** GitHub avatar, agent badge, open-PR author shown as a person when the login matches. The Owner filter can pick Humans or Agents.
- **Migration:** `fuda-test-board` migrated (commit 83210d7). `fuda-eliot` is migrated separately.

Not built: Azure DevOps login mapping; PR authors show only on the card's PR chip.

PR: https://github.com/milos-sikic-nimbus-tech/fuda/pull/11

- [x] A name-style owner on a board with `people.md` shows under Problems
- [x] "Assign to me" writes the logged-in user's id
- [x] Avatars and agent badges show on cards
- [x] Boards without `people.md` behave as before
- [ ] `fuda-eliot` migrated to ids
