# About page media

## What happened to the old photos?

The previous About carousel used missing / lost binaries (and later gray SVG placeholders). For 2.0 we shipped **illustrated SVG stand-ins** under `public/media/hero/` so the page is never empty. Captions (windsurf, snowboard, Leah, etc.) stayed.

## How to drop real photos

### Option A — local files (fast)

Place JPEGs/WebPs in `public/media/hero/` and update `AboutPage.vue` paths, or keep the SVG names and replace files:

| Slot | Suggested filename |
|------|-------------------|
| Portrait | `portrait.jpg` (or keep `portrait.svg`) |
| Lifestyle | `lifestyle-windsurf.jpg`, `lifestyle-snowboard.jpg`, … |

### Option B — Surfing / CDN (preferred)

1. Publish public media in Surfing.
2. Set build env `VITE_ABOUT_MEDIA_BASE=https://…cdn…/about` (no trailing slash).
3. About will load `{base}/portrait.jpg`, `{base}/windsurf.jpg`, etc.

Until `VITE_ABOUT_MEDIA_BASE` is set, the site uses the local SVG illustrations and shows an on-page note.
