# 03 — Login and tokens: how a user signs in so fuda can write as them

Researched 2026-10-08 against current GitHub Docs and Microsoft Learn. Every claim below links to
the page it came from. Things I could not confirm in a primary source are listed at the end.

## Direct answer

**GitHub: use a GitHub App and its *user access tokens*, not an OAuth App.**

- One fuda maintainer registers one GitHub App, free, under their account or the org that owns the
  task repos. Give it **Contents: Read and write**, turn on
  "Expire user authorization tokens", and turn on "Enable Device Flow" only if the desktop build
  uses it.
- The repo owner or org owner **installs** the app on just the task repos. Each user then
  **authorizes** it once. The token can do only what the app *and* the user can both do, and only
  in repos where the app is installed. Commits it makes are attributed to the user.
- **Web server:** web flow (authorization code + PKCE). The client secret stays on the server, the
  browser gets a cookie, and the server refreshes the token.
- **Desktop (Wails):** GitHub's own guidance prefers the authorization code + PKCE flow through a
  loopback redirect over the device flow. That means shipping the client secret inside the binary,
  which GitHub says you cannot keep secret anyway. The device flow needs no secret, and neither does
  refreshing a device-flow token, but GitHub warns it can be used for phishing. For non-developer
  users either flow works. Store the refresh token in the OS keychain.
- **Lifetimes:** the user token lasts **8 hours**. The refresh token lasts **6 months**, and each
  refresh issues a new pair. A user who doesn't open fuda for 6 months signs in again.
- **Smallest scope for "write one branch":** neither GitHub nor Azure has a token scope limited to
  one branch. The smallest GitHub grant is **Contents: write on only the installed repos**. An OAuth
  App can't get that close: to write a private repo it needs `repo`, which covers all of the user's
  repos. Limiting writes to one branch on the server side (protected branches or rulesets) is
  **not available for private repos on GitHub Free**, so at $0 **fuda itself must enforce the
  branch**.

**Azure DevOps: use Microsoft Entra ID OAuth if the Azure DevOps org is connected to Entra.
Otherwise, use a per-user PAT.**

- The old Azure DevOps OAuth apps have accepted **no new registrations since April 2025**, and they
  stop working in 2026. Entra ID OAuth is the only OAuth option left.
- Someone with rights in the org's Entra tenant registers **one free app registration**. Add a
  delegated permission on the **Azure DevOps** API (resource `499b84ac-1321-427f-aa17-267ca6975798`),
  and request **`vso.code_write`** rather than `user_impersonation`, which grants full access.
- **Web server:** a confidential client using the authorization code flow with a client secret, and a
  "Web" redirect URI.
- **Desktop:** a public client through MSAL Go: interactive sign-in with a loopback redirect on
  `127.0.0.1`, where any port matches, or the device code flow. Public client flows must be enabled.
  Cache the MSAL token cache in the keychain.
- **Lifetimes:** access tokens are short-lived; the docs say to use the returned expiry, not a fixed
  number. Refresh tokens last **90 days** and are replaced on every use.
- **Hard limit:** Entra apps "don't natively support Microsoft account (MSA) users" for Azure
  DevOps. If the org signs in with personal Microsoft accounts, OAuth is unavailable and the only
  route is a **PAT**: organization-scoped, scope **Code (Read & write)**, with an expiry the user
  picks within the org's policy. Global PATs stop working on 1 Dec 2026.
- Azure DevOps **branch permissions** can be set per branch at no extra cost, so an admin *can*
  limit "Contribute" to the board branch.

**Cost:** registering the GitHub App, the Entra app registration, and Entra ID Free all cost
nothing. **The one real cost risk is Azure DevOps licensing.** Writing code in a private project
needs **Basic** access. Stakeholder access, which is free, doesn't include Code. Basic is **free for
the first 5 users only**, so writers 6–20 cost money unless they hold a Visual Studio or GitHub
Enterprise licence. GitHub Free has no comparable per-user charge for collaborators.

## GitHub details

### OAuth App versus GitHub App

| | OAuth App | GitHub App (user access token) |
|---|---|---|
| Permission model | Broad scopes. Writing private code needs `repo` ("full access to public and private repositories") [scopes] | Fine-grained: Contents read/write, and only on the repos where it's installed [diff], [permreq] |
| What the token can reach | Everything the user can reach | The overlap of app permissions, user access, and installed accounts [behalf] |
| Token expiry | Long-lived by default. Expiring tokens are opt-in, per app or per sign-in with `offline_access` [oauth] | Expiring tokens are on by default, and GitHub "strongly recommends" leaving them on [register] |
| Org control | With OAuth app access restrictions on, only an owner can approve [diff] | An org owner or repo admin installs it on chosen repos [diff] |
| GitHub's advice | — | "In general, GitHub Apps are preferred over OAuth apps." [diff] |

