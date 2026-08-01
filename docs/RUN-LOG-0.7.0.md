# 0.7.0 run log — fixes / improvements as we execute PLAN-0.7.0-5x10

Rolling log of **≥10 improvements/fixes per task** while shipping the preview release.

**Start:** 2026-08-01 · **Preview live:** https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org/

---

## Task 1 — Preview scaffolding ✅

| # | Fix / improvement | Result |
|---|-------------------|--------|
| 1 | `dasmlab-home_preview-ocp.yaml` | Added |
| 2 | Host `dev-__OWNER__-dasmlab-home.apps…` | Confirmed |
| 3 | Port 8888 + nginx emptyDirs | Pod Ready |
| 4 | `deploy-preview.sh` → live-cicd | Pushed `lmcdasm.yaml` |
| 5 | `bootstrap-preview-ns.sh` GHCR pull | Secret created |
| 6 | `cleanup-preview.sh` | Added |
| 7 | Argo `dasmlab-home-previews` | Applied + Synced |
| 8 | `previews/README.md` | Added |
| 9 | Image `ghcr.io/lmcdasm/dasmlab-home` | Pushed tag |
| 10 | Docs TARGET + RUN-LOG | This file |

**Extra fixes while running:**
- PSA `privileged` on preview NS (nginx runs as root like prod)
- Hard Argo refresh when path first missing on `main`
- Re-bind anyuid SCC after SA exists

---

## Task 2 — CI ✅

| # | Fix / improvement | Result |
|---|-------------------|--------|
| 1 | Split prepare / build / preview / prod | `main.yml` rewritten |
| 2 | `is_prod` only on `main` | Done |
| 3 | `pull_request` + feature branches | Done |
| 4 | Concurrency `preview-{actor}` | Done |
| 5 | Wire `deploy-preview.sh` | Done |
| 6 | Keep Buildah/Kaniko + counters | Done |
| 7 | Prod dual-cluster for main only | Done |
| 8 | Remove FE-branch auto-merge | Done (PR-only) |
| 9 | Step summary in deploy script | Done |
| 10 | Branch `2026-lab-constellation` | Pushed |

**Extra:** `preview-cleanup.yml` on PR close; `.gitignore` for `in_progress:/`

---

## Task 3 — Ship + PR ✅

| # | Fix / improvement | Result |
|---|-------------------|--------|
| 1 | Branch `2026-lab-constellation` | ✅ |
| 2 | Commit UI (map/arch/WhatsNew out) | `e4d61f6` |
| 3 | Commit preview scaffolding | Same commit |
| 4 | Exclude `in_progress:/` | gitignore |
| 5 | `gh pr create` | https://github.com/lmcdasm/dasmlab_home/pull/10 |
| 6 | PR body + preview comment | Posted |
| 7 | live-cicd Argo PR | https://github.com/lmcdasm/dasmlab-live-cicd/pull/2 |
| 8 | Manual build while CI queued | Buildah OK |
| 9 | Image tag `v2026.08.01-e4d61f6` | GHCR |
| 10 | WhatsNew prune PR still open | https://github.com/lmcdasm/dasmlab-live-cicd/pull/1 |

---

## Task 4 — Cluster heal ✅

| # | Fix / improvement | Result |
|---|-------------------|--------|
| 1 | Apply Argo previews Application | Created |
| 2 | Bootstrap pull secret | OK |
| 3 | Wait pods Ready | 1/1 Running |
| 4 | Argo path-missing → hard refresh | Fixed |
| 5 | Route present | OK |
| 6 | HTTPS 200 | Verified (`curl -vk`) |
| 7 | IndexPage bundle has Lab map | Grep OK |
| 8 | Architecture strings present | In bundle |
| 9 | Surfing `/api/surfing` still proxies via nginx.conf | Same as prod image |
| 10 | Publish URL on PR | Commented |

---

## Task 5 — Harden (in progress / queued)

| # | Fix / improvement | Result |
|---|-------------------|--------|
| 1 | Stronger LabMap spoke contrast / stage border | Shipped in PR |
| 2 | Arch diagram shows R2 + cheapcloud path | Shipped |
| 3 | Mobile map stack CSS (900px) | Already in LabMap |
| 4 | Merge WhatsNew prune PR | Attempted / pending |
| 5 | preview-cleanup workflow | Added |
| 6 | R2 Phase 1 next (TARGET locked) | Surfing R2 live + album publish |
| 7 | cheapcloud MediaBroker next | Interface + `$20` media.example.yaml |
| 8 | Promote checklist on PR | In PR body |
| 9 | TARGET checkboxes | Update after merge |
| 10 | Steal next: mini-mock + running-translate preview patterns | Noted (appeared on live-cicd main) |

---

## CDN loop kickoff (2026-08-01)

| Item | Status |
|------|--------|
| Three-product loop locked in TARGET | mini-mock mock → cheapcloud profile → Surfing prod |
| Storage+CDN budget | **`$20/mo`** fail-closed (`surfing-cdn-storage`) |
| Surfing R2 backend | Live `v2026.08.01-48d8ec7-pub` |
| Lifecycle | Draft=PVC → `POST /days/:id/publish` → R2 albums |
| Bonaire migrate | **178/178 published** to `surfing/albums/{dayId}/original/…` |
| D2 diagrams | `surfing-service/diagrams/*.d2` + `.svg` |
| CI runner | `dev-2022-dasmlab-home` **online** |

## Working release coordinates

| Item | Value |
|------|-------|
| Preview URL | https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org/ |
| Home PR | https://github.com/lmcdasm/dasmlab_home/pull/10 |
| Home Image | `ghcr.io/lmcdasm/dasmlab-home:v2026.08.01-48d8ec7` |
| Surfing Image | `ghcr.io/lmcdasm/surfing-service:v2026.08.01-48d8ec7-pub` |
| Surfing CDN sample | `https://pub-973dadf41dd44837be5bcdd8913067a7.r2.dev/surfing/albums/…` |
| NS | `dasmlab-home-dev-lmcdasm` / `surfing-service-system` |
| Argo | `dasmlab-home-previews` + surfing-service live |

## Hotfix — HAProxy CERT (2026-08-01)

- Missed mini-mock `ensure-prod-cert.sh` pattern on first preview ship → `ERR_CERT_COMMON_NAME_INVALID` / HSTS
- Added `scripts/ci/ensure-preview-cert.sh` (ssh → 10.20.1.10, append CERTn, `./runme.sh`)
- Hooked from `deploy-preview.sh` (hands-free on every preview)
- Applied live: **CERT56**=`dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org`
- Verified: TLS CN matches FQDN, HTTPS 200 without `-k`

## Soft-hide + Windsurfing theme (2026-08-01)

| # | Change | Why |
|---|--------|-----|
| 1 | `MediaItem.Hidden` + DeleteMedia soft-hide | User asked: never hard-drop yet; UI remove must survive refresh |
| 2 | ListDays filters hidden + reloads manifest from disk | Multi-replica + preload reseed were likely refresh ghosts |
| 3 | Publish skips hidden | Don’t push hidden rows to CDN again |
| 4 | Surfing FE: Photos/Clips sections, share link, ocean theme | Sports motif + value-prop CTA |
| 5 | `docs/VALUE-PROP-PERSONAL-CDN.md` | Share-link vs Meta; premium NFT registry foreshadow |
| 6 | GOLDEN-CLIENT + TARGET decisions 8–9 | Lock the narrative |

## Theme MVP + ship (2026-08-01)

| # | Change | Why |
|---|--------|-----|
| 1 | `POST /days/:id/theme/generate` | Sample CDN photos → banner/bg + palette |
| 2 | OpenAI-compatible AI client | Dev=OpenAI; later=cheapcloud vLLM farm via BASE_URL |
| 3 | Theme art on R2 `…/theme/` | Same CDN as album media |
| 4 | FE applies banner/wash | See the imagined sport theme live |
| 5 | Secret `surfing-ai` + deploy env | Wire key without baking into image |
| 6 | Notes/kinds/Videos-Photos-More UI | Hot gallery MVP |
