# Research 02: cost of polling a branch head every 5 seconds

Date: 2026-10-08. Load assumed: up to 20 open apps (web or Wails desktop), each checking a branch head
every 5 s, so 12 checks/min and 720/h per app.

## Direct answer

1. **GitHub: yes, authenticated 304s are free on the primary rate limit.** The docs say so, and I
   checked it live: three authenticated `If-None-Match` requests that returned 304 left
   `X-Ratelimit-Used` unchanged. Unauthenticated 304s **do** count: each one used one of the 60/h.
2. **Cheapest GitHub call:** `GET /repos/{owner}/{repo}/commits/{branch}` with
   `Accept: application/vnd.github.sha` and `If-None-Match: "<last known sha>"`. The 200 body is just
   the 40-byte SHA, and the ETag is the SHA itself in quotes. So fuda can send the SHA it already
   synced, with no first call to fetch an ETag. A 304 means the head has not moved, and costs nothing
   on the primary limit. Unlike `git/ref/heads/{branch}`, it does not send `X-Poll-Interval: 300`.
3. **Cheapest Azure DevOps call:**
   `GET .../_apis/git/repositories/{repo}/refs?filter=heads/{branch}&api-version=7.1`, read
   `value[0].objectId`. There is **no ETag or 304** on it or on any other branch-head call I tried:
   every response had `cache-control: no-cache, no-store` and no `ETag`. Every poll is a full 200.
4. **Azure TSTUs at this load are not a concern.** The limit is 200 TSTUs per identity in a sliding
   5-minute window. Measured `X-RateLimit-Cost` was about 0.006 to 0.014 TSTU per call (the refs call
   sent no cost header at all). A 5 s poll makes 60 calls per 5 minutes, about 0.4 to 0.9 TSTU, under
   0.5 % of the limit. Even 20 apps sharing one identity (about 7 to 17 TSTU) stay under 10 %.
5. **Free push signal without a public server: none that is fit for production.** GitHub webhooks
   need a public HTTP receiver, and `gh webhook forward` is for testing only (one user per repo). The
   Events API is polled and can lag 30 s to 6 h. Azure DevOps service hooks can write to an **Azure
   Storage queue or Service Bus**, which a desktop app can read without a public server. That needs
   an Azure subscription (cheap, not free) and a project admin to create the subscription. Polling as
   in points 2 and 3 is the practical answer.

## Recommendations for fuda

- **Web server mode:** poll once per repo/branch on the server, not once per open browser. 20 viewers
  then cost the same as one: 12 checks/min. On GitHub those checks are free 304s, on Azure about
  0.8 TSTU per 5 min.
- **Desktop (Wails) mode:** each app polls with its own user's token, so each user's own budget
  (GitHub 5,000/h, Azure 200 TSTU per 5 min) absorbs only that user's 720 checks/h. Nearly all of them
  are free 304s on GitHub.
- Run a full sync only when the head SHA changes. That turns the expensive tree and blob reads into a
  per-push cost, not a per-poll cost.
- Always send `Authorization` on GitHub polls. Without it, 304s are charged and the limit is 60/h.
- Do the polling in Go, not in the webview. GitHub sends `Cache-Control: private, max-age=60`, so a
  browser HTTP cache could serve a stale head for up to 60 s. Go's `net/http` client does not cache.
- Honour `Retry-After`, plus GitHub's `x-ratelimit-*` and Azure's `X-RateLimit-Remaining`, and back
  off when they appear.

## Details

### GitHub

**What the docs say.** The best-practices page says a conditional request "does not count against
your primary rate limit if a `304` response is returned and the request was made while correctly
authorized with an `Authorization` header" (via Context7 `/websites/github_en_rest`, and the page
itself). It also says "You should subscribe to webhook events instead of polling" and, if polling,
to honour `x-poll-interval`.

**Primary limits** (rate-limits page): unauthenticated 60/h per IP; personal access token 5,000/h;
GitHub App installation 5,000/h, scaling to 12,500 (15,000 on Enterprise Cloud); `GITHUB_TOKEN`
1,000/h per repo.

