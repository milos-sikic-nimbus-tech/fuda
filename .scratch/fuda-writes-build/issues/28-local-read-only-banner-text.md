# 28: Say why a local board is read-only

**What to build:** A server on the `local` host always shows the board read-only (only the desktop app writes to local folders), but the banner in `client/src/routes/index.tsx` says "This board is read-only: your account can read this repository but not write to it." A local folder has no account, so the text is wrong. For the local host, show: "Local folders are editable only in the desktop app."

**Blocked by:** none

**Status:** ready-for-agent

- [ ] The banner on a local board reads "Local folders are editable only in the desktop app."
- [ ] GitHub and Azure boards keep the current account text
- [ ] A client test covers both texts
- [ ] The Guide page `07-run-locally.md` still matches
- [ ] A screenshot of the new banner on a `local` board is attached to the commit or PR
- [ ] `make check` passes
