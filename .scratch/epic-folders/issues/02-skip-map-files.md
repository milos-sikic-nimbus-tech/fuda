# 02: Skip map.md files

**What to build:** A file named `map.md` in the tasks folder or an epic folder is not a task and is not listed under Problems.

**Blocked by:** 01

**Status:** done

Why: wayfinder keeps one `map.md` per epic next to its tickets. It has no `id`, so today it shows as a Problem. Only `map.md` is skipped for now; any other file without an `id` stays a Problem, so a forgotten id is still visible.

- [x] `map.md` is skipped in `tasks/` and `tasks/<epic>/`
- [x] Other id-less files are still Problems; tests cover both
- [x] Guide `03-task-format.md` mentions it
