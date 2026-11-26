---
status: accepted
---

# fuda's own tickets stay in `.scratch/` for now

fuda's own tickets stay as markdown files under `.scratch/<feature>/issues/`, with ADRs in `docs/adr/` and the glossary in `CONTEXT.md`. Upstream had removed them after reviewing the first fork PR, because they pointed every contributor's agent at one person's setup. Agreed with the upstream maintainer on 25 November 2026, on the condition that the tickets later move onto a fuda board so fuda shows its own work (`.scratch/own-board/issues/01-scratch-to-board.md`).

## Considered Options

- Upstream's rules now (no tickets, `decisions.md` as the only place for why): rejected for now, the merge plan itself lives in these tickets.
- A `fuda-fuda` board repository at once: deferred to the own-board ticket, which decides where the board lives.

## Consequences

- Changes still go through PRs that the maintainer reviews (CODEOWNERS and the PR template stay).
- `decisions.md` keeps product and architecture decisions and the TBD list; ADRs hold one decision each.
- Moving `CONTEXT.md` to `docs/glossary.md` and folding ADRs into `decisions.md` are decided with the own-board ticket.
