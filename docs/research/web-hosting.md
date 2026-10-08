# Free hosts for the fuda web shell

Checked 2026-10-08, using only the providers' own pricing and docs pages. Free tiers change often, so recheck them before you commit to one.

What fuda needs: one Go binary in one container, HTTPS, a secret env var for cookie encryption, a small persistent
disk for `FUDA_CACHE_DIR`, and either no sleep or a fast wake.

## Comparison

| Host | Free compute | Persistent disk on free | Sleep / cold start | Custom domain + HTTPS | Card required |
|---|---|---|---|---|---|
| Fly.io | None. New orgs get no free tier, only a trial of 2 machine-hours or 7 days [1] | No. Volumes cost $0.15/GB/mo [1] | Paid machines run always-on. The cheapest is about $2.19/mo [1] | Yes (paid) | Yes, after the trial [1] |
| Render | 750 instance-hours/month per workspace. CPU/RAM not stated on the free page [2] | No. "Persistent disks" are paid only [2] | Spins down after 15 min idle, and waking takes "about one minute" [2] | Custom domains supported on free [2] | No for the free tier. Without a card, going over bandwidth suspends the service instead of billing [2] |
| Koyeb | 1 free instance: 0.1 vCPU, 512 MB RAM, 2 GB SSD, Frankfurt or Washington only [4] | No. Free instances "can't be used with Volumes" [4] | Scales to zero after 1 h without traffic [4] | 10 custom domains free [3] | Not stated [3] |
| Google Cloud Run | 180k vCPU-s, 360k GB-s and 2M requests per month, 1 GB egress from North America [5] | No local disk. You would need a bucket or other external storage, which this check did not cover | Scales to zero, and min-instances cost money (not confirmed on a primary page, see note) | Yes | Yes, a billing account is required [5] |
| Google Compute Engine e2-micro | 1 e2-micro in us-west1, us-central1 or us-east1 [5] | Yes, 30 GB standard PD [5] | No sleep (it is a VM) | Yes, but you run Caddy or similar yourself | Yes [5] |
| Oracle Cloud Always Free | Ampere A1 with 2 OCPU and 12 GB total, or 2x E2.1.Micro (1/8 OCPU, 1 GB) [6] | Yes, 200 GB block storage total [6] | No sleep, but Oracle may reclaim an instance that stays under 20% CPU, network and memory (A1) for 7 days [6]. A1 capacity is often "out of host capacity" [6] | Yes, DIY (Caddy etc.). 10 TB/month egress [6] | Uncertain. The docs page doesn't say. Signup is known to ask for a card, but this check did not confirm it on a primary page |
| Railway | Free plan: $1 credit/month, 1 vCPU and 0.5 GB per service. Trial: $5 for 30 days [7] | Yes, volumes up to 500 MB [7] | Not stated for services [7]. $1/month likely runs out before the month ends for an always-on service (my estimate) | Not checked | Trial needs no card [7]. Not stated for the Free plan |
| Northflank | Sandbox: 2 free services, "Always-on-compute, no sleeping" [8]. CPU/RAM figures not given | Not stated for free. Disk is $0.15/GB/mo [8] | No sleep [8] | Not checked | Yes, "all users must add a payment method" [9] |

Note: the Cloud Run pricing page [10] did not render for me, so the min-instance billing detail above is unverified.

## Recommendation

- **Real free, always-on, with disk:** an **Oracle Cloud Always Free** VM. Use an A1 shape if capacity allows, otherwise
  E2.1.Micro. Put Caddy in front for automatic HTTPS, run fuda under systemd or Docker, set the secret in an env file,
  and point `FUDA_CACHE_DIR` at the boot volume. The costs: you maintain the VM yourself, Oracle can reclaim it when idle,
  and a card is likely needed at signup. **Google e2-micro** works the same way, with less RAM and 1 GB egress.
- **Zero-ops and free, if you accept a cold start:** **Render free**. The cache does not survive a restart, so fuda
  re-syncs on every wake. That is acceptable because a failed or missing cache only means one fresh sync.
  Waking takes about a minute.
- **Cheapest managed option with disk and no sleep:** **Fly.io** at roughly $2-3/month (shared-cpu-1x plus a 1 GB
  volume at $0.15). It is not free, but it fits fuda best.
- Skip Koyeb free (no volume, 0.1 vCPU), Cloud Run (no disk) and Northflank (card required, disk not free) for this use.

## Sources

1. https://docs.fly.io/about/pricing
2. https://render.com/docs/free
3. https://www.koyeb.com/pricing
4. https://www.koyeb.com/docs/reference/instances
5. https://docs.cloud.google.com/free/docs/free-cloud-features
6. https://docs.oracle.com/en-us/iaas/Content/FreeTier/freetier_topic-Always_Free_Resources.htm
7. https://railway.com/pricing
8. https://northflank.com/pricing
9. https://northflank.com/docs/v1/application/billing/pricing-on-northflank
10. https://cloud.google.com/run/pricing (did not render)