**Secondary limits** (same page): 100 concurrent requests; 900 points/min for REST, where most GETs
cost 1 point; 90 s of CPU time per 60 s of real time. The page does not say whether a 304 costs a
secondary point. Even if it does, a server polling once per branch makes 12 points/min. The
worst case, 20 apps sharing one token and polling separately, makes 240/min. Both are under 900.

**Live probe** (2026-10-08, public repo `golang/go`, branch `master`):

| Request | Auth | First response | Conditional responses | Primary limit used by 304s |
| --- | --- | --- | --- | --- |
| `GET /repos/golang/go/git/ref/heads/master` | none | 200, weak ETag | 304 | **1 each** (remaining went 56 to 55) |
| `GET /repos/golang/go/commits/master` | none | 200, weak ETag | 304 | **1 each** (54 to 53) |
| `GET /repos/golang/go/git/ref/heads/master` | token | 200, `X-Poll-Interval: 300` | 304 ×3 | **0** (used stayed 308) |
| `GET /repos/golang/go/commits/master`, `Accept: application/vnd.github.sha` | token | 200, body = 40-byte SHA, `Etag: "<sha>"` (strong) | 304 ×3 | **0** (used stayed 309) |
| `GET /repos/golang/go/branches/master` | token | 200, large JSON | 304 ×3 | **0** (used stayed 310) |

How the endpoints compare:

- `commits/{ref}` with the `vnd.github.sha` media type is the best choice. It has the smallest body
  on a change, its ETag equals the SHA so no extra state is needed, and it sends no poll-interval
  hint. The commits reference documents this use: the media type returns the SHA-1, which can be
  compared "by providing the local SHA-1 reference as the ETag". `ref` accepts a branch name or
  `heads/BRANCH`.
- `git/ref/heads/{branch}` works the same way, but its response sends `X-Poll-Interval: 300`. The
  best-practices page says to wait at least that long before polling the same endpoint again, so
  polling it every 5 s ignores GitHub's stated guidance.
- `branches/{branch}` returns a large body on change. It is no better.
- GraphQL has no conditional requests and is charged in points, so it is not cheaper for this.

### Azure DevOps

**Branch-head calls compared** (live, anonymous, public repo `dnceng-public/public/dotnet-public-wiki`,
api-version 7.1):

| Call | Body | `X-RateLimit-Cost` (TSTU) | ETag |
| --- | --- | --- | --- |
| `refs?filter=heads/main` | small JSON, `objectId` = head SHA | no header sent | none |
| `stats/branches?name=main` | ~640 B | 0.01176 | none |
| `commits?searchCriteria.itemVersion.version=main&$top=1` | ~650 B | 0.01425 | none |
| `items?path=/&versionDescriptor.version=main` | ~950 B, includes `commitId` | 0.00581 | none |

Every response had `cache-control: no-cache, no-store, must-revalidate` and no `ETag` or
`Last-Modified`. I found no Azure DevOps documentation of conditional GETs on Git REST resources.
`refs?filter=` is a prefix match ("starts with"), so `heads/main` also matches `heads/main-old`. Pick
the entry whose `name` equals `refs/heads/{branch}` exactly. The scope needed is `vso.code`.

The cheapest change detector is `refs?filter=heads/{branch}`, compared with the last synced SHA. The
fuda Azure source already reads items by branch. Only the head check is new.

**How TSTU limits work** (rate-limits page):

- One TSTU is the average load of a typical user over five minutes. Normal activity spikes to 10 or
  fewer, and the limit is 200 TSTUs in any sliding 5-minute window, per user identity (and per
  pipeline).
- When an identity goes over, requests are first **delayed**, from milliseconds up to 30 s, and can
  later be **blocked** with HTTP 429 `TF400733`. Delays stop within five minutes of consumption
  falling.
- Headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After`, and
  `X-RateLimit-Delay`. All except the delay header are sent **before** delays start, so a client can
  slow down ahead of time. `X-RateLimit-Cost`, when present, gives the TSTU cost of that request.
- TSTUs cannot be calculated with a formula. The usage-monitoring page shows real costs.

**At this load:** 60 calls per 5 min at ≤0.015 TSTU each is ≤0.9 TSTU per identity, against a limit of
200. A full sync on each push costs more (tree plus N blob reads) but happens only per push. A
web-server deployment on a single PAT or service identity has all of its usage counted against that
one identity. Polling once per branch, not once per viewer, keeps that cost flat.

### Push signals without a public server

- **GitHub webhooks** are HTTP POSTs to a URL GitHub can reach, so a desktop app cannot receive them
  without a public endpoint or relay. `gh webhook forward` "is only designed for use during testing
  and development". It is not supported in production, and only one person per repo or org can use it
  at a time.
- **GitHub Events API** (`/repos/{o}/{r}/events`) is still polling. Docs: "not built to serve
  real-time use cases", latency "anywhere from 30s to 6h", and it has its own `X-Poll-Interval`
  (typically 60). It is worse than polling the branch head.
- **Azure DevOps service hooks** support "Code pushed" events. Consumers include Webhooks (needs a
  public URL), **Azure Storage** ("Insert a message in a Storage queue"), and **Azure Service Bus**
  (queue or topic). A desktop app could read a Service Bus subscription or poll a Storage queue with
  no inbound port. Costs: an Azure subscription and resource, a shared credential distributed to
  every desktop app, and *Edit subscriptions* permission (project admins by default) to set it up.
  It is not free, and not something fuda can switch on by itself.
- Third-party relays such as smee.io are not first-party and are not meant for production. They
  are not evaluated here.

## Unverified

- Whether GitHub 304s cost a **secondary**-limit point. The docs are silent. The load above stays
  under 900 points/min either way.
- Whether GitHub's API ever serves a cached head for up to 60 s behind `Cache-Control: max-age=60`.
  The `Vary: Authorization` header and authenticated `private` caching suggest it does not, but I did
  not test how quickly a push shows up.
- Azure behaviour was probed **anonymously** on a public project. Authenticated responses (PAT or
  Entra) may have different headers, including ETag, and different costs. Repeat the probe with a PAT
  against the target org before relying on the numbers.
- The `refs` call sent no `X-RateLimit-Cost`. I assume it is at most as expensive as the other calls,
  but that is not measured.
- GitHub `X-Poll-Interval` values can change under load. The value of 300 was seen once on
  `git/ref`.

## Sources

- GitHub, best practices (conditional requests, polling, x-poll-interval):
  https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api
- GitHub, rate limits (primary and secondary):
  https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
- GitHub, get a commit (`application/vnd.github.sha`, SHA as ETag):
  https://docs.github.com/en/rest/commits/commits#get-a-commit
- GitHub, git references: https://docs.github.com/en/rest/git/refs#get-a-reference
- GitHub, events API latency: https://docs.github.com/en/rest/activity/events
- GitHub, webhook forwarding for testing:
  https://docs.github.com/en/webhooks/testing-and-troubleshooting-webhooks/using-the-github-cli-to-forward-webhooks-for-testing
- Context7 `/websites/github_en_rest`, the same best-practices text, retrieved 2026-10-08.
- Azure DevOps, rate and usage limits (TSTU, headers):
  https://learn.microsoft.com/en-us/azure/devops/integrate/concepts/rate-limits?view=azure-devops
- Azure DevOps, integration best practices:
  https://learn.microsoft.com/en-us/azure/devops/integrate/concepts/integration-bestpractices?view=azure-devops
- Azure DevOps, Refs - List 7.1:
  https://learn.microsoft.com/en-us/rest/api/azure/devops/git/refs/list?view=azure-devops-rest-7.1
- Azure DevOps, Git limits: https://learn.microsoft.com/en-us/azure/devops/repos/git/limits?view=azure-devops
- Azure DevOps, service hooks (consumers, permissions):
  https://learn.microsoft.com/en-us/azure/devops/service-hooks/overview?view=azure-devops
- Live probes on 2026-10-08 with `curl` and `gh api -i` against `api.github.com/repos/golang/go/...`
  and `dev.azure.com/dnceng-public/public/_apis/git/repositories/2459d599-fdb2-4d28-9810-daeec061cf90/...`.
