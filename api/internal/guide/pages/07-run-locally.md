# Run locally

Point fuda at a checkout on disk to see the board before deploying, or to preview task edits
before you push them.

```sh
git clone <fuda> && cd fuda
cp .env.example .env          # set FUDA_LOCAL_PATH to your repository
make dev                      # API on :8080, app on http://localhost:5173
```

With the local source the board shows your working tree, uncommitted edits included. If the
checkout has `origin/main` and `FUDA_WATCH_MAIN=true`, "in prod" works too.

| Setting | Use |
|---|---|
| `FUDA_SOURCE=local`, `FUDA_LOCAL_PATH=…` | Read a checkout |
| `FUDA_SOURCE=github`, `FUDA_GITHUB_REPO`, `FUDA_GITHUB_TOKEN=$(gh auth token)` | Read GitHub with your own login |
| `FUDA_SYNC_INTERVAL=5s` | Pick up local edits on their own |

`make webhook` (or `make webhook HOST=azure`) sends a simulated webhook to the running API,
signed when `FUDA_WEBHOOK_SECRET` is set.

Prerequisites: Go 1.27+, Node 24+ and pnpm 10. `make check` runs every lint and test.
