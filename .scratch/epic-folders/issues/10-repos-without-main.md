# 10: Repos without main that are not fuda- repos

**What to build:** Decide and build what fuda does with a repo that has no `main` branch and whose name does not start with `fuda-`.

**Blocked by:** 05

**Status:** needs-triage

Why: ticket 05 shows an error page only for `fuda-` repos. For other repos the behaviour is undefined today.

Options to weigh: the same error page; hiding the repo from the board picker; a mark in the picker. Consider that the GitHub and Azure listings may include many repos that were never meant as boards.

- [ ] The behaviour is decided and recorded in `docs/decisions.md`
- [ ] It is built, with tests, and the Guide says what happens
