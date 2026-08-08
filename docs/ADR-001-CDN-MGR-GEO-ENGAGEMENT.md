# ADR 001 — dasmlab-cdn-mgr + GEO/SEO/engagement (2.x → 3.0)

**Status:** Accepted for radar / not blocking home 2.0 visual polish  
**Date:** 2026-08-07  
**Deciders:** Technologies DASMLAB Inc.

## Context

`dasmlab_home` is becoming an Engineering Knowledge Network with public Activity (Phase A), demo visitors, and topical hubs. Surfing already moves media toward object/CDN origins. A separate platform — **dasmlab-cdn-mgr** — is intended for realms, BYO backends, index, publish, and migrate (see TARGET-0.7.0 constellation notes).

GEO (generative engine optimization), SEO spine, and engagement measurement will eventually need a **shared media + edge layer**, not one-off nginx proxies per site.

## Decision

Keep **dasmlab-cdn-mgr** on the 2.x / 3.0 incorporation path:

| Horizon | Scope |
|---------|--------|
| **2.0.x** | Home content/UX + Activity Phase A; Surfing origin work; demo facades |
| **2.x** | Wire home hubs / Activity tracks to CDN-mgr concepts (index, publish hooks) where cheap |
| **3.0** | Full CDN-mgr product: multi-realm, GEO-aware edge, engagement sinks, commercial envelopes |

Do **not** block current showcase polish on CDN-mgr shipping.

## Consequences

- Document links from project hubs (`surfing`, `cheapcloud`, `mock-me`) back to CDN-mgr ADR.
- Engagement events may later dual-write to CDN-mgr analytics without changing the Activity contract.
- Open-core CE stays separate from production CDN keys / realm ops.

## Related

- [TARGET-2.0.0.md](./TARGET-2.0.0.md)
- [DEMO-VISITOR-CONTRACT.md](./DEMO-VISITOR-CONTRACT.md)
- Blueprint: `/home/dasm/dasmlab-cdn-mgr/docs/TARGET.md` (when present)
