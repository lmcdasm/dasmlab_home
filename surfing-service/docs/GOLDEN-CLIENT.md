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

## Platform successor

`dasmlab-cdn-mgr` + mock-me **Content Management → Self-Serve Cloud Personal CDN**.
