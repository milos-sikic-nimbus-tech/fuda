# 16: Folder epics show in the filter when labels.md sets groups

**What to build:** The Epic filter lists every epic folder, also on boards whose `labels.md` configures its own groups.

**Blocked by:** 01

**Status:** done

Why: with `groups:` in `labels.md`, `labelGroups` (`api/internal/board/facets.go`) returns only the configured groups. The `epic:<folder>` label that ticket 01 gives every Task in `tasks/<epic>/` never reaches the facets, so fuda-eliot (groups: only `wayfinder`) has no Epic filter.

- [x] A board with configured groups gets an Epic group with every epic folder that has a Task
- [x] Configured groups stay strict otherwise: other unlisted values get no filter entry
