# 05: Upstream's owner features on person ids

**What to build:** Make upstream's "No owner" filter choice and its owner line check work with person ids from `people.md`. "No owner" sits next to the fork's Humans, Agents and My tasks choices.

**Blocked by:** 02

**Status:** needs-triage

Why: we agreed owners are person ids (idea 3). Upstream built both features on names.

Decided: switch strictly now. An owner that is not an id in `people.md` is a Problem. Agents can fix existing boards quickly given clear steps, so the README carries them.

- [ ] The Owner filter offers No owner, Humans, Agents, My tasks and the people, and No owner stays out of the names pill
- [ ] Assign refuses a value that is not an id in `people.md`, replacing upstream's name-shape check
- [ ] An owner written as a name is listed under Problems, naming the task and the value
- [ ] README has a section an agent can follow to move a board to ids: give each `people.md` entry an `id`, `github` and `kind`; rewrite each task's `owner` from name to id; check fuda shows no owner Problems
- [ ] Tests for each filter choice and for Assign