Attribution: "API requests made by an app on behalf of a user will be attributed to that user." The
audit log shows the access type as "GitHub App user-to-server token" [behalf].

### Permission needed to commit

`PUT /repos/{owner}/{repo}/contents/{path}`, `POST git/blobs`, `git/trees`, `git/commits`, and
`PATCH git/refs/{ref}` all need **Contents: write** [permreq]. The contents endpoint and the refs
endpoint are also listed under **Workflows**. That fits the usual requirement for changing
`.github/workflows/*`, which fuda never needs to write.

### Web flow (server)

1. Redirect to `https://github.com/login/oauth/authorize` with `client_id`, `state`, `redirect_uri`
   and PKCE `code_challenge` (only `S256` is accepted) [gen].
2. Exchange the `code` at `https://github.com/login/oauth/access_token` with `client_id`,
   `client_secret`, `code` and `code_verifier` [gen].
3. You get an `access_token` (`ghu_` prefix) and a `refresh_token` (`ghr_` prefix) [gen].
4. A GitHub App takes up to 10 callback URLs [register], enough for a few deployments plus dev.

### Device flow (desktop or headless)

- Turn on "Enable Device Flow" in the app settings first. Otherwise the error is
  `device_flow_disabled` [oauth].
- `POST https://github.com/login/device/code` with `client_id`. The user enters the `user_code` at
  `https://github.com/login/device`. Poll `access_token` with `grant_type`
  `urn:ietf:params:oauth:grant-type:device_code`, respecting `interval`; `slow_down` adds 5 s
  [gen], [oauth]. Codes expire after 900 s [oauth].
- "The `client_secret` is not needed for the device flow" [oauth]. Refresh also needs no secret:
  the docs mark it "Required unless the user access token was generated using the device flow"
  [refresh].

### Which flow for the Wails app

GitHub's best-practices page says, for native apps: "you cannot secure your client secret" and "You
do have to ship the client secret in the application's code". It recommends PKCE and says "It is
preferable to use the authorization code with PKCE over the device flow". It also advises "do not
enable the device flow for your application unless you are using the app in a constrained
environment", because of phishing [best]. The OAuth docs allow loopback redirects such as
`http://127.0.0.1:1234/path` against a registered `http://127.0.0.1/path`, with any port [oauth].
For tokens in native apps: "store tokens via the mechanism recommended for your app's platform",
which for fuda means the OS keychain [best].

Practical reading for fuda: in the desktop build, prefer the loopback + PKCE flow, accept that the
embedded secret is public, and keep the device flow as a fallback only if loopback turns out
unreliable.

### Lifetimes and refresh

- User access token: 8 hours (`expires_in` 28800). Refresh token: 6 months
  (`refresh_token_expires_in` 15897600) [gen], [refresh].
- Refreshing returns a new pair. For OAuth apps: "Once you use a refresh token, that refresh token
  and the old access token will no longer work" [oauth]. So the server must store each new refresh
  token atomically, or concurrent refreshes will lock the user out.
- An expired refresh token means the user goes through the web flow or device flow again [refresh].

### Limiting writes to one branch

- Neither scopes nor GitHub App permissions mention branches. Both work per repo [scopes], [permreq].
- Branch restrictions ("only users, teams, or apps that have been given permission can push") work
  in public repos of Free orgs, and "in all repositories owned by an organization using GitHub Team
  or GitHub Enterprise Cloud" [protected]. GitHub's plans page lists protected branches for private
  repos only under Pro and Team [plans], and the rulesets page names Team and Enterprise [rulesets].
- Result: for **private repos at $0**, the only branch guard is fuda's own code, such as always
  targeting the configured branch and refusing other refs. Anyone with write access can still push
  any branch with plain git.

### Who registers what (GitHub)

| What | Who | Cost |
|---|---|---|
| GitHub App registration (client ID, client secret, private key not needed for user tokens) | A personal account, an org owner, or an app manager. Up to 100 apps per account [register] | No charge found. See Unverified |
| Install on the task repos | The org owner, or a repo admin on a user-owned repo [diff] | — |
| Authorize the app | Each user, once [behalf] | — |

## Azure DevOps details

### Options

