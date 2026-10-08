# Shipping the Wails v3 desktop app at $0

Researched 2026-10-08. Every claim has a source. Items marked **unverified** were not confirmed from a primary source.

## Answer

- **The real $0 cost is user friction.** Unsigned builds run on all three OSes, but macOS and Windows show warnings that users have to click past. Linux has no equivalent gate.
- **macOS 15 Sequoia removed the Control-click "Open" bypass.** Users now open System Settings > Privacy & Security and allow the app there (or run `xattr` in a terminal). The only real fix is Developer ID signing plus notarization, which needs the Apple Developer Program at **$99/year**. No free route exists for an individual.
- **Windows:** unsigned apps show "Windows protected your PC"; users click "More info" then "Run anyway". Each new unsigned version starts again with zero reputation. Smart App Control on Windows 11 can block unsigned files outright.
- **Free Windows signing:** SignPath Foundation signs OSS projects for free, with conditions: OSI license, no proprietary parts, automated CI builds from your repo, MFA, defined roles, and a "Code signing policy" page. Azure Artifact Signing (formerly Trusted Signing) costs **from $9.99/month** and needs a paid Azure subscription. Individual developers are accepted **only from the US or Canada**; organizations from the US, Canada, EU, UK and a few other countries.
- **Signing does not stop SmartScreen right away.** Even signed (and EV-signed) apps warn until reputation builds over "several weeks and hundreds of clean installs".
- **Package managers do not remove the cost on macOS.** Homebrew deprecated unsigned casks and said it would disable casks failing Gatekeeper in Sept 2026. A self-hosted tap is the fallback. winget and Scoop have no signing requirement that I could find (**unverified**).
- **Wails v3 has a built-in updater (`app.Updater`)** with a GitHub Releases provider, signature/checksum verification and a binary swap with rollback. Wails v3 is still beta (v3.0.0-beta.28), so treat the updater as pre-release.
- **GitHub Releases is free:** up to 1000 assets per release, each under 2 GiB, and no limit on total size or bandwidth.

## macOS: Gatekeeper

- In Sequoia, users "will no longer be able to Control-click to override Gatekeeper" for software that isn't signed correctly or notarized. They must go to System Settings > Privacy & Security to allow it. Source: https://developer.apple.com/news/?id=saqachfa
- In practice, the user tries to open the app, gets blocked, goes to Privacy & Security, clicks the allow button, confirms again and enters their password. Source (secondhand blog): https://mjtsai.com/blog/?p=43993. The exact button label ("Open Anyway") is **unverified** from Apple.
- `sudo spctl --global-disable` no longer works in Sequoia (secondhand): https://appleinsider.com/inside/macos-sequoia/tips/whats-changed-in-runtime-protection-for-macos-sequoia
- Terminal workaround: remove the quarantine attribute (`xattr -d com.apple.quarantine Fuda.app`). This is general knowledge and is what Homebrew's `--no-quarantine` did (per https://brew.sh/2025/11/12/homebrew-5.0.0/). **Unverified** against Apple docs.
- Apple Developer Program: 99 USD per membership year. Fee waivers apply only to eligible nonprofits, educational institutions and government entities: https://developer.apple.com/programs/whats-included/. Apple recommends notarization for software distributed outside the Mac App Store: https://developer.apple.com/news/?id=saqachfa. That Developer ID and notarization are included in the $99 fee is standard knowledge but was not stated on the fetched page (**unverified** there).
- An unsigned arm64 build still needs at least an ad-hoc signature to run on Apple silicon. **Unverified** here; Homebrew cites "the requirement for valid signatures on Apple silicon".

## Windows: SmartScreen

- Unsigned: "Windows protected your PC". The user must choose "Run anyway". Enterprise policy can prevent continuing at all. Self-signed behaves the same as unsigned. Source: https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation
- Unsigned files build reputation from zero for every new version: same page.
- "EV certificates no longer bypass SmartScreen": same page.
- Signed apps still warn until reputation builds, which can take "several weeks and hundreds of clean installs": same page.
- Windows 11 Smart App Control "will block execution of unsigned files unless the file has a positive reputation": same page.
- Microsoft Store distribution avoids warnings entirely: same page. Store fees are not researched here.

## Free or cheap signing

### SignPath Foundation (free for OSS)
Conditions from https://signpath.org/terms:
- OSI-approved license for all components, with no commercial dual-licensing and no proprietary components (system libraries allowed).
- Actively maintained, already released in the form to be signed, and functionality described on the download page.
- Each signed binary must come from an automated build of the repo, verified via origin verification: https://docs.signpath.io/origin-verification
- MFA for all team members, plus defined author, reviewer and approver roles. Every signing request is approved manually.
- A "Code signing policy" section on the home page and download page.
- The certificate is issued to SignPath Foundation, so it shows as the publisher.
- Builds from a public repository: stated in secondary write-ups. Not explicit on the terms page (**unverified**).
- Platform: SignPath is known for Windows Authenticode. Whether it can sign or notarize macOS apps is **unverified**; assume it does not cover Gatekeeper.
- Implication: fuda's license must be OSI-approved for this to apply.

