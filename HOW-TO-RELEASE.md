# How to release

A release is a GitHub Release with two macOS disk images: `fuda-desktop_darwin_arm64.dmg` (Apple
Silicon) and `fuda-desktop_darwin_amd64.dmg` (Intel). The `release` workflow
(`.github/workflows/release.yml`) builds them from a tag.

## Before the first release

Set the repository variables `FUDA_GITHUB_CLIENT_ID` and `FUDA_GITHUB_APP_SLUG`; the build fails
without them. `FUDA_GITHUB_BOARDS` and `FUDA_AZURE_BOARDS` are optional and passed to the build when set.

## Cut a release

```bash
echo 0.2.1 > VERSION
git commit -am "Release 0.2.1"
git tag "v$(cat VERSION)"
git push origin main "v$(cat VERSION)"
```

- `VERSION` at the repository root is the one version. The tag must be `v` plus its contents; the
  workflow fails on a mismatch.
- `X.Y.Z` is a stable release. `X.Y.Z-rc.N` (any `-suffix`) is a prerelease, which the in-app update
  check does not offer. Any other tag fails the workflow.
- The version goes into the binary (`main.version`) and into `Info.plist` (`CFBundleVersion`,
  `CFBundleShortVersionString`). The Makefile reads `VERSION`, so local builds carry it too;
  `VERSION=x` on the `make` line overrides it.
- The workflow creates a draft release with generated notes, builds both DMGs, uploads them and
  publishes the release only after both are uploaded. If a build fails, the draft stays hidden: fix
  it and re-run the failed jobs, or delete the draft and the tag.

## What users get

The app checks the newest published release on start and from Help → Check for updates…, and when
it is newer offers to open its release page. It does not install anything. A build made with
`VERSION=dev` on the `make` line never checks; a plain local build carries the `VERSION` file and does.

The DMG is ad-hoc signed, not notarized. Installing it is described in the README.