| | Microsoft Entra ID OAuth | Azure DevOps OAuth (legacy) | PAT |
|---|---|---|---|
| Status | Recommended [oauth-ado] | "New app registrations are no longer accepted as of April 2025. The service is scheduled for full deprecation in 2026." [oauth-ado] | Supported, but "Avoid PATs when a more secure authentication method is available" [pat] |
| Users with personal Microsoft accounts (MSA) | **Not supported** for the Azure DevOps resource [entra-oauth] | Supported, but closed to new apps | Supported |
| Expiry | Access token short-lived, refresh token 90 days and replaced on each use [rt] | — | The user picks it, capped by org policy [pat] |
| Scope | Azure DevOps scopes such as `vso.code_write`. "Both Microsoft Entra ID OAuth and Azure DevOps OAuth use the same scope definitions." [oauth-ado] | same | Same scope list, chosen in the UI [pat] |

### Scopes

- `vso.code_write`, Code (read and write): "read, update, and delete source code … create and manage
  pull requests" [oauth-ado]. This is the smallest scope that can push a commit.
- `user_impersonation`: "Grants full access … Request or consent this scope with caution" [oauth-ado].
  Avoid it.
- With Entra, request `499b84ac-1321-427f-aa17-267ca6975798/vso.code_write`, or `/.default` for all
  the permissions configured on the registration [entra-oauth].

### Desktop (Wails)

- Use a public client in MSAL Go: `public.New(clientID, public.WithAuthority(...))`, then
  `AcquireTokenSilent` and, if that fails, `AcquireTokenInteractive` [msal-go]. The token cache
  persists through a cache accessor, which fuda can back with the OS keychain [msal-go].
- Redirect: register `http://localhost` (or `127.0.0.1`) under "Mobile and desktop applications".
  For loopback URIs "the port component … is ignored", and Microsoft says to "Prefer 127.0.0.1 over
  localhost". An `http://127.0.0.1` URI has to be added through the app manifest, because the
  portal box rejects it [reply].
- Device code flow: needs no redirect URI [reply]. It does need a tenant-specific authority:
  "Device Code flow is not supported under /common or /consumers endpoint" [devcode].

### Web server

- Use a confidential client: a "Web" platform redirect URI on https [reply], the authorization code
  flow, and a client secret stored on the server. The legacy Azure DevOps OAuth FAQ says the same:
  the web server flow "requires securely storing an app secret" [oauth-ado].
- Refresh tokens for confidential clients survive password changes [rt]. The server must store each
  new refresh token, because "Refresh tokens replace themselves with a fresh token upon every use"
  [rt].

### Consent

"By default, all users are allowed to consent to applications for permissions that don't require
administrator consent." With the stricter built-in policy, users can still consent to "apps that
are registered in your tenant", but only for permissions classified low impact [consent]. So
whether non-admin users can consent to `vso.code_write` without an admin depends on the tenant. See
Unverified.

### PAT route (MSA-backed orgs, or if Entra registration isn't possible)

- Create one under User settings > Personal access tokens. Pick **one organization**, an expiry, and
  **Code (Read & write)** [pat].
- Send it as HTTP Basic with an empty username: `base64(":" + PAT)` [pat].
- "Azure DevOps is retiring global PATs. All existing global PATs stop working on December 1, 2026."
  [pat]
- Admins can restrict PAT creation, scopes, and maximum lifetime. In Entra-backed orgs a PAT goes
  inactive if the user doesn't sign in within 90 days [pat].
- The user pastes the PAT once. fuda keeps it server-side or in the keychain. There is no refresh;
  the user makes a new one when it expires.

### Branch limits (Azure)

Branch security can set **Contribute** per user or group on a single branch. Branch creators get
Contribute, Force push and Manage permissions on branches they create [branch]. So an admin can
deny Contribute everywhere except the board branch, at no extra licence cost.

### Who registers what (Azure DevOps)

| What | Who | Cost |
|---|---|---|
| Entra app registration (desktop as a public client, server as a confidential client: one registration with both platforms, or two) | A user allowed to register apps in the tenant the Azure DevOps org is connected to | Entra ID Free "has no product charges" [entrafree] |
| Optional admin consent for `vso.code_write` | A tenant admin, if the consent policy requires it [consent] | — |
| Basic access for each writer | An Azure DevOps org admin | Free for the first 5 Basic users, then paid [billing] |
| PAT (fallback) | Each user | — |

## Cost summary

- GitHub App or OAuth App registration: no charge found in the docs (see Unverified). GitHub Free
  orgs allow private repos.
- GitHub branch protection for private repos needs Pro or Team, so **skip it at $0** and enforce
  the branch in fuda [plans], [protected].
