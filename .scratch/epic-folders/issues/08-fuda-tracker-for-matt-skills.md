# 08: fuda as the issue tracker for Matt Pocock's skills

**What to build:** A tracker doc that lets Matt Pocock's skills (`wayfinder`, `to-tickets`, `to-spec`, `triage`, `implement`) use a fuda board as their issue tracker, as they already use local markdown or GitHub Issues. fuda ships it as a copy-paste section of the Agent guide page; a board keeps its one copy in `docs/agents/issue-tracker.md`, and code repositories point to it.

**Blocked by:** 01, 02, 12

**Status:** ready-for-human

Why: eliot and the next projects plan with these skills. They find tickets through an issue-tracker doc, and none describes a fuda board, so they fall back to `.scratch/` in the code repository, which ticket 12 rules out.

- **One copy.** The doc lives in the board repository. A code repository's `CLAUDE.md` only points to it; the developer's board checkout path is the `fuda board:` line in their `CLAUDE.local.md` (gitignored, different on each machine).
- **Location:** tickets in `<board>/docs/board/tasks/<epic>/<ID>-<slug>.md`; an optional map in `<board>/docs/board/tasks/<epic>/map.md` with `label: wayfinder:map` frontmatter.
- **Ids:** one sequence per repository (`EL-N` for eliot), the highest number in `tasks/` and `archive/` plus one, after a pull. No per-epic prefixes: they collide.
- **Statuses:** the triage roles are statuses with the same names, listed in `stages.md`; won't-do is `wontfix` in its own column and never archived, because fuda counts archived tasks as finished work (ticket 03).
- **Wayfinding operations:** claim = `owner`, `status: "in progress"`, `claimed`, pushed before work; blocking = `blocked_by`; frontier = backlog, unowned, every blocker archived or `wontfix`; resolve = `## Answer`, move to `archive/<epic>/` as `done`, a gist line in the map's Decisions so far.
- **Commits:** as in the `fuda-tasks` skill, pushed straight to `main`.
- **The conversion of eliot's existing tickets** is not part of fuda: it is a one-shot prompt run in an eliot session, `.scratch/epic-folders/eliot-conversion-prompt.md`.

- [x] Agent guide page: Board setup (`stages.md`, `labels.md`), the tracker doc, and the code-repository pointer plus `CLAUDE.local.md` line
- [ ] Ticket 12's code-repository rule block points to the tracker doc and `CLAUDE.local.md`
- [ ] A wayfinder session in eliot, using the doc, claims and resolves a ticket in `fuda-eliot`, and fuda shows the map and the archived ticket with no Problems
