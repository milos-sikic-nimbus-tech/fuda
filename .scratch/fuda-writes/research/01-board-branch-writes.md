# Board-branch writes over the GitHub and Azure DevOps REST APIs

Research date: 2026-10-08. Sources are official docs, read on that date; the URLs are listed at the end.
The fuda code referenced here: `api/internal/source/github/github.go` (REST, `X-GitHub-Api-Version: 2022-11-28`)
and `api/internal/source/azure/azure.go` (`api-version=7.1`).

## Direct answer

- **GitHub:** a commit that changes N Task files takes **one write call** with the GraphQL mutation
  `createCommitOnBranch`, or **three write calls** with the REST Git Database API: `POST git/trees`
  with inline `content` and a `base_tree`, then `POST git/commits`, then `PATCH git/refs/heads/{branch}`.
  Both need the head commit first. fuda's `Head()` already returns it, and the REST path also needs
  that commit's tree SHA (one `GET git/commits/{sha}`).
  - Stale-parent detection: GraphQL takes `expectedHeadOid`, a real compare-and-swap.
  - REST `PATCH ref` with `force: false` accepts only a fast-forward. Because fuda's commit has the
    head it read as its only parent, that refuses any concurrent advance. It is still not a strict
    CAS (see GitHub details).
- **Azure DevOps:** **one write call**, `POST _apis/git/repositories/{repo}/pushes`. It carries
  `refUpdates[{name, oldObjectId}]` and one commit whose `changes` hold N `edit` entries with
  `rawtext` content. `oldObjectId` is the CAS: the ref-update status for a mismatch is
  `staleOldObjectId` ("the caller lost a race to update the ref").
- **Branches:** a non-default branch works on both hosts.
  - GitHub REST refs only need `refs/` plus at least two slashes, so `refs/fuda/*` probably works
    over REST (unverified). GraphQL `createCommitOnBranch` is **branches only** (`refs/heads/`).
  - On Azure, custom `refs/fuda/*` refs are undocumented (unverified).
- **Token permissions:**
  - GitHub fine-grained token: **Contents: write** (plus Metadata: read, which is implicit).
    Classic OAuth or PAT: `repo`, or `public_repo` for public repos.
  - Azure: scope **`vso.code_write`** (Code read & write), and the user needs the **Contribute**
    permission. **Create branch** is needed for new branches, and **Bypass policies when pushing**
    is needed if the branch has required policies.
- **Load:** one commit every few seconds across up to 20 users (roughly 10-30 commits per minute in
  total, about 1-2 per minute per user) is **well within the documented limits on both hosts**.
  - GitHub's limits count per user. With per-user tokens, the binding limit is the content-creation
    cap of 80 per minute and 500 per hour.
  - One user who writes non-stop every 3 s through REST (3 writes per commit) would reach the
    500-per-hour cap after about 8 minutes. Through GraphQL (1 write per commit) it would take about
    25 minutes.
  - The real risk is not quota. It is **CAS conflicts when 20 people commit to the same branch**.
    fuda should serialize writes per branch inside the server, and retry on conflict (re-read the
    head, re-apply the frontmatter edit, re-commit).

## GitHub

### Option A: GraphQL `createCommitOnBranch` (recommended)

One `POST https://api.github.com/graphql` with input:

- `branch`: either `{repositoryNameWithOwner, branchName}` or a branch node ID. The docs say it "must
  refer to a git branch (i.e. the fully qualified path must begin with refs/heads/, although including
  this prefix is optional)".
- `expectedHeadOid` (required): "The git commit oid expected at the head of the branch prior to the
  commit."
- `fileChanges.additions[{path, contents}]` for each changed Task file (`contents` is base64), and
  `fileChanges.deletions[]` if needed.
- `message`.

Properties:

- **Calls per commit:** 1 write. fuda already has the head SHA from `Head()`, the same value the board
  snapshot was built from.
- **Stale parent:** a mismatch on `expectedHeadOid` fails the mutation. The exact error text is not in
  the reference (unverified), so fuda should treat any mutation error after a re-read head
  differs as a conflict.
- **Author:** "authored by the owner of the credential which authenticates the API request". The
  author cannot be set ("will not add support for this in the future").
- **Signing:** "automatically signed by GitHub if supported and will be marked as verified". This is
  useful if the target branch requires signed commits.
- **Limitation:** branches only, so it cannot write to `refs/fuda/*`.
- **Cost:** a mutation costs 5 secondary points. The GraphQL secondary budget is 2,000 points per
  minute. It also counts once against the content-creation cap.

### Option B: REST Git Database API

1. `GET /repos/{o}/{r}/git/ref/heads/{branch}` returns the head SHA H. fuda already does this in
   `Head()`.
