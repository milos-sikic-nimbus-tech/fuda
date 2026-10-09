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
| `make desktop-macos` | The desktop app for macOS in `bin/fuda-desktop`. Needs a Mac with Xcode command line tools. `ARCH=amd64` for Intel. |
| `make desktop-windows` | The desktop app for Windows in `bin/fuda-desktop.exe`. Cross-builds from any host. |
| `make clean` | Removes build output and the local cache. |

## Sources

| `FUDA_SOURCE` | Needs |
|---|---|
| `github` (default) | A GitHub App: `FUDA_GITHUB_CLIENT_ID`, `FUDA_GITHUB_CLIENT_SECRET`, plus `FUDA_BASE_URL` and `FUDA_COOKIE_SECRET`. People log in; there is no server token. |
| `local` | `FUDA_LOCAL_PATH`: a checkout on disk. Shows the working tree, uncommitted edits included. No login. |
| `azure` | A Microsoft Entra ID app: `FUDA_AZURE_CLIENT_ID`, `FUDA_AZURE_CLIENT_SECRET`, optionally `FUDA_AZURE_TENANT` (default `organizations`), plus `FUDA_BASE_URL` and `FUDA_COOKIE_SECRET`. People log in; there is no server token. |

`FUDA_SOURCE=github,azure` serves both hosts at once, each with its own login and settings. `local`
stands alone.

All settings are in [`.env.example`](.env.example).

## Entra ID app locally

Register an app in Microsoft Entra ID (App registrations) with:

- A Web redirect URI `http://localhost:5173/auth/azure/callback`, and a client secret.
- API permission: Azure DevOps, delegated `user_impersonation`. Whether people can consent to it
  themselves or an admin must grant it depends on the tenant's consent policy.
- For the desktop app: "Allow public client flows" on, and run it with `FUDA_AZURE_CLIENT_ID`. It
  logs in with the device code and keeps the token in the OS keychain. With `FUDA_GITHUB_CLIENT_ID`
  set too, the app serves both hosts.

## GitHub App locally

Create a GitHub App (Settings → Developer settings → GitHub Apps) with:

- Callback URL `http://localhost:5173/auth/github/callback`. Vite proxies `/auth` to the API, so the
  login cookie lands on the page you browse.
- "Expire user authorization tokens" on, and "Request user authorization (OAuth) during installation" off.
- Repository permissions: Contents read, Pull requests read, Metadata read.
- No webhook.

Install it on your `fuda-` repositories, then set `FUDA_GITHUB_CLIENT_ID`, `FUDA_GITHUB_CLIENT_SECRET`,
`FUDA_BASE_URL=http://localhost:5173` and any `FUDA_COOKIE_SECRET` in `.env`.

The browser polls the Board about every 5 seconds. The server asks GitHub for the head commit first
(a conditional request, free when nothing changed) and reads files only when it moved. With the local
source the same poll picks up edits to your working tree.

## Start a new Board

Create a repository named `fuda-<something>` and install the fuda GitHub App on it. Open it in fuda:
with no Task files, or no `develop` branch yet, it shows an empty Board with the default Stages and the hint "No Tasks yet: add
files in `tasks/`". Add Task files under `docs/board/tasks/` and they appear on the next sync.

## Desktop app

Set `FUDA_GITHUB_CLIENT_ID` (a GitHub App with "Enable Device Flow" on) and run `bin/fuda-desktop`, or
build it in with `make desktop-macos GITHUB_CLIENT_ID=<id>`. For a local run, turn off "Expire user
authorization tokens" on the App: the app has no client secret, so it cannot refresh a token.

To release, push a tag such as `v0.1.0`. The `release` workflow builds macOS (arm64, amd64) and
Windows, and uploads `fuda-desktop_<os>_<arch>` archives to the GitHub Release. Set the repository
variable `FUDA_GITHUB_CLIENT_ID` first. The app checks the latest release on start and from
Help → Check for updates…, then replaces its own binary. A build made with `VERSION=dev` never updates.

## Layout

```
api/      Go module: cmd/fuda (web server) and cmd/fuda-desktop (Wails app) and internal/ packages; serves the API and the built app
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
