# 18: Open folder in the web UI

**What to build:** "Open folder…" on desktop lives in the File menu today (`api/cmd/fuda-desktop/folder.go`). Move it into the web UI, so people find it where they pick a Board. Not refined yet.

**Blocked by:** None

**Status:** needs-triage

- [ ] To be written when the task is refined

Spec: `.scratch/fuda-writes/spec.md`.

## Open questions (to grill)

Not decided yet: where it sits (the Board picker, the empty state, both), how the page asks the desktop shell for the folder dialog (a small endpoint or a Wails binding), what the web server shows instead (nothing, or a hint about `FUDA_LOCAL_PATH`), and whether the File menu item stays. Related to 17: workspace picker redesign.
