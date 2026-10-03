# fuda

A read-only kanban board and doc reader for teams that keep their tasks as markdown files in git.

Each task is one file in `docs/board/tasks/` with YAML frontmatter (`status`, `owner`, `labels`, …).
fuda reads them from GitHub, Azure DevOps or a local checkout, and shows them as a board with filters,
a task reader and the repo's docs. It never writes to the repo: people (and their AI agents) move
tasks with ordinary commits.

> Status: early development. The Guide inside the app (coming) explains how to set up a repo.

## Develop

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

### Sources

| `FUDA_SOURCE` | Needs |
|---|---|
| `local` | `FUDA_LOCAL_PATH`: a checkout on disk. Shows the working tree, uncommitted edits included. |
| `github` | `FUDA_GITHUB_REPO` (`owner/repo`), `FUDA_GITHUB_TOKEN` (fine-grained: Contents + Pull requests, read). Locally: `FUDA_GITHUB_TOKEN=$(gh auth token)`. |
| `azure` | `FUDA_AZURE_ORG`, `FUDA_AZURE_PROJECT`, `FUDA_AZURE_REPO`, and `FUDA_AZURE_PAT` (Code: Read). Locally you can use `FUDA_AZURE_BEARER` from `az account get-access-token`. |

All settings are in [`.env.example`](.env.example).

## Layout

```
api/      Go module: cmd/fuda (entry) and internal/ packages; serves the API and the built app
client/   Vite + React + TypeScript app; builds into api/internal/web/dist
```

## Run the image

```sh
docker run -p 8080:8080 -v fuda-data:/data \
  -e FUDA_SOURCE=github -e FUDA_GITHUB_REPO=owner/repo -e FUDA_GITHUB_TOKEN=… fuda
```

Protect it with your proxy's auth, or set `FUDA_AUTH_PASSWORD` (and optionally `FUDA_AUTH_USER`, default
`fuda`) for built-in basic auth. `/healthz` and `/api/webhooks/*` stay open either way.

## License

MIT
