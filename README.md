# fuda

A kanban board and doc reader for teams that keep their tasks as markdown files in git.

Each task is one file in `docs/board/tasks/` with YAML frontmatter (`status`, `owner`, `labels`, …).
fuda reads them from GitHub, Azure DevOps or a local checkout, and shows them as a board with filters,
a task reader and the repo's docs. People (and their AI agents) move tasks with ordinary commits.
On GitHub a person can also drag a card to another Stage: fuda commits that one `status` change as
the logged-in user, and pick a Task's Owners from the short names in `people.md`.

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
  -e FUDA_BASE_URL=https://fuda.example.com \
  -e FUDA_COOKIE_SECRET=<a long random string> \
  -e FUDA_GITHUB_CLIENT_ID=<GitHub App client id> \
  -e FUDA_GITHUB_CLIENT_SECRET=<GitHub App client secret> fuda
```

People log in with the fuda GitHub App. fuda lists every `fuda-` repository they can read and the
App is installed on, and reads each Board with that person's own token. There is no server token and
no shared password. Boards open at `/github/<owner>/<repo>/`; `/` lists them, or opens the only one.
Set up the App in the [Self-host guide](api/internal/guide/pages/08-self-host.md).

## Desktop app

Download the archive for your system from the latest GitHub Release, unpack it and open `fuda-desktop`
(macOS) or `fuda-desktop.exe` (Windows). Log in with GitHub: fuda shows a code and opens GitHub, where
you type it. The login token stays in your OS keychain. The app updates itself from new releases.

The builds are not signed, so your system warns the first time:

- **macOS:** "fuda-desktop cannot be opened". Right-click the file, choose Open, then Open again. Or
  run `xattr -d com.apple.quarantine fuda-desktop`.
- **Windows:** SmartScreen says "Windows protected your PC". Choose More info, then Run anyway.

## License

MIT
