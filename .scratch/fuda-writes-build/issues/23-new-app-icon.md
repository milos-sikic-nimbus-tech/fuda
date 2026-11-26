# 23: New app icon

**What to build:** Design a new icon for fuda and use it everywhere: the desktop app (macOS Dock, app switcher, window, Windows taskbar) and the web app favicon. Keep one vector source file in the repo and make every size from it with one command.

**Blocked by:** None

**Status:** needs-triage

Where the icons are now:

- `api/cmd/fuda-desktop/icon.png` and `icon-macos.png`, both 1024x1024 RGBA. `appIcon()` in `main.go` picks one at runtime.
- `api/cmd/fuda-desktop/icon.icns` is copied into `fuda.app` by the Makefile.
- `api/cmd/fuda-desktop/rsrc_windows_amd64.syso` and `rsrc_windows_arm64.syso` hold the Windows icon.
- `client/public/favicon.svg` is the web app icon.

Design notes:

- The icon must be clear at 16x16. Use a simple shape with few details.
- The macOS version uses the standard rounded-square shape and padding.
- The icon must look good on light and dark backgrounds.

- [ ] A new design is agreed before the files are replaced
- [ ] One vector source is in the repo (for example `api/cmd/fuda-desktop/icon-source.svg`), with the command that makes every size
- [ ] `icon.png`, `icon-macos.png`, `icon.icns`, both `.syso` files and `favicon.svg` use the new icon
- [ ] `icon.icns` has all standard sizes (16 to 512, plus @2x); the Windows resource has 16, 24, 32, 48, 64 and 256
- [ ] Screenshots from a real build show the icon in the Dock, the app switcher, the window, the Windows taskbar and the browser tab
- [ ] `make check` passes; `CONTRIBUTING.md` says how to make the icons again
