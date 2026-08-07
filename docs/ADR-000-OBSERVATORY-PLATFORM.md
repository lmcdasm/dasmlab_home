# ADR-000 — Link to DASMLAB Observatory Platform (DOP)

**Status:** Accepted  
**Date:** 2026-08-07  
**Deciders:** Technologies DASMLAB Inc.

## Context

`dasmlab_home` is the public Engineering Knowledge Network. Company-wide Observatory architecture now lives in a dedicated platform repo so product ADRs (home CDN-mgr, etc.) do not collide with platform numbering.

## Decision

- **Platform ADRs** live at `/home/dasm/dasmlab-observatory-platform` (branch `2026-dop-v0`), starting with ADR-0001 Vision and ADR-9999 Innovation Principles.
- **First product specialization:** Digital Presence Observatory — platform `docs/adr/ADR-0400-digital-presence-observatory.md` (`/home/dasm/dasmlab-observatory-platform`).
- Home remains Activity producer and topical authority; DPO consumes Activity + GSC + edge + GitHub.
- Home [ADR-001 CDN-mgr GEO](./ADR-001-CDN-MGR-GEO-ENGAGEMENT.md) stays valid and local; CDN-mgr is not DOP.
- **Live repo:** https://github.com/lmcdasm/dasmlab-observatory-platform  
- **Live DPO (prod-1):** https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org  
- GitOps: Argo Application `dpo` → `dasmlab-live-cicd` `clusters/2026-prod-1/dpo/live`

## Consequences

- New Observatory features are gated by platform ADR-9999.
- When GitHub org `dasmlab` is ready, mirror this link to `github.com/dasmlab/observatory-platform`.

## Related

- Platform index: `/home/dasm/dasmlab-observatory-platform/docs/adr/README.md`
- [TARGET-2.0.0.md](./TARGET-2.0.0.md)
