# Value prop — Self-Serve Cloud Personal CDN

**Status:** pseudo-design / product narrative (2026-08)  
**Golden UX today:** Surfing (Windsurfing Trips) on dasmlab.org  
**Platform path:** `dasmlab-cdn-mgr` · cheapcloud · mock-me Content Management

---

## Why people do this (instead of Facebook)

| Pain with social upload | What we sell |
|-------------------------|--------------|
| Meta/Google keep a copy of your source | You hold keys; bytes live on **your** bucket / pass-through cloud bill |
| Algorithm owns distribution | You share a **curated box link** — friends open your gallery, not a feed |
| Hard to move later | Index + migrate across backends (corp → private → other regions) |
| “Delete” is theater | Soft-hide now; hard drop later under your policy |

**Core pitch:** put the **link** on Facebook (or nowhere at all). The media stays in a disk/object store you control, encrypted with your key at the bucket and (where applicable) at the infra layer, with geo/region choices at scale.

Two audiences, one stack:

1. **Don’t “do social”** — private/curated invite links to friends/family.  
2. **Do social carefully** — outbound share button; socials fetch *your* CDN, not a copy you uploaded forever.

---

## Commercial sketch (unchanged)

- Cloud storage / egress = **pass-through** on the user’s account.  
- ~**$1** figurative DASMLAB fee = index, orderly catalog, site-CDN, moves.  
- User holds encryption / access keys (OAuth2/SSO later).

---

## Premium foreshadow

| Capability | Notes |
|------------|--------|
| Unique watermark on outbound pulls | Tied to account / realm identity |
| Content provenance timestamp | “Published at …” for later AI-era claim of origin |
| `dasmlab-block-nft-registry` | Optional on-chain / NFT-style content registry for premium assurance |
| Unlock codes on shares | Private gallery invite; email/SMS delivery later |
| AI audio clean / overlay | Level loud/unmastered clips; starts muted in gallery today |

Registry is a **longer-term scaffold**: associate content to cloud provider pass/id/usage as people port from corp → private → other ops styles. Not required for golden Surfing path.

Stub name for future service/repo: **`dasmlab-block-nft-registry`**.

---

## Share tollgate (MVP)

Album/media share buttons mint a **DASMLAB-fronted** link (`GET /s/:token`): allow → meter → redirect to CDN.  
Social posts carry **metadata text + tollgate URL** (not a source upload).  
FinOps story: hits feed cheapcloud so storage/CDN use stays positive/neutral.

Detail: `docs/SHARE-TOLLGATE.md`.

---

## Videos visualization (interview-me pattern)

Album **Videos** use a Discovery-style satellite map (hub = Videos) with a **tight side panel** (not a popup):

- Plays / tags / CDN status  
- Publish · Share · Play muted  
- **Name tags** — plain text only, owner-approved; no links, no friend graph  
- Double-click / Cards → card browser grid  

Same control story: you own visibility; tags are just names on content.

---

## Surfing as golden client

- Theme: **Windsurfing Trips** (sports motif template; other sports later).  
- Gallery sections: **Videos** (open-out), **Photos**, **More** (Garmin / iPhone / activity shares).  
- Per-item **notes** + caption + optional outbound link (`PATCH …/media/:id`).  
- Share-link CTA in-app; publish → R2/CDN.  
- Hide = soft tag (`hidden`) in manifest — **no hard drop yet**.

See also: `surfing-service/docs/GOLDEN-CLIENT.md`.

---

## Later: AI sport-theme generator (running-translate pattern)

**MVP live path:** `POST /days/:id/theme/generate`

1. Sample a few album photos (CDN absolute URLs).  
2. Chat model returns palette + style brief (vision when samples exist).  
3. Image model returns **banner** + soft **background** art.  
4. Bytes land on R2 under `surfing/albums/{dayId}/theme/…`; day.theme stores URLs + palette.  
5. Surfing FE paints hero banner + workspace wash from that theme.

**Today (dev):** `SURFING_AI_PROVIDER=openai` → `api.openai.com`.  
**Next:** same OpenAI-compatible client pointed at **cheapcloud** managed AI farm (vLLM on spot capacity) — flip `SURFING_AI_BASE_URL` / provider; no FE change.

User story unchanged: a few pics + “my Bonaire trip / camping / first day of school” → page skin tailored to *their* CDN content.