2. `GET /repos/{o}/{r}/git/commits/{H}` returns `tree.sha`.
3. `POST /repos/{o}/{r}/git/trees` with `{base_tree: <tree sha>, tree: [{path, mode: "100644",
   type: "blob", content: "<new file text>"}, ...]}`. Inline `content` makes GitHub write the blobs,
   so no separate `POST git/blobs` call is needed. Without `base_tree`, any file not listed shows as
   deleted.
4. `POST /repos/{o}/{r}/git/commits` with `{message, tree, parents: [H]}`. The author defaults to "the
   authenticated user and the current date". The commit is unsigned unless fuda passes a `signature`.
5. `PATCH /repos/{o}/{r}/git/refs/heads/{branch}` with `{sha: <new commit>, force: false}`.

Properties:

- **Calls per commit:** 3 writes (steps 3-5) and 1-2 reads.
- **Stale parent:** `force` "Indicates whether to force the update or to make sure the update is a
  fast-forward update". With `false` the docs say "you're not overwriting work". If the branch moved
  to H′ after H was read, the new commit (parent H) is not a descendant of H′, so the PATCH is
  rejected.
  - The documented codes are 409 Conflict and 422 "Validation failed, or the endpoint has been
    spammed". The non-fast-forward case is commonly reported as 422 "Update is not a fast forward",
    but the reference does not spell it out (unverified).
  - Caveat: this check is fast-forward-only, not a strict CAS. If someone *force-moved the branch
    backwards* to an ancestor of H, the PATCH would be accepted. Re-reading the head just before the
    PATCH narrows that window but does not close it. GraphQL `expectedHeadOid` has no such gap.
- **Orphans:** objects from a failed attempt (tree, commit) are left unreferenced. GitHub cleans them
  up; this is harmless.
- **Single-file alternative:** `PUT /repos/{o}/{r}/contents/{path}` with the file's blob `sha` is
  also CAS-guarded, but it writes **one file per commit**. It does not meet the "one commit, several
  files" requirement.

### Non-default branch and custom refs

- **Non-default branch:** works in both options. Name the branch in `branch` (GraphQL) or in the ref
  path (REST).
- **Custom refs (`refs/fuda/*`):**
  - `POST git/refs` rejects a `ref` unless it starts with `refs` and has at least two slashes, so
    `refs/fuda/board` passes that rule. The list endpoint returns "Anything in the namespace … not just
    heads and tags".
  - Creating and PATCHing a custom ref is therefore *likely* to work over REST. **No doc confirms it;
    probe on a scratch repo.**
  - `createCommitOnBranch` cannot target a custom ref.
  - "You are unable to create new references for empty repositories".
- **Rulesets** target branches and tags (plus push rules that apply to every push). A custom
  `refs/fuda/*` ref is outside branch rulesets. A protected board branch that requires PRs or status
  checks would reject direct API commits; the exact error is unverified.

### Permissions (user's token)

- **Fine-grained PAT / GitHub App user token:** **Contents: write**. The docs list it for
  `POST git/trees`, `POST git/commits`, `POST git/refs`, `PATCH git/refs/{ref}` and
  `PUT contents/{path}`. Contents: read is listed for `GET git/ref/{ref}`.
  - `POST/PATCH git/refs` and `PUT contents` carry an "Additional permissions" mark. They also appear
    under **Workflows**, which matters only when touching `.github/workflows/*`. Task files do not.