- Entra app registration: no charge (Entra ID Free) [entrafree].
- Azure DevOps: Stakeholder is free but has no Code access. **Basic is free for 5 users**, then paid
  per user. Visual Studio subscribers and GitHub Enterprise users count as Basic at no extra charge
  [access], [billing]. With 20 users where everyone writes, expect up to 15 paid Basic seats unless
  only a few people write.

## Sources

- [gen] Generating a user access token for a GitHub App — https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-user-access-token-for-a-github-app
- [refresh] Refreshing user access tokens — https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/refreshing-user-access-tokens
- [behalf] Authenticating with a GitHub App on behalf of a user — https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-with-a-github-app-on-behalf-of-a-user
- [best] Best practices for creating a GitHub App — https://docs.github.com/en/apps/creating-github-apps/about-creating-github-apps/best-practices-for-creating-a-github-app
- [register] Registering a GitHub App — https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app
- [diff] Differences between GitHub Apps and OAuth apps — https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/differences-between-github-apps-and-oauth-apps
- [oauth] Authorizing OAuth apps (device flow, expiring tokens, `offline_access`, loopback redirects) — https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
- [scopes] Scopes for OAuth apps — https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps
- [permreq] Permissions required for GitHub Apps — https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps
- [protected] About protected branches — https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches
- [rulesets] About rulesets — https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets
- [plans] GitHub's plans — https://docs.github.com/en/get-started/learning-about-github/githubs-plans
- [oauth-ado] OAuth 2.0 authentication for Azure DevOps REST APIs (deprecation and scope table) — https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/oauth
- [entra-oauth] Build Azure DevOps integrations with Microsoft Entra OAuth apps — https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/entra-oauth
- [entra] Authenticate to Azure DevOps with Microsoft Entra ID — https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/entra
- [pat] Use personal access tokens — https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate
- [access] About access levels — https://learn.microsoft.com/en-us/azure/devops/organizations/security/access-levels
- [billing] Manage paid access for users — https://learn.microsoft.com/en-us/azure/devops/organizations/billing/buy-basic-access-add-users
- [branch] Set Git branch security and permissions — https://learn.microsoft.com/en-us/azure/devops/repos/git/branch-permissions
- [rt] Refresh tokens in the Microsoft identity platform — https://learn.microsoft.com/en-us/entra/identity-platform/refresh-tokens
- [devcode] Device code flow (desktop) — https://learn.microsoft.com/en-us/entra/identity-platform/scenario-desktop-acquire-token-device-code-flow
- [reply] Redirect URI best practices and limitations — https://learn.microsoft.com/en-us/entra/identity-platform/reply-url
- [consent] Configure how users consent to applications — https://learn.microsoft.com/en-us/entra/identity/enterprise-apps/configure-user-consent
- [entrafree] Microsoft Entra ID Free — https://learn.microsoft.com/en-us/azure/cost-management-billing/manage/microsoft-entra-id-free
- [msal-go] MSAL for Go (README and wiki, via Context7 `/azuread/microsoft-authentication-library-for-go`) — https://github.com/AzureAD/microsoft-authentication-library-for-go

## Unverified or open

- **Loopback redirects for GitHub *Apps*.** The rule that any port matches `http://127.0.0.1/path`
  is documented on the OAuth-app page [oauth]. I didn't find the same statement for GitHub Apps.
  Test it with a real app before relying on it.
- **Explicit "free" statement for GitHub App or OAuth App registration.** No GitHub page says it
  costs nothing. The registration docs list no charge and no plan requirement [register].
- **Whether users can self-consent to `vso.code_write`.** The docs don't say whether it's marked as
  needing admin consent. It's flagged "High privilege" in the scope table [oauth-ado]. The user
  consent policy is tenant-specific [consent]. Check in the target tenant.
- **Price of an extra Azure DevOps Basic user.** The docs say only "paid for 6 or more users"
  [billing]. The amount is on the Azure pricing page, which I didn't check.
- **Exact date in 2026 when legacy Azure DevOps OAuth stops working.** The docs say only "2026"
  [oauth-ado]. It doesn't matter for fuda, which can't register a legacy app anyway.
- **MSA support in Entra OAuth for Azure DevOps.** "Microsoft is currently working on native
  support" [entra-oauth]. Recheck before deciding to rely on PATs for MSA-backed orgs.
- **Entra client secret maximum lifetime** (the server's confidential client). Not checked. The
  secret expires and has to be rotated by whoever registered the app.
- **GitHub org OAuth app access restrictions default.** Mentioned only in the differences page
  [diff]. The default state for new orgs wasn't checked. It doesn't affect GitHub Apps.
