---
status: accepted
---

# Owners are person ids

`owner` and `tester` hold the `id` of a `people.md` entry, such as `owner: nemanja`. An entry is
`id` (required, unique), `name` (shown), `aliases` (search only), `github` (login) and `kind`
(`human` or `agent`, default `human`). Without `people.md`, people are still derived from the Tasks.

## Why

Linking a GitHub login to a person needs one stable key. Display names change and aliases made
matching fuzzy: "Nemanja", "nemanja" and "Nemanja Mudrinic" all meant one person, so nothing could
safely say which person a login or a PR author is. Agents such as `research-agent` own Tasks too and
need to be told apart from humans.

## Considered Options

- Keep names and fold aliases, adding `github` next to them: rejected. Matching stays fuzzy, and a
  renamed person breaks every Task that used the old name.
- Accept ids and fall back to aliases for old Tasks: rejected. Boards would never be cleaned up;
  strict ids with Problems show what to fix.

## Consequences

- Breaking: every `people.md` needs an `id` per entry, and Tasks must use those ids. A Task with any
  other value is listed under Problems; its card still shows the value.
- Assign writes ids and refuses an id that is not listed.
- The logged-in user's login maps to a person: "Assign to me" and a "My tasks" filter.
- Cards show the GitHub avatar, an agent badge, and the author of an open PR as a person when the
  login matches. The Owner filter can pick Humans or Agents.
- Azure DevOps logins are not mapped yet; only `github` exists.
