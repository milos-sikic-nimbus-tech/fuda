# 24: Tell people the GitHub App must be installed

**What to build:** GitHub login works even when the fuda GitHub App is not installed. Then the Board list is empty, or a Board URL gives 404, and nothing says why. Detect a missing install, show a banner with an install link, and document the App install for operators.

**Blocked by:** none

**Status:** ready-for-agent

Decided with the user, 2026-10-13:

- **Install link from a new env var.** `FUDA_GITHUB_APP_SLUG` builds `https://github.com/apps/<slug>/installations/new`. It is required when a GitHub host is configured: fuda does not start without it. This breaks existing deploys, so the release notes say so. The desktop app embeds a slug at build time, the same way as the client id (ticket 20).
- **Detect the install.** After login, call `GET /user/installations` with the person's own token. The behaviour lives in `board`/`source/github`; the handler stays thin.
- **Host status field.** Add `appInstalled` (true / false / unknown) next to `LoggedIn` and `Error`. If the installations call fails, the value is unknown and no banner shows.
- **Banner under the app bar** when `appInstalled` is false. Text: install the fuda GitHub App; if you cannot, send this link to an org owner. People can close it; it comes back on reload or the next login.
- **Existing hints get the link.** The empty Board list (`BoardList.tsx`) and the 404 text (`__root.tsx`) keep their text and add the install link. "Installed, but zero `fuda-` Boards" keeps the current text.
- **Docs for operators:** `README.md`, `.env.example`, Guide `08-self-host.md` and `02-setup-a-repo.md`, `docs/architecture.md`, `docs/decisions.md`. Say where to install the App (the org or user that owns the `fuda-` repos) and add the slug var.

- [ ] fuda refuses to start for a GitHub host without `FUDA_GITHUB_APP_SLUG`
- [ ] The desktop build embeds the slug
- [ ] Host status returns `appInstalled`; unit tests cover installed, not installed, and call failed (unknown)
- [ ] The banner shows only when `appInstalled` is false, has the install link, and can be closed until reload
- [ ] The empty Board list and the 404 page show the install link
- [ ] Docs listed above are updated; release notes mention the new required var
- [ ] `make check` passes
