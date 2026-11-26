# 07: Warnings next to Problems

**What to build:** A warning level for notices about files that still work but need attention, shown apart from Problems (files fuda could not read).

**Blocked by:** 01

**Status:** needs-triage

Why: ticket 01 lists every deprecated `epic:`/`theme:` label as a Problem. Mixing harmless notices with broken files fills the Problems list and teaches people to skip it. Ticket 16 in `fuda-writes-build` (warn when people are missing) needs the same kind of notice.

- Decide the API shape: a severity on each Problem, or a separate list.
- The UI shows warnings apart from Problems, so a board with only warnings does not look broken.
- Move the `epic:`/`theme:` deprecation notice from ticket 01 to a warning.
- Check whether ticket 16 in `fuda-writes-build` should use it.

- [ ] The API returns warnings separately from Problems, with tests
- [ ] The client shows them apart from Problems
- [ ] The `epic:`/`theme:` deprecation is a warning
- [ ] Guide and `docs/architecture.md` describe warnings
