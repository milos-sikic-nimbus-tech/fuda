# 31: Copy button for the task ID

**What to build:** People paste a task ID into branch names, commits, PRs and chat. Add a small copy button next to the ID in the task sheet header (`client/src/components/organisms/TaskSheet.tsx`, the `SheetDescription` that shows `search.task`). Clicking it puts the ID, e.g. `FUDA-12`, on the clipboard and shows "Copied" for about 1.5 seconds, like the code block copy button in `client/src/components/molecules/RichContent.tsx`. Works on read-only boards too.

**Blocked by:** none

**Status:** ready-for-agent

- [ ] A copy button sits next to the task ID in the task sheet header, with an accessible label ("Copy task ID")
- [ ] Clicking it copies just the ID and briefly shows that it was copied
- [ ] It works on read-only and local boards
- [ ] A client test covers the copy
- [ ] A screenshot of the button is attached to the commit or PR
- [ ] `make check` passes
