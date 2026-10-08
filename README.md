# fuda

A read-only kanban board and doc reader for teams that keep their tasks as markdown files in git.

Each task is one file in `docs/board/tasks/` with YAML frontmatter (`status`, `owner`, `labels`, …).
fuda reads them from GitHub, Azure DevOps or a local checkout, and shows them as a board with filters,
a task reader and the repo's docs. It never writes to the repo: people (and their AI agents) move
tasks with ordinary commits.

> Status: early development, in use on two repositories.

## Documentation

- **Using fuda on your repository:** the Guide inside the app, or
  [`api/internal/guide/pages`](api/internal/guide/pages): set up a repo, task format, board
  config, workflow, the agent guide, self-hosting.
- **How it works:** [docs/architecture.md](docs/architecture.md).
- **Why it works that way:** [docs/decisions.md](docs/decisions.md).
- **Changing fuda:** [CONTRIBUTING.md](CONTRIBUTING.md).

## Quick start

```sh
cp .env.example .env      # point FUDA_LOCAL_PATH at any checkout that has docs/board/tasks
make dev                  # API on :8080 and the app on http://localhost:5173
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
