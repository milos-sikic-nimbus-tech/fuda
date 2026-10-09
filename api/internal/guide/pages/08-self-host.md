# Self-host

fuda is one container. People log in with a GitHub App, and fuda reads each Board with that
person's own token. fuda keeps no server token and no shared password. Access is decided by GitHub.

## The GitHub App

Create one app for fuda (Settings → Developer settings → GitHub Apps):

- **Callback URL:** `https://<fuda>/auth/github/callback`
- **Expire user authorization tokens:** on. fuda renews tokens with the refresh token and sends the
  person back to login only when that fails.
- **Webhook:** off.
- **Repository permissions:** Contents read, Pull requests read, Metadata read.

Install the app on your organisation for the `fuda-` repositories only. fuda then sees only those.

```sh
docker run -p 8080:8080 -v fuda-data:/data \
  -e FUDA_BASE_URL=https://<fuda> \
  -e FUDA_COOKIE_SECRET=<a long random string> \
  -e FUDA_GITHUB_CLIENT_ID=<client id> \
  -e FUDA_GITHUB_CLIENT_SECRET=<client secret> \
  fuda
```

## Settings

| Variable | Default | Meaning |
|---|---|---|
| `FUDA_SOURCE` | `github` | `github`, `azure`, `local`, or `github,azure` for both hosts at once |
| `FUDA_BASE_URL` | | Public address of fuda, used for the login callback |
| `FUDA_COOKIE_SECRET` | | Encrypts the login cookie. Changing it logs everyone out |
| `FUDA_GITHUB_CLIENT_ID`, `FUDA_GITHUB_CLIENT_SECRET` | | The GitHub App's client id and secret |
| `FUDA_AZURE_CLIENT_ID`, `FUDA_AZURE_CLIENT_SECRET` | | The Microsoft Entra ID app's client id and secret (when `FUDA_SOURCE` includes `azure`) |
| `FUDA_AZURE_TENANT` | `organizations` | The Entra tenant that may log in. Use your tenant id for a single-tenant app |
| `FUDA_TITLE` | repository name | Shown in the header of the `local` Board |
| `FUDA_WATCH_MAIN` | `false` | Also read `main` for the "in prod" badge |
| `FUDA_SYNC_COOLDOWN` | `30s` | Minimum time between manual syncs and pull-request re-reads |
| `FUDA_CACHE_DIR` | `/data` in the image | Where the last read copy is kept; served at startup until the next read |

`/healthz` never asks for a login. The `local` source has no login, so keep such
a server on a private network or behind your proxy's authentication.

## How people see Boards

- Not logged in: every Board page sends you to GitHub or Microsoft to log in, then back.
- After login the Board button in the top bar lists the `fuda-` repositories the app is installed
  on and you can read. `/` opens the only Board, or lists them.
- No access: the page says so. Check that the app is installed on the repository and that your
  account can read it.
- The token lives in an encrypted cookie. fuda stores no user secret.

## Updates

The browser asks for changes about every 5 seconds. fuda checks the head commit with a conditional
request, which costs no rate limit when nothing changed, and reads files only when the head moved.
There are no webhooks.
