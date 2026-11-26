# 01: Throwaway spike: a Board from GitHub Issues, with Move

**What to build:** A prototype that shows one GitHub repository's Issues as a Board and lets a person Move a Task, wired straight into the current code on a branch that is never merged. Its only output is the answers below, written into this file under `## Answer`.

**Blocked by:** none

**Status:** ready-for-human

**Type:** prototype

Why: we will put a `Backend` interface (see `CONTEXT.md`) between `board` and where a Board is stored. A second, non-file case must shape that interface before ticket 02 fixes it.

Scope of the spike:

- One repository, hard-wired. Read with the caller's GitHub token, the same as today.
- Task: one issue. Status from a `status:<name>` label, Owners from assignees, labels as labels, a closed issue is Done.
- Move: swap the `status:` label. Refuse when the Status is no longer the one the person saw.
- No Assign, no Docs, no Archive, no tests, no docs. Cut every corner that does not change an answer.

Questions it must answer:

- **Stages:** where does an Issues Board get its Stages and their order: from the labels that exist, a config file in the repository, or defaults?
- **Discovery:** how does fuda know a repository is an Issues Board, and how does it sit next to a markdown Board in the same repository? Proposed Board id: `/github-issues/<owner>/<repo>/`.
- **Task id:** what is the id (`#12`?), and what breaks where the code expects a prefix or a file path (`Card.Prefix`, `Card.Path`).
- **Freshness and cost:** what is the cheap "did anything change" check (ETag of the issue list, newest `updated_at`), and what does a 5-second poll cost in rate limit. How many issues can one Board hold before the full list is too slow.
- **Move conflict:** can the Backend tell who changed the Status first, for the "Ben moved this to Testing just now" message.
- **In review:** does GitHub's own issue-to-PR link replace matching by id in the PR title.
- **Relations and people:** which native fields map to blocked-by, mentions and people (id, name, avatar).
- **Body:** what in an issue body does not render with our markdown (`#123`, mentions, attachment URLs that need a token).
- **Token:** does the fuda GitHub App need more permissions for issues read and write.
- **UI:** every place the client breaks or shows git wording on this Board.
- **Interface:** does the proposed `Backend` shape in ticket 02 hold, and what must change.

- [ ] The spike shows a real repository's issues as a Board in the running app
- [ ] A Move changes the `status:` label on GitHub and a conflicting Move is refused
- [ ] Every question above has an answer under `## Answer`
- [ ] Tickets 02, 03 and 04 are updated from the answers
- [ ] ADR "Backends behind a domain seam" is written and `docs/decisions.md` gets a TBD item with the proposal
- [ ] The spike branch is deleted, nothing from it is merged
