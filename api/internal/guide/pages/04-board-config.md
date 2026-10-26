# Board config files

Four optional files in `docs/board/`. Each is normal markdown with a YAML frontmatter block at
the top that fuda reads; the rest of the file is for people. A missing file means fuda derives
the same thing from the tasks. An invalid block is listed under Problems and ignored.

## stages.md: the columns

```markdown
---
stages:
  - name: Backlog
    status: [backlog]
  - name: In progress
    status: [in progress]
  - name: In review
    when: pr-open
  - name: Merged
    status: [merged]
  - name: Testing
    status: [testing, in qa]
  - name: Validated
    status: [validated]
  - name: Done
    status: [done]
---
```

- The list order is the column order. A column can hold several statuses; the card then shows its
  own status icon.
- `when: pr-open` is the derived column: a task that an open pull request names shows there,
  whatever its status.
- A status no column lists gets its own column at the end, marked unknown.

## labels.md: label groups, values and colours

```markdown
---
groups:
  type:  [bug, feat, impr, refactor, test, chore, docs, question]
  epic:  [billing, onboarding]
  area:  [api, web, db, infra]
colors:
  "type:bug": "#ef4444"
  "epic:billing": "#22c55e"
---
```

- Each group becomes a filter that offers **only these values**, in this order. A value used in
  a task but not listed still shows on the card, it just has no filter entry.
- `colors` is optional. Without it, fuda gives every value in a group its own colour from a
  palette, and `type` keeps fixed defaults (bug red, feat blue, …).
- `theme:` and `epic:` are treated as the same group, shown as Epic. Both labels are deprecated:
  put the task in `tasks/<epic>/` instead, which sets `epic:<folder>`.

## people.md: people and agents

```markdown
---
people:
  - id: nemanja
    name: Nemanja Mudrinic
    aliases: [Nemanja]
    github: nmudrinic
  - id: research-agent
    name: Research agent
    kind: agent
---
```

| Key | Required | Meaning |
| --- | --- | --- |
| `id` | yes | The stable value tasks store in `owner` and `tester`, such as `owner: nemanja`. Unique. |
| `name` | yes | The name fuda shows. |
| `aliases` | no | Extra words the people search finds. Not used to match `owner`. |
| `github` | no | The GitHub login. Shows the person's avatar, maps the logged-in user to them ("Assign to me", "My tasks"), and shows them as the author of an open pull request. |
| `kind` | no | `human` (default) or `agent`. Agents get a bot badge; the Owner filter can pick Humans or Agents. |

With the file, `owner` and `tester` hold ids, matched exactly. A value that is not an id is listed
under Problems for that task, and the card shows it as written. A missing or repeated `id`, or an
unknown `kind`, makes the whole file a Problem; fuda then derives people from the tasks. Assign
writes ids, and only ids listed here.

Without the file, people come from the names used in tasks, matched case-insensitively, and a
unique first name is matched to the full name.

## repos.md: code repositories for In review

```markdown
---
code_repos: [acme/app, acme/web]
---
```

fuda lists the open pull requests of each repository, on any base branch, with the logged-in
person's token. A task whose id is in a pull request's title or branch name (`SS-12` in
`SS-12: add login` or `feature/ss-12-login`) shows in the In review column and is locked. One pull
request can name several tasks. If the person cannot see a repository, the board still loads
without that repository's In review.
