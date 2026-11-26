# 02: The Backend interface, with markdown behind it

**What to build:** `board` declares a `Backend` interface that speaks Tasks, Stages, Move and Assign. Today's markdown-files code becomes the first Backend in `internal/backend/markdown`. Nothing changes for a person using fuda: same Boards, same URLs, same payload.

**Blocked by:** 01

**Status:** needs-triage

The shape below is the proposal agreed before the spike. Ticket 01 may change it; update this ticket from its answers, then set `ready-for-agent`.

```go
type Backend interface {
    Load(ctx) (Data, error)
    Task(ctx, id) (TaskDetail, error)
    Move(ctx, id, from, to Status) error
    Assign(ctx, id, from, to []PersonID) error
    CanWrite(ctx) (bool, error)
    Capabilities() []Capability
}

type Data struct {
    Tasks  []Task        // id, title, Status, Owners, labels, blocked-by, mentions, linked PRs, dates
    Stages []Stage       // name and the Statuses it groups
    People []Person      // id, name, avatar URL, kind
    Labels []LabelGroup  // name, values, colours
}
```

Rules:

- **The Backend returns normalised data.** Everything that reads a markdown convention moves into the markdown Backend: `people.md` matching and its Problems (`board/people.go`), `group:value` labels and `labels.md` colours (`board/facets.go`), mention and `blocked_by` parsing (`board/references.go`), matching a pull request to a Task by id (`board/reviews.go`), body rendering and link rewriting, search over bodies.
- **`board` keeps the product rules,** applied once for every Backend: placing Tasks in Stages, the In review lock, the flagged Stage for an unknown Status, valid Move targets, filter lists, epic counts, the payload. `board` no longer imports `taskfiles`.
- **Conflicts.** Move carries the Status the person saw and Assign the Owners they saw. The Backend refuses with a conflict that says who changed it and to what. File versions and the retry loop (`board/move.go`, `board/assign.go`) move inside the markdown Backend.
- **Sync.** Each Backend owns how it stays fresh. Today's snapshot, disk cache and cooldown (`board/sync.go`, `board/cache.go`) become a helper a Backend may use; markdown uses it.
- **Two rules bind every Backend:** a failed read keeps the last good state, and the caller's access is checked before anything cached is served.
- **Capabilities.** Docs, Archive, In review, Problems and epic maps are optional. Markdown declares all of them. Their read methods live on optional interfaces, not on `Backend`.
- **File stores stay.** `Source`, `TaskEditor` and `ReviewSource` move to the markdown package, which declares them; `source/{github,azure,local}` satisfy them unchanged.
- `workspace` builds a markdown Backend per Board; `board.Registry` holds Boards as before.

- [ ] `board` declares `Backend` and builds the Board from `Data`; unit tests for the rules use a fake Backend, not markdown files
- [ ] `internal/backend/markdown` passes the existing parsing, sync, Move and Assign tests, moved with the code
- [ ] `board` does not import `taskfiles` or `markdown`
- [ ] A contract test any Backend can run: keeps the last good state on a failed read, refuses a stale Move and a stale Assign
- [ ] The API payload and the HTTP smoke test are unchanged
- [ ] `CONTEXT.md` "In code" lines, `docs/architecture.md` (packages, diagrams) and `docs/decisions.md` (the TBD item moves up; "Plain Go packages by responsibility" is reworded) are updated
- [ ] `make check` passes
