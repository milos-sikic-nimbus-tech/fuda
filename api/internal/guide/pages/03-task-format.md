# Task format

One file per task in `docs/board/tasks/`. The filename starts with the id: `<ID>-<slug>.md` or
`<ID>.md`. A file named `map.md` in `tasks/`, `archive/` or an epic folder is not a task and is skipped; any
other file without an `id` is listed under Problems.

An epic folder's `map.md` with `label: wayfinder:map` in its frontmatter is the epic's map. Select that
epic in the Filter menu and the **Map** button next to it opens the map text with the epic's task count
per column and its archived tasks. A `map.md` without the label is ignored; one with broken frontmatter
is listed under Problems.

```markdown
---
id: "SS-12"
title: "Invitation email failure rolls back the insert"
status: "backlog"
blocked_by: ""
owner: ""
tester: ""
labels: ["type:bug", "area:api"]
added: "2026-10-01"
claimed: ""
pr: []
done: ""
---
# SS-12: Invitation email failure rolls back the insert

What, why, where (files and symbols).

## Done when
The check that proves it is finished.

## Evidence
What was checked or changed, with dates and commit shas.
```

## Fields

| Field | Required | Meaning |
|---|---|---|
| `id` | yes | Unique, never reused |
| `title` | yes | One line |
| `status` | yes | Decides the column. Default set: `backlog`, `in progress`, `merged`, `testing`, `validated`; `done` only in `archive/` |
| `blocked_by` | no | Task ids (`"SS-3, SS-9"`) or text. Non-empty means blocked, at any status |
| `owner`, `tester` | no | One or more names (`"Ron, Chinmay"` or a list) |
| `labels` | no | A list of `group:value` strings, e.g. `type:bug`, `area:api`. The epic comes from the task's folder (`tasks/<epic>/`); `epic:` and `theme:` labels are deprecated: the task still loads, but each one is listed under Problems, and inside an epic folder they are ignored |
| `added`, `claimed`, `done` | no | `YYYY-MM-DD`; `added` is what Insights and the Added filter use |
| `pr` | no | Pull request numbers, or a branch name |
| anything else | no | Custom fields, shown in the task panel, not used for columns or filters |

Values are any valid YAML. Quoted JSON style, as above, is the safest for simple checkers.

## Body

The body is GitHub-flavoured markdown, including tables and mermaid diagrams. These sections
are recommended:

- `## Done when`: shown first in the task panel.
- `## Evidence`: what was checked, with dates and commit shas.
- `## Group context`: once per label group, for background that many tasks share.
- `## Answer`: on a question once it is answered.

Task ids written anywhere in a body become links, and the panel lists which tasks block or
mention each other.

## When a file is wrong

fuda skips the file and lists it under **Problems** with the reason: broken YAML, a missing
`id`, `title` or `status`, a filename that does not start with the id, `done` outside the archive,
or a duplicate id. Everything else on the board keeps working.
