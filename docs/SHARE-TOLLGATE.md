# Share tollgate — social second, private first

**Status:** MVP in Surfing (2026-08)  
**Why:** We host the **index / pointers + metadata**, not a Meta-style source dump. Shares must go through **us** so we can allow/deny, meter pulls, and keep the cloud storage envelope **positive or neutral** via cheapcloud.

## Flow

```text
Owner hits Share → POST /shares (mint token + copy from metadata)
Friend opens tollgate URL → GET /s/:token
  ├─ social crawler → OG HTML (title/text/image from metadata)
  └─ human → allow-check → hit++ → 302 to CDN / album page
Hit event → log + optional POST SURFING_CHEAPCLOUD_HIT_URL
cheapcloud → watches daily burn, later forecasts / “if you did these” proposals
```

## API

| Method | Path | Role |
|--------|------|------|
| POST | `/shares` | Mint album or media share (`day_id`, optional `media_id`, optional `unlock_code`) |
| GET | `/s/:token` | Tollgate resolve |
| GET | `/shares/:token` | JSON meta / hit count |

## FE share sheet

- **Copy share link** — tollgate (metered)
- **Copy CDN link** — direct edge (power users)
- **Hottest:** WhatsApp, Facebook, Threads, X
- **Also:** LinkedIn, Instagram/TikTok (copy+paste — no web intent), device share sheet
- Share **text** is built from album/media caption + notes + location

## Premium foreshadow

- `unlock_code` on create → hash stored; resolve requires `?unlock=`
- Later: email / SMS delivery of codes to specific people
- Watermark + NFT registry remain separate premium tracks

## cheapcloud tie-in

Tollgate hits are the **consumption signal** for “how hard are friends pulling our hosted origin?”

Near-term:

1. Surfing logs `share_hit` + optional emit to `SURFING_CHEAPCLOUD_HIT_URL`
2. cheapcloud R2 / free-tier scrapers already show “where we are today”
3. Later: nightly RHOAI retrain on hit patterns + adapter price markets → simple price buckets (pop CDN, pop mgmt cluster, …) and end-user “if you did these” proposals

See also: `cheapcloud/docs/mediabroker.md`, `docs/VALUE-PROP-PERSONAL-CDN.md`.
