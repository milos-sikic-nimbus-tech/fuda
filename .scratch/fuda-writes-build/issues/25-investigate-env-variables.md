# 25: Investigate all env variables

**What to build:** A review of every `FUDA_*` environment variable, to decide which ones to keep, rename, merge or drop. Nothing is changed in this ticket. It produces a report and a list of follow-up tickets.

**Blocked by:** none

**Status:** needs-triage

Why now: the writes PR added `FUDA_HOSTS` (with `FUDA_SOURCE` as an old name) and `FUDA_GITHUB_APP_SLUG`, and the desktop app reads some of the same variables with built-in defaults. Local Boards cannot be written on a server, so there may be a case for an opt-in switch. The variables are listed in `.env.example`, `CONTRIBUTING.md`, Guide `08-self-host.md` and `api/internal/config/config.go`, and these lists may not agree.

Questions to answer:

- **Inventory.** Every variable fuda reads, in `config.go` and anywhere else (`os.Getenv`, `os.LookupEnv`, the Makefile, the Dockerfile, `release.yml`). For each: name, default, who needs it (web, desktop, both), and which host.
- **Docs match code.** Does each variable appear in `.env.example`, `CONTRIBUTING.md`, the Guide and `docs/architecture.md`, with the same default and meaning?
- **Dead or ignored variables.** `FUDA_GITHUB_CLIENT_SECRET` is only warned about. What else is read but unused? How long should the warnings stay?
- **Names.** Is the prefix and wording consistent (`FUDA_HOSTS` and the old `FUDA_SOURCE`, `FUDA_LOCAL_PATH`, `FUDA_WATCH_MAIN`, `FUDA_SYNC_COOLDOWN`)? When can the `FUDA_SOURCE` alias go?
- **Validation.** Which variables are checked at start and which fail late? Are the error messages clear? (The cookie secret now needs at least 32 characters.)
- **Web versus desktop.** Which variables does the desktop app honour over its built-in values? Is that what we want?
- **Local Boards on a server.** Should there be an opt-in variable to write Local Boards on a server (for example `FUDA_LOCAL_WRITES`)? Today only the desktop app writes them. Decide, and write the risk down.
- **Secrets.** Which variables hold secrets? Do any end up in logs, the cache dir or the image?

- [ ] A report is written, as a table with one row per variable, under `## Answer` in this file
- [ ] Every mismatch between the code and the docs is listed
- [ ] Each proposed change is a follow-up ticket in this folder, or a TBD entry in `docs/decisions.md`
- [ ] The Local Boards write switch has a decision
