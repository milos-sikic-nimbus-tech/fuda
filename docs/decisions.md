# Decisions

The choices fuda is built on, each with the reason. These are settled: change one only with a PR
that updates this file and the code together. Open items are at the end under [TBD](#tbd), each
with a proposal.

## Product

- **Read-only.** fuda never writes to a repository. Every change is an ordinary commit by a person
  or their agent, so git stays the single source of truth and fuda needs only read access.
- **develop is the board.** Columns and filters come from develop. Claims are committed to develop
  so everyone sees them. main is optional and only adds the "in prod" badge; it never moves a
  card, because anything on main is on develop too.
- **In review comes from open PRs, matched by content.** A PR delivers a task when its version of
  the task file sets `merged` or adds the PR's number. Branch names are never used: they are
  conventions, the diff is a fact.
- **The repository decides its board.** Columns come from an optional `stages.md`, label groups
  from `labels.md`, people from `people.md`. Without them fuda derives everything from the tasks,
  so a repo works on day one.
- **Never break on repository content.** An invalid file becomes a Problem shown in the UI and is
  left off the board. An unexpected status gets its own flagged column instead of disappearing.
- **A claimed task stays where its status says.** An owner on a backlog task is not a Problem;
  people move tasks explicitly.

## Architecture

- **One binary, one container per repository.** The API and the built client ship together
  (`go:embed`). Several repositories means several containers, not tenants inside one process.
- **No database.** The files are the data. A parsed snapshot lives in memory and a copy on disk
  (`FUDA_CACHE_DIR`) so a restart serves the last board at once.
- **Webhooks plus a fallback interval.** A webhook means "read again now" and its payload is never
  trusted. Without webhooks the repository is checked every `FUDA_SYNC_INTERVAL`. A cooldown stops
  bursts; webhooks inside it are deferred, never dropped.
- **Cheap change detection.** A sync first compares branch head SHAs and downloads files only when
  they changed.
- **Plain Go packages by responsibility.** No ports-and-adapters layering. `board` holds the
  behaviour and declares the small interfaces it needs; sources satisfy them implicitly. Handlers
  stay thin.
- **The standard library first.** `net/http` routing with a tiny middleware chain, `log/slog`,
  plain HTTP clients for GitHub and Azure instead of SDKs. Dependencies: env parsing, YAML,
  goldmark, singleflight.
- **Filtering runs in the browser.** The board payload carries frontmatter only, and every filter
  is a URL parameter, so views are shareable and the server stays simple.
- **Optional built-in basic auth.** One user and password from the environment for hosts without
  proxy auth. `/healthz` and webhooks stay open; webhooks can require a secret instead.
- **Configuration is environment only.** Typed and validated at start; no config files for fuda
  itself.

## Code

- **Names explain the code.** No comments that restate it, and no doc comments on functions or
  methods, in Go or TypeScript. A comment only for a non-obvious why.
- **No speculative abstractions (YAGNI).** Build for the two real repositories, not imagined
  ones. Remove features that don't work rather than leave them half-wired.
- **Tests where the logic is.** Unit tests for parsing and rules, one HTTP smoke test, no mocking
  framework, no end-to-end suite. Small fixtures in `testdata/`.
- **The Makefile uses only `go`, `pnpm` and `docker`,** so anyone can clone and run it.
- **Docs change with the code.** README, Guide pages and these docs are updated in the same PR as
  the behaviour they describe.

## TBD

Decided as a direction, not built. Each has a proposal; a PR that builds one moves it above and
updates [architecture.md](architecture.md).

### History-based dates and insights

- **Settled:** fuda only knows the dates written in frontmatter (`added`, `claimed`, `done`).
  "Time in progress", "merged per week" and an "archive candidate" badge all need the date a
  task's status changed.
- **Open:** how far back to read, and the cost on large repositories.
- **Proposal:** for each task, read the file's commit history once from the git host (GitHub
  commits by path, Azure commits by `itemPath`), find the first commit where each status appears,
  and keep the result in `FUDA_CACHE_DIR` keyed by task and blob SHA, so only changed files are
  read again. Expose the dates per card; the client computes the insights. The archive-candidate
  badge (`FUDA_ARCHIVE_AFTER_DAYS`) returns on top of the same dates.

### CI for pull requests

- **Settled:** contributors open PRs; only maintainers merge, and `main` requires a review.
- **Open:** nothing runs `make check` on a PR yet.
- **Proposal:** a GitHub Actions workflow that installs Go, Node and pnpm and runs `make check`
  on every PR, set as a required status check on `main`.

### More git hosts

- **Settled:** a host is a package under `internal/source` implementing `board.Source` (and
  `board.ReviewSource` for In review).
- **Open:** which host is next.
- **Proposal:** GitLab, when someone needs it: project archive for `docs/`, merge requests and
  their changes for In review.

### Live updates

- **Settled:** the client refetches the board every minute, on window focus, and shortly after
  the Sync button.
- **Open:** whether a push channel is worth its complexity.
- **Proposal:** a server-sent events endpoint that emits the snapshot's head SHAs after each
  swap; the client refetches when they change. No change to the data model.
