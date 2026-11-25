---
status: accepted
---

# One version for fuda, and a Board format numbered by it

A root `VERSION` file is fuda's one version. Builds, the API, the app footer, the Docker image, the client package and the Guide follow it, and checks fail when anything drifts. The Board format, the file layout a Board follows, is numbered by the same major.minor: fuda 0.2 reads 0.2 boards. A Board states its format with `format: 0.2` in `stages.md`; a `stages.md` without it is 0.1, and a Board without `stages.md` uses fuda's defaults and counts as current. Agreed on 25 November 2026.

## Why

The merge with upstream changes three things existing boards rely on: the board's repository and branch, owners as person ids, and the Done Stage. Everything before the merge is 0.1. They land together, so there is one migration, "Moving a board from 0.1 to 0.2", written for an agent to follow.

## Considered Options

- One version per change (0.2, 0.3, 0.4): rejected, no fuda release ever read a board with only some of them.
- A format number apart from fuda's release: rejected, two numbers to keep in step for no gain while they change together.
- No marker in the board, spotting 0.1 by its symptoms: rejected, an explicit marker lets fuda say up front which migration a board needs.
- The git tag as the only source of the version: rejected, local builds and the server never knew their version.

## Consequences

- fuda reads a board of any format as well as it can and shows a notice when the format is older or newer than its own; parts it cannot read are Problems.
- The Board format is a Go constant, checked against `VERSION` in `make check`, so tests and `go run` work without build flags.
- The release workflow refuses a tag that is not `v` plus `VERSION`.
- The first release of the merged codebase is 0.2.0.
