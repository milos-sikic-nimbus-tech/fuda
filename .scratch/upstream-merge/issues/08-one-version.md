# 08: One version for all of fuda

**What to build:** A root `VERSION` file (`0.2.0`) is the one source of fuda's version. Every build, the API, the app, the Docker image, the client package and the Guide use it, and checks fail when anything drifts from it.

**Blocked by:** 02

**Status:** ready-for-agent

Why: the Board format (ticket 09) is numbered like fuda's releases ("fuda 0.2 reads 0.2 boards"). Today only the desktop release knows its version, from the git tag; `make` builds `dev`, the server and Docker image have none, and `client/package.json` says `0.0.0`.

- The Makefile reads `VERSION` instead of defaulting to `dev`, and passes it to both binaries (`fuda` and `fuda-desktop`) and to `Info.plist`.
- The Board format fuda reads is a Go constant (`0.2`), so `go test` and `go run` work without build flags. A test in `make check` fails when it is not the major.minor of `VERSION`.
- `/api/version` returns fuda's version and the Board format it reads. The app shows a small "fuda 0.2.0" footer.
- The Docker image is built with `VERSION`, tagged and labeled with it.
- `client/package.json` carries `VERSION`; `make check` fails when they differ.
- The Guide pages that name a format (board config, migration) name it from the constant, or a check fails when they differ.
- The release workflow fails when the pushed tag is not `v` plus `VERSION`.
- CONTRIBUTING's release steps say to bump `VERSION` in a PR, then push the matching tag.

- [ ] `VERSION` exists and `make build`, `make docker` and the desktop targets build that version
- [ ] `/api/version` and the app footer show it
- [ ] Docker image tagged and labeled with it
- [ ] `client/package.json`, the format constant and the Guide are checked against it in `make check`
- [ ] Release workflow refuses a tag that does not match
- [ ] CONTRIBUTING describes the release steps; an ADR records one version for fuda and the Board format
