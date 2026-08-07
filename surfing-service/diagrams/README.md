# Surfing diagrams

Sources are `*.d2` (same sparse style as mini-acm / monday). Render locally:

```bash
cd surfing-service/diagrams
d2 surfing-cdn-flow.d2 surfing-cdn-flow.svg
d2 surfing-cdn-migrate.d2 surfing-cdn-migrate.svg
d2 surfing-cdn-layout.d2 surfing-cdn-layout.svg
```

| File | Purpose |
|------|---------|
| `surfing-cdn-flow.d2` | You → UI → API → PVC stage → Publish → R2 → CDN → visitors |
| `surfing-cdn-migrate.d2` | PVC-era `/serve` → publish job → CDN URLs → optional PVC cleanup |
| `surfing-cdn-layout.d2` | Bucket folder layout + three-product loop |

## Lifecycle (locked)

1. **Draft** — upload lands on PVC (`/data/media`). Manifest `url=/serve?id=…`, `published=false`.
2. **Publish** — `POST /days/:id/publish` copies each unpublished item to  
   `surfing/albums/{dayId}/original/{mediaId}{ext}`, rewrites `url` to `R2_PUBLIC_BASE_URL/…`, sets `published=true`.
3. **Serve** — FE uses absolute CDN URLs; `/serve` 302s if URL is already https.
4. **Cleanup (optional)** — after verify, delete PVC bytes for published items (manifest stays on OCP as source of truth).

## R2 key layout

```
surfing/albums/{dayId}/original/{mediaId}{ext}
surfing/albums/{dayId}/day.json          # optional snapshot
surfing/legacy/media/{mediaId}{ext}      # old flat keys if any
```
