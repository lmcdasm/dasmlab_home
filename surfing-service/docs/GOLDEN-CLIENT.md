# Surfing — golden Personal CDN client (fast path)

Media bytes are on **Cloudflare R2 + CDN**. The browser must load/download from
absolute `https://pub-….r2.dev/...` URLs (or a future custom domain) — **not**
through `/api/surfing/serve` or the basement PVC.

## Current proof (2026-08-01)

| Check | Result |
|-------|--------|
| Manifest CDN URLs | 178/178 `https://` |
| Sample ~5.8 MB via CDN | ~0.31 s total |
| Same via `/serve` (PVC path) | ~0.46 s (slower; avoid) |
| FE `mediaUrl()` | passes through absolute URLs |
| FE download | uses CDN URL directly (no `?download=1` on R2) |

## Why it feels “fast fast”

1. Publish rewrote album URLs to R2 public base  
2. SPA links `<img>` / `<video>` / download straight at CF edge  
3. New publishes set `Cache-Control: public, max-age=31536000, immutable`

## Limits

- `*.r2.dev` is **rate-limited** (CF public development URL). For production
  speed/SLA, attach a custom domain on Cloudflare (`media.surfing.dasmlab.org`)
  once `dasmlab.org` is on CF.
- PVC remains staging for drafts; published gallery should never stream via OCP.

## Soft-hide (not hard drop)

`DELETE /days/:id/media/:mediaId` sets `hidden: true` on the manifest row and
**does not** delete PVC or R2 bytes. `GET /days` filters hidden items.
Keeping the row also stops `/data/preload` from re-importing the same file on
pod restart. Hard drop lands later in `dasmlab-cdn-mgr`.

`ListDays` reloads the manifest from disk so multi-replica pods see hides
written by siblings.

## Notes + kinds

| Field | Role |
|-------|------|
| `caption` | Short title on the card |
| `notes` | Longer session story |
| `kind` | `photo` \| `video` \| `other` (gallery section) |
| `external_url` | Open-out link (videos, Garmin, iPhone shares) |

- `PATCH /days/:id/media/:mediaId` — update notes/caption/kind/link  
- `POST /days/:id/media/link` — add link-only **More** items (no bytes)
- `POST /days/:id/theme/generate` — sample photos → AI banner/background on R2
  (`SURFING_AI_*`; OpenAI now, cheapcloud farm later)

## Platform successor

`dasmlab-cdn-mgr` + mock-me **Content Management → Self-Serve Cloud Personal CDN**.

Value prop (share-link vs Meta, keys, premium watermark / NFT registry, AI theme):
`docs/VALUE-PROP-PERSONAL-CDN.md`.
