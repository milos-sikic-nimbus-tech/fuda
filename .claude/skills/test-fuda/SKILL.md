---
name: test-fuda
description: Run fuda against a local fixture repo and drive it in the built-in browser with screenshots. Use for any change to the app: to try it, check it works, debug it, or after adding UI features.
---

# Test fuda in a browser

The `local` host needs no login and no env file. Never read the real env file (a hook blocks it).

## Start

```bash
S=$TMPDIR/fuda-test   # or the session scratchpad
rm -rf $S && mkdir -p $S && cp -R .claude/skills/test-fuda/fixture $S/fixture
git -C $S/fixture init -q && git -C $S/fixture add -A && git -C $S/fixture -c user.email=a@b -c user.name=t commit -qm init
go -C api build -o $S/fuda ./cmd/fuda
FUDA_HOSTS=local FUDA_LOCAL_PATH=$S/fixture FUDA_ADDR=:8099 FUDA_CACHE_DIR=$S/cache $S/fuda > $S/fuda.log 2>&1 &
```

- Client changes: run `pnpm -C client build` first (the binary embeds `api/internal/web/dist`).
  Or run `pnpm -C client dev` for Vite on :5173; it proxies `/api` to :8080, so start the server
  with `FUDA_ADDR=:8080` then.
- Go changes: rebuild and restart. Stop the old server first: `pkill -f "$S/fuda"`.
- Open `http://localhost:8099` with `mcp__Claude_Browser__navigate`. The board is at `/local/fixture`.

## Drive

`find` + `computer` (click, type, screenshot) to act; `get_page_text` / `read_page` to check
content; `read_console_messages` and `read_network_requests` for client errors; `$S/fuda.log` for
server errors. Clicking a card opens the task drawer (URL gets `?task=T-1`). Wait a moment before
screenshotting animations.

## Limits

- A server shows a local board read-only. Move and Assign only work in the desktop app.
- The fixture has no `stages.md`, so stages derive from the tasks. Add files to the fixture to test
  config, labels, people, Problems (invalid task files) and docs pages.
- GitHub and Azure login are not covered.

## After a UI change

Always run the app, exercise the feature, and take screenshots (mobile too via `resize_window`,
dark mode if relevant). Show them to the user and say what each one confirms works. Reset the
viewport with preset `desktop` afterwards, and stop the server when done.
