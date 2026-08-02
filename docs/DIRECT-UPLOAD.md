# Direct-to-R2 draft uploads

Batch uploads no longer need to double-hop through the cluster PVC.

## Flow

1. `POST /days/:id/media/presign` — reserve media id in album **index** (`origin=r2-draft`, key under `surfing/albums/{day}/draft/…`)
2. Browser **PUT** bytes to the returned R2 URL (DASMLAB only holds the index)
3. `POST /days/:id/media/:mediaId/complete` — HEAD verify on R2
4. `POST /days/:id/publish` — **promote** `draft/` → `original/` (CopyObject), flip `published` — no re-upload

Legacy multipart `POST /days/:id/media` remains as fallback (PVC draft).

FE queue: concurrency **3**, up to **3** retries, debounced gallery refresh.

## R2 CORS (required for browser PUT)

In Cloudflare R2 → bucket → Settings → CORS, allow PUT from prod + preview:

```json
[
  {
    "AllowedOrigins": [
      "https://dasmlab.org",
      "https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org"
    ],
    "AllowedMethods": ["GET", "PUT", "HEAD"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3600
  }
]
```

Without CORS, the FE falls back to cluster multipart (slower / more failure-prone for big batches).
