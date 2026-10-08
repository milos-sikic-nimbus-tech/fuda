# Contributing

Thanks for helping. fuda is small on purpose; read [docs/architecture.md](docs/architecture.md)
for how it works and [docs/decisions.md](docs/decisions.md) for why, before larger changes.

## Set up

Prerequisites: Go 1.27+, Node 24+, pnpm 10 (`corepack enable`), golangci-lint v2, Docker (for images).

```sh
cp .env.example .env      # point FUDA_LOCAL_PATH at any checkout that has docs/board/tasks
make dev                  # API on :8080 and the app on http://localhost:5173
```

`make dev` runs two processes: the Go API (`make dev-api`) and Vite (`make dev-client`, proxying
`/api` to :8080). Run them in separate terminals if you prefer.

| Command | Does |
|---|---|
| `make check` | Go lint + tests, client lint + format check + types + tests. Run before committing. |
| `make fmt` | Formats Go and client code. |
| `make build` | Builds the client into `api/internal/web/dist`, then the binary `bin/fuda`. |
| `make docker` | Builds the image (`fuda`). |
| `make clean` | Removes build output and the local cache. |

## Sources

| `FUDA_SOURCE` | Needs |
|---|---|
| `local` | `FUDA_LOCAL_PATH`: a checkout on disk. Shows the working tree, uncommitted edits included. |
| `github` | `FUDA_GITHUB_REPO` (`owner/repo`), `FUDA_GITHUB_TOKEN` (fine-grained: Contents + Pull requests, read). Locally: `FUDA_GITHUB_TOKEN=$(gh auth token)`. |
| `azure` | `FUDA_AZURE_ORG`, `FUDA_AZURE_PROJECT`, `FUDA_AZURE_REPO`, and `FUDA_AZURE_PAT` (Code: Read). Locally you can use `FUDA_AZURE_BEARER` from `az account get-access-token`. |

All settings are in [`.env.example`](.env.example).

## Webhooks locally

Point your git host at `POST /api/webhooks/github` or `POST /api/webhooks/azure` (pushes and pull requests).
A webhook means "re-read now"; fuda never trusts the payload. Locally, `make webhook` (or
`make webhook HOST=azure`) simulates one against the running dev API, signed if `FUDA_WEBHOOK_SECRET` is set.
With the local source you can also set `FUDA_SYNC_INTERVAL=5s` so task edits show up on their own. Without webhooks fuda still checks every
`FUDA_SYNC_INTERVAL` (3 min). Set `FUDA_WEBHOOK_SECRET` to require GitHub's signature or Azure's
`X-Fuda-Secret` header.

## Layout

```
api/      Go module: cmd/fuda (entry) and internal/ packages; serves the API and the built app
client/   Vite + React + TypeScript app; builds into api/internal/web/dist
docs/     how fuda works (architecture) and why (decisions)
```

The Guide pages shown inside the app are in `api/internal/guide/pages/`: they are for people who
use fuda on their repository, not for people changing fuda.

## Rules for code

- **Names explain the code.** No comments that restate it, and no doc comments on functions or
  methods, in Go or TypeScript. A comment only for a non-obvious why.
- **Behaviour lives in `internal/board`.** HTTP handlers parse, call the service and write.
  Interfaces are declared where they are used.
- **No speculative abstractions.** Build what is needed now; leave the rest in
  [docs/decisions.md § TBD](docs/decisions.md#tbd).
- **Repository content never breaks fuda.** Invalid files become Problems; a failed sync keeps the
  last snapshot.
- **Tests where the logic is.** Unit tests for parsing and rules, small fixtures in `testdata/`, no
  mocking framework.
- **Client:** shadcn components (`components/ui` is generated, don't edit it), Tailwind tokens from
  `index.css`, every view state in the URL.
- **The Makefile uses only `go`, `pnpm` and `docker`.**

## Pull requests

1. Branch from `main`. One topic per PR.
2. `make check` passes (Go lint + tests, client lint + format + types + tests).
3. Docs change in the same PR as the behaviour: README, the Guide pages, `docs/architecture.md`,
   and `docs/decisions.md` when a decision changes.
4. Commit messages say what changed and why, in the imperative ("Restore the last board from the
   cache at startup").

`main` is protected: every PR needs an approving review from a maintainer, and review threads must
be resolved before merging.

## License

By contributing you agree that your contribution is licensed under the [MIT License](LICENSE).
