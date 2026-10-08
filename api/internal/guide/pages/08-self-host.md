# Self-host

fuda is one container. It needs read access to the repository and, optionally, a webhook from it.

```sh
docker run -p 8080:8080 -v fuda-data:/data \
  -e FUDA_SOURCE=github \
  -e FUDA_GITHUB_REPO=owner/repo \
  -e FUDA_GITHUB_TOKEN=… \
  -e FUDA_AUTH_PASSWORD=… \
  fuda
```

## Settings

| Variable | Default | Meaning |
|---|---|---|
| `FUDA_SOURCE` | `local` | `github`, `azure` or `local` |
| `FUDA_GITHUB_REPO`, `FUDA_GITHUB_TOKEN` | | GitHub repository and token |
| `FUDA_AZURE_ORG`, `_PROJECT`, `_REPO`, `_PAT` | | Azure DevOps repository and personal access token |
| `FUDA_TITLE` | repository name | Shown in the header |
| `FUDA_WATCH_MAIN` | `false` | Also read `main` for the "in prod" badge |
| `FUDA_SYNC_INTERVAL` | `3m` | How often fuda checks the repository when no webhook arrived |
| `FUDA_SYNC_COOLDOWN` | `30s` | Minimum time between manual syncs; webhooks inside it are deferred |
| `FUDA_AUTH_USER`, `FUDA_AUTH_PASSWORD` | `fuda`, empty | Built-in basic auth, off while the password is empty |
| `FUDA_WEBHOOK_SECRET` | empty | Require GitHub's signature or Azure's `X-Fuda-Secret` header |
| `FUDA_CACHE_DIR` | `/data` in the image | Where the last synced copy is kept; served at startup until the first sync |

`/healthz` and `/api/webhooks/*` never ask for basic auth. If your proxy offers authentication,
prefer it and leave `FUDA_AUTH_PASSWORD` empty.

## Tokens

- **GitHub:** a fine-grained token for the one repository with **Contents: read** and
  **Pull requests: read**.
- **Azure DevOps:** a personal access token scoped to **Code: Read**.

## Webhooks

| Host | Where | URL | Events |
|---|---|---|---|
| GitHub | Repository → Settings → Webhooks | `https://<fuda>/api/webhooks/github`, JSON | Pushes, Pull requests |
| Azure DevOps | Project settings → Service hooks → Web Hooks | `https://<fuda>/api/webhooks/azure` | Code pushed (develop, main), Pull request created / updated / merged |

A webhook only means "read again now"; fuda never trusts its payload. Without webhooks the
board still refreshes every `FUDA_SYNC_INTERVAL`.