- **Classic OAuth app / PAT:** `repo` ("full access to public and private repositories including read
  and write access to code"), or `public_repo` for public repositories only.
- GraphQL `createCommitOnBranch` needs the same repository write access (assumed from the Contents
  model; the GraphQL reference does not list permissions per mutation).

### Rate limits (github.com)

**Primary limit:** 5,000 requests per hour per user, shared across that user's PATs and the apps
acting for them. It is 15,000 for apps owned by an Enterprise Cloud org. GraphQL is 5,000 points per
hour per user.

**Secondary limits** "are subject to change without notice":

- No more than 100 concurrent requests, shared between REST and GraphQL.
- No more than 900 points per minute for REST. Most GETs cost 1 point; most POST, PATCH, PUT and
  DELETE calls cost 5.
- 2,000 points per minute for GraphQL. A mutation costs 5, a query 1.
- No more than 90 s of CPU time per 60 s.
- **Content creation:** "no more than 80 content-generating requests per minute and no more than 500
  content-generating requests per hour". "Some endpoints have lower content creation limits", and the
  docs do not name them. The cap "include[s] actions taken on the GitHub web interface as well as via
  the REST API and GraphQL API".
  - The docs do not say whether each Git Database POST or PATCH counts as one content-generating
    request (unverified). The worst case assumed here: 3 per REST commit, 1 per GraphQL commit.

**Handling:**

- A limit returns 403 or 429.
- Honour `retry-after`. If `x-ratelimit-remaining` is 0, wait until `x-ratelimit-reset`. Otherwise
  wait at least a minute, then back off exponentially.
- The best-practices guide says to "make requests serially instead of concurrently" and, for many
  writes, "wait at least one second between each request".

**Arithmetic for fuda:**

- Secondary limits are tracked per user or token. The docs don't say which; per user is assumed.
  With each user committing through their own token, 20 users do not share a bucket.
- At about 1-2 commits per minute per user: REST is ≤6 writes per minute (≤30 points), and GraphQL is
  ≤2 mutations per minute. Both are far under 80 per minute, 500 per hour and 900 points per minute.
- A single user sustaining one commit every 3 s:
  - REST: 60 writes per minute is under 80 per minute, but reaches 500 per hour in about 8 minutes.
  - GraphQL: 20 per minute reaches 500 per hour in about 25 minutes.
  - People dragging cards do not sustain this. If fuda batches rapid edits from one user (for
    example, debounce for 2-3 s into one commit), the issue goes away.
- If fuda ever wrote through **one shared server token**, the 20 users would share one 80-per-minute
  and 500-per-hour bucket. 500 per hour is about 8 commits per minute (GraphQL), which "a commit every
  few seconds" would exceed. Per-user tokens are the right model.

## Azure DevOps

### Pushes - Create (one call)

```http
POST https://dev.azure.com/{org}/{project}/_apis/git/repositories/{repo}/pushes?api-version=7.1
{
  "refUpdates": [{ "name": "refs/heads/{branch}", "oldObjectId": "<head sha fuda read>" }],
  "commits": [{
    "comment": "fuda: move T-12 to Done",
    "changes": [
      { "changeType": "edit", "item": { "path": "/docs/board/tasks/T-12.md" },
        "newContent": { "content": "<full new file text>", "contentType": "rawtext" } },
      { "changeType": "edit", "item": { "path": "/docs/board/tasks/T-13.md" },
        "newContent": { "content": "...", "contentType": "rawtext" } }
    ]
  }]
}
```

- **Calls per commit:** 1 write. fuda's `Head()` (`GET refs?filter=heads/{branch}`) supplies
  `oldObjectId`. The "Multiple changes" example in the docs puts several changes in one commit. The
  docs list 200 as the response code; the examples show 201.
- **Content:** `contentType` is `rawtext` or `base64encoded`. Each change sends the whole new file;
  there is no patch format.
- **Stale parent:** `oldObjectId` must equal the current ref value. Refs - Update Refs documents the
  status `staleOldObjectId`: "the old object ID presented in the request was not the object ID of the
  ref when the database attempted the update. The most likely scenario is that the caller lost a race
  to update the ref."
  - For a *push*, the failure surfaces as an HTTP error, commonly `TF401028` ("reference has already
    been updated by another client"), reportedly HTTP 409. **The exact code and status are not
    documented on the Pushes page (unverified).**
  - Unlike GitHub's REST path, this is a real compare-and-swap.
- **Other relevant ref statuses:** `rejectedByPolicy` (branch policies), `forcePushRequired`,
  `writePermissionRequired`, `createBranchPermissionRequired`, `invalidRefName`, `locked`.
- **Author:** the authenticated identity, shown as `pushedBy` on the push. The docs examples show
  `author` and `committer` in the response; whether request-supplied author fields are honoured is
  unverified.
- **New branch:** a `refUpdates` name for a branch that does not exist, with `oldObjectId` set to an
  existing commit, creates it ("Update a file in a new branch"). All-zeros is used for an initial
  commit.

### Non-default branch and custom refs

- **Non-default branch:** works. Set `refUpdates[].name` to any `refs/heads/...`.
- **Custom refs (`refs/fuda/*`):** Refs - Update Refs takes any `name`, and the permission model
  names branches, tags and notes. Nothing documents arbitrary namespaces. `invalidRefName` exists as
  a status. **Unverified; probe before relying on it.**
- **Branch policies:** if the board branch has required policies (for example, require PR), direct
  pushes fail with `rejectedByPolicy` unless the user has **Bypass policies when pushing**.

### Permissions

- **Token scope:** `vso.code_write` ("read, update, and delete source code … branches"), shown in the
  PAT UI as Code (Read & write). The same scope applies to a Microsoft Entra OAuth token.
- **Repository permissions:**
  - **Contribute**: "can push their changes to existing branches … Doesn't override restrictions in
    place from branch policies".
  - **Create branch**, if fuda creates the board branch.
  - **Bypass policies when pushing**, for a policy-protected branch.
  - The Contributors group has Contribute, Create branch, Create tag and Manage notes by default.
- **Note:** global PATs "stop working on December 1, 2026". Use organization-scoped PATs or Entra
  OAuth. Microsoft recommends Entra tokens over PATs for applications.

### Rate limits

- **Global consumption limit:** 200 TSTUs within any sliding five-minute window, **per user** (and
  per pipeline). One TSTU is "the average load generated by a typical Azure DevOps user over five
  minutes". Normal activity spikes at 10 TSTUs or fewer.
- **Over the limit:** requests are delayed by milliseconds up to 30 s, then blocked with 429
  (`TF400733`).
- **Headers:** `Retry-After`, `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`,
  `X-RateLimit-Delay`, and `X-RateLimit-Cost` (the TSTUs this request consumed).
- **Git limits** cover size only: pushes up to 5 GB, files up to 100 MB, path length. There is no
  documented push-frequency limit.
- **Per-push cost:** not published ("You can't calculate usage in TSTUs for an action with a
  formula"). A small push of a few KB of markdown is ordinary Git work. 1-2 pushes per minute per
  user should be nowhere near 200× a typical user (unverified). Log `X-RateLimit-Cost` in the first
  build to confirm.

## Contention and recommended shape (inference, not from docs)

- **All users commit to the same branch.** When two writes race, one loses the CAS. A single fuda
  process can remove most races by **serializing writes per branch** with a mutex or queue. That also
  follows GitHub's "requests serially" advice. Only pushes from outside fuda can then conflict.
- **On conflict:**
  1. Re-read the head.
  2. Fetch the affected Task files at the new head (`FileAt`).
  3. Re-apply the frontmatter change.
  4. If the same field was changed meanwhile, surface it to the user rather than overwrite.
  5. Retry a bounded number of times.
- **Batching:** coalescing a user's rapid edits into one commit keeps GitHub's hourly
  content-creation cap out of reach.
- **Sync:** after a successful write, fuda knows the new head SHA and can refresh its snapshot without
  waiting for the next sync.

## Unverified (needs a probe against a scratch repo)

1. **GitHub custom refs:** whether `POST git/refs` and `PATCH git/refs/fuda/...` accept a
   `refs/fuda/*` ref.
2. **GitHub error bodies:** the exact status and message for a non-fast-forward `PATCH ref` (422
   "Update is not a fast forward" expected), and for an `expectedHeadOid` mismatch.
3. **GitHub content creation:** which Git Database calls count as "content-generating", and whether
   secondary limits are tracked per user or per token.
4. **Azure push conflict:** the HTTP status and message for a stale `oldObjectId` on Pushes - Create
   (409 / TF401028 expected).
5. **Azure custom refs:** whether `refs/fuda/*` refs are accepted, and which permission governs them.
6. **Azure push cost:** the TSTU cost of one small push (read `X-RateLimit-Cost`).
7. **Azure author fields:** whether `author` and `committer` in the push request body are honoured.

## Sources

- GitHub Git references: https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28
- GitHub Git trees: https://docs.github.com/en/rest/git/trees?apiVersion=2022-11-28
- GitHub Git commits: https://docs.github.com/en/rest/git/commits?apiVersion=2022-11-28
- GitHub GraphQL `createCommitOnBranch` and `CreateCommitOnBranchInput`: https://docs.github.com/en/graphql/reference/commits
- GitHub REST rate limits: https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28
- GitHub REST best practices: https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api?apiVersion=2022-11-28
- GitHub GraphQL rate limits: https://docs.github.com/en/graphql/overview/rate-limits-and-query-limits-for-the-graphql-api
- GitHub fine-grained PAT permissions: https://docs.github.com/en/rest/authentication/permissions-required-for-fine-grained-personal-access-tokens?apiVersion=2022-11-28
- GitHub OAuth scopes: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps
- GitHub rulesets: https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets
- Azure Pushes - Create (7.1): https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pushes/create?view=azure-devops-rest-7.1
- Azure Refs - Update Refs (7.1): https://learn.microsoft.com/en-us/rest/api/azure/devops/git/refs/update-refs?view=azure-devops-rest-7.1
- Azure rate and usage limits: https://learn.microsoft.com/en-us/azure/devops/integrate/concepts/rate-limits?view=azure-devops
- Azure Git limits: https://learn.microsoft.com/en-us/azure/devops/repos/git/limits?view=azure-devops
- Azure Git repository permissions: https://learn.microsoft.com/en-us/azure/devops/repos/git/set-git-repository-permissions?view=azure-devops
- Azure permissions reference: https://learn.microsoft.com/en-us/azure/devops/organizations/security/permissions?view=azure-devops
- Azure PATs: https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops
