# About page media

Portrait and lifestyle photos are on **Cloudflare R2** (same Personal CDN path as Surfing).
The browser loads absolute `https://pub-….r2.dev/about/…` URLs — not the cluster nginx
tree in the basement. Local WebPs under `public/media/hero/` stay in the image as a
fallback if R2 is unreachable.

## Live objects

Bucket: `dasmlab-home`  
Public base: `https://pub-29bde7a836c744729bebe74bfd4008a2.r2.dev/about`  
Keys: `about/<file>.webp` with `Cache-Control: public, max-age=31536000, immutable`

| Slot | Object | Source file |
|------|--------|-------------|
| Portrait | `portrait.webp` | `/home/dasm/me-suit.jpg` |
| Windsurf | `dasm_ride.webp` | `src/assets/dasm_ride.jpg` |
| Snowboard | `me_ride_2.webp` | `src/assets/me_ride_2.jpg` |
| Country | `me_home.webp` | `src/assets/me_home.jpg` |
| Garden | `me_plant_1.webp` | `src/assets/me_plant_1.jpg` |
| Leah | `me_leah_1.webp` | `src/assets/me_leah_1.jpg` |
| Sail | `me_sail_1.webp` | `src/assets/me_sail_1.jpg` |
| Baking | `me_baking_1.webp` | `src/assets/me_baking_1.jpg` |
| Pies | `me_baking_2.webp` | `src/assets/me_baking_2.jpg` |

Re-publish after converting new files:

```bash
python3 scripts/publish-about-media.py
```

`dasm_ride.jpg` is a 6000×4000 / ~11 MB original; the WebP is 1600px on the long edge.

## CORS

R2 bucket Settings → CORS (same origins as Surfing):

- `https://dasmlab.org`
- `https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org`
- `https://dev-dasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org`

Methods: `GET`, `PUT`, `HEAD`. Headers: `*`.

## Override / custom domain

`*.r2.dev` is rate-limited (Cloudflare public development URL). For production SLA,
attach a custom domain on the bucket (e.g. `media.dasmlab.org`) once `dasmlab.org`
is on Cloudflare, then set build env:

`VITE_ABOUT_MEDIA_BASE=https://media.dasmlab.org/about`

Until then, `src/data/hubs.js` `ABOUT_MEDIA_BASE` defaults to the pub URL above.

The homepage `design-carousel-service` feed is a separate spinner (infra / MCP diagrams).
These personal photos belong on About, not that service.