### Azure Artifact Signing (formerly Trusted Signing)
- Starts at $9.99/month (Microsoft Learn): https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation. Basic includes 5,000 signatures per month and Premium 100,000: https://azure.microsoft.com/en-us/pricing/details/artifact-signing/
- Billing is not pro-rated, and free, trial or sponsored Azure subscriptions are not supported: https://learn.microsoft.com/en-us/azure/artifact-signing/faq
- Eligibility: Public Trust certificates are available to organizations in the US, Canada, EU, UK, Australia, New Zealand, Japan, South Korea, Singapore, Switzerland, Norway and Israel. **Individual developers must be in the US or Canada.** Source: https://learn.microsoft.com/en-us/azure/artifact-signing/quickstart
- An EU-based individual can't use it. An EU-registered organization can, after identity validation that takes 1 to 20 business days (same page).
- No EV certificates: FAQ above.
- Windows only.

## Package managers

- **Homebrew cask:** Homebrew 5.0.0 deprecated casks without codesigning and said it "will disable all Homebrew/homebrew-cask casks that fail Gatekeeper checks in September 2026". `--no-quarantine` is deprecated too: https://brew.sh/2025/11/12/homebrew-5.0.0/. The acceptable-casks rules say a cask must not require Gatekeeper to be bypassed: https://docs.brew.sh/Acceptable-Casks. So the official cask repo effectively needs notarization. A self-hosted tap (your own repo) is not covered by those rules, but users would still hit Gatekeeper. Whether the September 2026 disable was actually enforced is **unverified**.
- **winget:** you submit a manifest PR to microsoft/winget-pkgs, and automated validation checks policy and malware: https://learn.microsoft.com/en-us/windows/package-manager/package/repository. I found no signing requirement (**unverified**). Installing via winget may still trigger SmartScreen or Smart App Control (**unverified**).
- **Scoop:** buckets are plain git repos, and you can host your own bucket. Its signing policy and SmartScreen behaviour were not researched from primary sources (**unverified**).
- **Linux:** no OS-level signing gate. Shipping AppImage, .deb or a tarball on GitHub Releases costs nothing. General knowledge, not sourced.

## Wails v3 auto-update

- Built in as `app.Updater`, which is already wired into every `*application.App`. You call `Init` after `application.New` and before `app.Run()`: https://v3.wails.io/tutorials/04-self-update-a-wails-app/
- Providers: GitHub Releases (picks assets by OS and architecture, optional checksum file), keygen.sh, Sparkle AppCast, a Wails Update Manifest JSON endpoint, or a custom provider. Providers run as an ordered fallback chain. Source: https://docs-preview.wails.io/guides/updater/
- Verification: Ed25519, Ed25519ph or ECDSA P-256 signatures, or digest-only SHA-256/512. Digest-only guards against corruption, not tampering. `Config.PublicKey` is the trust root. The `wails3 updater` CLI generates keys, signs and verifies (same page). This is independent of OS code signing, so it works with no paid certificate.
- Install: a helper process backs up the app, swaps the binary, relaunches, and rolls back on a failed launch. Supports single binaries, `.zip` (recommended for macOS `.app`), and `.tar.gz` with one top-level entry. `.dmg`, `.pkg` and `.msi` are not supported in v1 (same page). The GitHub provider excludes Windows installer assets by default (secondhand, from release notes).
- Includes a default update window with Markdown release notes, skip-version, remind-later and periodic checks (same page).
- Status: Wails v3 is at v3.0.0-beta.28 (pre-release, "API is stable"). beta.21 fixed an updater relaunch bug: https://github.com/wailsapp/wails/releases. The guide is on docs-preview, so the published docs may differ.
- Unknown: whether a self-updated unsigned macOS app gets re-quarantined after the swap. The guide only mentions that Gatekeeper and code signing need consideration (**unverified**). Test on Sequoia.

## GitHub Releases as the update host

- Up to 1000 assets per release, each file under 2 GiB, and "no limit on the total size of a release, nor bandwidth usage": https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases
- Update checks call the GitHub API. Unauthenticated rate limits for many clients are not researched here (**unverified** risk).

## Cost summary

| Route | Cost | macOS | Windows |
|---|---|---|---|
| Unsigned + GitHub Releases + Wails updater | $0 | Settings > Privacy & Security allow, or xattr | "More info" > "Run anyway" on every new version; Smart App Control may block |
| + SignPath (OSS) | $0 | no change | Signed; warnings until reputation builds |
| + Artifact Signing | $9.99/mo (US/CA individuals; listed-country orgs) | no change | same as above |
| + Apple Developer Program | $99/yr | Notarized: no warning; Homebrew cask possible | — |
