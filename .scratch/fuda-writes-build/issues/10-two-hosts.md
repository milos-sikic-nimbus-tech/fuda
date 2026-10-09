# 10: Both hosts at once

**What to build:** A person can be logged in to GitHub and Azure DevOps at once. The Board picker lists Boards from both. Login and logout are per host, with one cookie per host. A "Log out of all" button clears both. An expired token sends only that host back to login.

**Blocked by:** 09: Azure DevOps login, Boards and writes

**Status:** ready-for-agent

- [ ] Picker shows Boards from both hosts
- [ ] Logging out of one host keeps the other working
- [ ] "Log out of all" clears both
- [ ] The desktop app has one keychain entry per host and the same per-host logout

Spec: `.scratch/fuda-writes/spec.md`.
