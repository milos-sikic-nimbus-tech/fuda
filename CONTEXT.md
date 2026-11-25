# fuda

A kanban board and doc reader over a repository's markdown task files.

## Language

**Board**:
One tasks repository, shown as Stages and Tasks. One fuda can show several Boards.
In code: `board.Service` is one open Board, `board.Registry` holds them, `board.ID` names one.

**Local Board**:
A Board read from a folder on disk instead of from a git host. A Move or an Assign changes the file on disk; the person commits it. Only the desktop app writes to a Local Board.
In code: host `board.HostLocal`, `source/local`, folders kept by `workspace.Folders`.

**Task**:
One unit of work, stored as one markdown file in the repository.
In code: `taskfiles.Task` as parsed, `board.Card` as drawn.
_Avoid_: ticket, issue, card (a card is only how a Task is drawn on the board)

**Doc**:
A markdown file in the repository's docs that is not a Task: anything outside the Board's tasks and Archive folders. Board-config files such as `stages.md` are Docs too. Shown in the Docs view.
In code: `snapshot.docs`; the list is `Service.Docs`, served at `/api/…/docs/list`.
_Avoid_: page, note

**Status**:
The field on a Task that says where it is in the workflow.
In code: the `status` line, changed by `taskfiles.StatusEdit`.

**Done**:
The Status `done`, a step on the board after Validated, with the `done` date. A Task is Done before, and independent of, being archived.
In code: `board/columns.go` default Stages, the `done` line.

**Board format**:
The version of the file layout a Board follows, numbered like fuda's releases (`0.2` is read by fuda 0.2). A Board states it with `format:` in `stages.md`; a `stages.md` without it is format 0.1, and a Board without `stages.md` is current. fuda reads a Board of any format as well as it can and shows a notice when the format is older or newer than its own.
In code: the `format` key in `stages.md`.

**Done Stage**:
The one Stage that shows Done Tasks. It shows every Done Task still on the board, and archived Done Tasks whose `done` date is within the last `days` days (default 14), newest first. Only archived cards in it are locked. In `stages.md` it is the Stage that lists the status `done`. A board whose own `stages.md` has none shows no Done Stage.

**Stage**:
A column on the board. One Stage groups one or more Statuses.
In code: `board.Column`, configured in `stages.md`.
_Avoid_: lane, state

**Owner**:
The person or people working on a Task.
In code: the `owner` line, changed by `taskfiles.ApplyAssign`.
_Avoid_: assignee

**Move**:
Changing a Task's Status from fuda, so it appears in another Stage.
In code: `board.Service.Move`.

**Assign**:
Changing a Task's Owner from fuda.
In code: `board.Service.Assign`.

**Archive**:
The place for Tasks that are off the board: `archive/`, and `archive/<epic>/` for an epic's. A place, not a Status: an archived Task keeps whatever Status it has. Archived Tasks are listed in the Archive view with that Status and are read-only in the app; archiving is a file move in git that changes nothing in the file.
In code: `taskfiles.Task.Archived`, `/api/archive`.

**Problem**:
A repository file fuda could not read. It is shown to people, never fatal.
In code: `board.Problem`.

**In review**:
The Stage a Task is in while an open pull request names it, or when its Status is one the Stage lists (`in review` by default), set by a Move.
In code: `board/reviews.go`, `board.LinkedPR`.

**Read-only Board**:
A Board the logged-in person cannot write to. It shows no drag handles and refuses Move and Assign.
In code: `Service.CanWrite`, `TaskEditor.CanWrite`.

**Host**:
Where a Board lives: `github`, `azure` or `local`.
In code: `board.Host`, `FUDA_HOSTS`.

**Backend**:
What a Board is stored in and written to, for example markdown files in a repository. A Board has exactly one Backend, and that Backend is its only source of truth. A Host is where a person logs in; one Host can serve more than one Backend.
_Avoid_: storage, tracker, provider

**Board id**:
A host and a repository, for example `github/owner/repo`. It is also the path of the Board's URL.
In code: `board.ID`.

**Login**:
How a person proves who they are to a host: device flow (GitHub) or OAuth redirect with PKCE (Azure). The token is kept in a sealed cookie on the web and in the OS keychain on desktop.
In code: `login.NewCookieDevice`, `login.NewCookieOAuth`, `login.NewKeychainDevice`.

**Token**:
The host's access token for one person. fuda has no server token; every read and write uses the caller's own.
In code: `board.WithToken`, `board.TokenFrom`.

## Relationships

- A **Stage** groups one or more **Statuses**; a **Task** has exactly one **Status**.
- A **Move** and an **Assign** are the only changes fuda makes to a **Task**. Creating and editing Task text happens in git.

## Flagged ambiguities

- "ticket" and "issue" were used for **Task**. Resolved: the term is **Task**.
- "assign" was used for setting the **Owner**. fuda's file field is `owner`; the action is **Assign**.
