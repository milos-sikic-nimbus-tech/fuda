# 01: Epic folders

**What to build:** fuda reads task files one level deep under the tasks folder, `docs/board/tasks/<epic>/*.md`. The folder is the task's epic: the task gets the label `epic:<folder>`.

**Blocked by:** none

**Status:** done

Why: a team (eliot) wants a separate tickets repo with one folder per epic, flat inside. Today `nestedUnder` in `api/internal/taskfiles/taskfiles.go` reports every file in a subfolder as a Problem.

- One level only. A file two folders deep stays a Problem.
- Today an epic is only an `epic:<name>` or `theme:<name>` label (`board/facets.go` folds both into the Epic group used by the Filter menu, Insights and List view). For a task in an epic folder, any `epic:` or `theme:` label in its frontmatter is ignored, so the folder is the one source and moving a ticket to another epic is just `git mv`. Tasks directly in `tasks/` keep these labels for now.
- **Deprecate `epic:` and `theme:` labels.** Every task that has one, in a folder or not, still loads, but it is listed as a Problem: "`epic:` and `theme:` labels are deprecated and will be removed; move the task into `tasks/<epic>/`". For now it is a regular Problem with clear wording; a separate warning level is ticket 07. Record the deprecation and a removal plan under TBD in `docs/decisions.md`.
- Duplicate-id and filename-starts-with-id rules apply across all folders: ids are unique in the whole repo.
- Archive: decide whether `archive/` keeps epic subfolders too, and write it down.
- Writes (moves, owner changes) must commit to the file's real path inside its epic folder.

- [x] Files in `tasks/<epic>/` load as tasks; deeper files are Problems
- [x] Tasks in a folder get `epic:<folder>` and their own `epic:`/`theme:` labels are ignored; flat tasks keep them; unit tests cover both
- [x] Any `epic:` or `theme:` label produces a deprecation notice; a test covers it
- [x] Writes keep working for tasks inside epic folders
- [x] Guide `02-setup-a-repo.md`, `03-task-format.md` and `docs/architecture.md` describe the layout and mark `epic:` and `theme:` labels as deprecated
