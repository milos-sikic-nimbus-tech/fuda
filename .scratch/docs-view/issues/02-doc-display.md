# 02: Show a doc the way a task body is shown

**What to build:** The selected doc renders in the Docs view's main pane with the same body styling as a task in the task sheet, but in the page, not in the side sheet.

**Blocked by:** 01

**Status:** ready-for-agent

`routes/docs/$.tsx` already renders `doc.html` with `RichContent` and lists backlinks. Bring it in line with the task sheet (`client/src/components/organisms/TaskSheet.tsx`): pull the shared pieces (title, file path under the title with the full path on hover, body, "Mentioned by" footer) into a component both use, rather than copying markup. Docs have no stage, owner or labels, so that metadata list stays task-only.

- [ ] Doc title, path and body match the task sheet's typography and spacing
- [ ] Shared rendering lives in one component used by the task sheet and the doc page
- [ ] Missing doc still shows the not-found message, inside the new layout next to the tree
- [ ] Checked in the running app in light and dark mode (by hand)
