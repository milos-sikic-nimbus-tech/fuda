# 11: Desktop app (Wails)

**What to build:** A Wails v3 desktop app mounts the same Go HTTP handler and React UI. Login uses GitHub device flow and MSAL device code; tokens live in the OS keychain, one entry per host. "Open folder…" opens a Local Board and fuda remembers opened folders. Releases go to GitHub Releases and the app updates itself.

**Blocked by:** 09: Both hosts at once, 10: Local Board writes

**Status:** ready-for-agent

- [ ] App builds for macOS and Windows from `make` with only go, pnpm and docker
- [ ] Login, Board picker, Move and Assign work in the app
- [ ] Opened folders are remembered
- [ ] Update from a newer GitHub Release works (checked by hand)
- [ ] ADR 0004 accepted; README explains the unsigned-app warnings

Spec: `.scratch/fuda-writes/spec.md`.
