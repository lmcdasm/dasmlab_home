# About page media

Portrait and lifestyle photos ship as local WebP under `public/media/hero/`. The About page
frames and captions stay the same; only the bytes inside the frames changed.

Sources:

| Slot | WebP (served) | Source file |
|------|----------------|-------------|
| Portrait | `portrait.webp` | `/home/dasm/me-suit.jpg` |
| Windsurf | `dasm_ride.webp` | `src/assets/dasm_ride.jpg` |
| Snowboard | `me_ride_2.webp` | `src/assets/me_ride_2.jpg` |
| Country | `me_home.webp` | `src/assets/me_home.jpg` |
| Garden | `me_plant_1.webp` | `src/assets/me_plant_1.jpg` |
| Leah | `me_leah_1.webp` | `src/assets/me_leah_1.jpg` |
| Sail | `me_sail_1.webp` | `src/assets/me_sail_1.jpg` |
| Baking | `me_baking_1.webp` | `src/assets/me_baking_1.jpg` |
| Pies | `me_baking_2.webp` | `src/assets/me_baking_2.jpg` |

Original JPEGs remain in `src/assets/` (already in git). The About page does **not** import
those binaries into the JS bundle — it loads the converted WebP from `/media/hero/`.

`dasm_ride.jpg` is a 6000×4000 / ~11 MB original; the WebP is resized to 1600px on the
long edge.

## Optional CDN override

Set build env `VITE_ABOUT_MEDIA_BASE=https://…cdn…/about` (no trailing slash) and About
will load `{base}/portrait.webp`, `{base}/dasm_ride.webp`, etc. Until that is set, local
WebP is used.

The homepage `design-carousel-service` feed is a separate spinner (infra / MCP diagrams).
These personal photos belong on About, not that service.
