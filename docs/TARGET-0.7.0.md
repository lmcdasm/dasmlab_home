# Target map: 0.7.0 — Lab Constellation + Surfing CDN

**Branch (proposed):** `2026-lab-constellation`  
**Baseline:** 0.6.x spring FE update (hero/orbs, approach switcher, teal/sage tokens)  
**Target release:** **0.7.0**  
**Vision:** Kill the dashboard-card soup. Make dasmlab.org feel like the same family as **interview-me**, **mini-acm**, and **thao-rip** — one composition, interactive map, clear CTAs — while getting Surfing video bytes **off basement spinning disk**.

---

## Why now

| Pain | Fix |
|------|-----|
| Homepage still looks like a Quasar tutorial + commit carousel | Lab constellation + architecture stage (steal from interview-me / mini-acm) |
| WhatsNew is stale noise (2023–2024 commits) | Deprecate service end-to-end |
| Surfing streams from OCP PVC (`lvms-vg1` / basement) | Origin → cheap object store + edge CDN |
| cheapcloud is dry-run / Azure-shaped, underused | Evolve into the **cheap origin broker** for Surfing (and friends) |

---

## Workstreams

### A. Deprecate WhatsNew (do first — low risk)

**Frontend (`dasmlab_home`):**
- Remove `WhatsNew.vue`, IndexPage fetch to `whatsnew.svc.dasmlab.org`, “Live signal” block
- Drop or retire `whatsnew-service` card on Backend Projects
- Update README / voice docs

**Cluster / GitOps (`dasmlab-live-cicd`):**
- Disable Argo Applications for `whatsnew-service` on `2026-prod-1`, `2026-prod-2-1` (+ legacy `dasmlab-prod-1` if still live)
- Empty/remove `whatsnew-service/live/` so prune deletes Deployment/Service/Route
- Delete bootstrap RBAC + namespace secrets (`whatsnew-service-secrets`)
- Optional later: archive `dasmlab/whatsnew-service` repo, revoke GitHub App

**Keep:** DesignCarousel (independent), VisitCounter, approach switcher (simplify later).

---

### B. Homepage redesign — “Lab Constellation”

Stay on **Vue 3 + Quasar + Vite** (same stack as interview-me / mini-acm / thao-rip). Do **not** rewrite to Next/React.

#### Steal from our own boilerplate (priority order)

| # | Source | Pattern to port |
|---|--------|-----------------|
| 1 | interview-me `DiscoveryMindMap.vue` | Hub-and-spoke SVG: hub = **DASMLAB**, spokes = Frontend / Backend / AI·ML / Cloud / Infra / Security / Surfing |
| 2 | interview-me HomePage stage + detail aside | Click a lane → inspector panel (stats, links, produce/consume of knowledge) — **no overlapping map labels** |
| 3 | mini-acm `ArchitectureBoxView.vue` | Second section: “How the lab is wired” full-bleed systems diagram |
| 4 | thao-rip StudioPage motion + typography | Staggered `rise`, display/body font pairing — **recolor to `--dasm-*` teal/sage**, avoid cream+terracotta AI default |
| 5 | etcd/mini-acm action cards | Hover-lift only where interaction needs a container |

#### First viewport (hard rules)

One composition, brand-first:

1. **DASMLAB** as hero-level brand  
2. One headline + one short supporting line  
3. One CTA group (Surfing · Projects · Contact)  
4. Dominant visual = **Lab map** (edge-to-edge stage, not inset card soup)

Then **in the same release**: architecture wiring diagram + demoted carousel + lane chips (not 6 equal Quasar cards in the hero).

**Always improve:** thicker borders, darker ink, stronger spoke contrast, rise motion — pull forward from whatever interview-me / mini-acm / thao-rip shipped most recently. The map is a starting pattern, not a ceiling.

No WhatsNew. No card grid in the hero. Stats / project lists live in the side panel or below the fold.

#### External UI inspiration (curate, don’t clone)

Hot 2025–26 portfolio patterns that map cleanly to our stack:

| Idea | Takeaway for DASMLAB | Skip |
|------|----------------------|------|
| [Linear](https://linear.app) / [Vercel](https://vercel.com) product landings | Restraint, one visual idea, strong type | Dark-mode-by-default purple glow |
| [delowarhossain.dev](https://delowarhossain.dev) Creative-Folio | Typography-first, lab/experiments section | Heavy WebGL if it tanks Core Web Vitals |
| Recruiter Mode vs Cinematic Mode (Victor Alves-style) | Maps to our existing **approach switcher** (scan vs deep) | OS-simulator chrome |
| Liquid-glass / glassmorphism portfolios | Soft blur rails only (thao-rip already does this lightly) | Full cinematic video backgrounds |

**DASMLAB differentiator:** interactive **systems map** (interview-me/mini-acm DNA), not another glassmorphic resume.

#### Suggested page flow after 0.7.0

```
Hero (brand + one line + CTAs)
  └─ Lab map (hub-spoke + inspector)     ← parallel track A
How the lab is wired (architecture SVG) ← parallel track B
Featured experiments (DesignCarousel — demoted)
Project lanes (chip row, not card soup)
```

Keep iterating edges/shades/contrast against sibling repos; never freeze on the first map draft.

---

### C. Surfing CDN — options matrix

**Today:** Go API writes files to PVC `/data/media`, serves via `GET /serve` through OCP Route / nginx proxy. Origin = basement disk. Max upload 500MB, PVC 20Gi.

**Target split:**
- **API stays on OCP** — days, metadata, auth, upload orchestration  
- **Bytes leave the cluster** — object storage + CDN URL in manifest  

#### Option 1 — Cloudflare cache / R2 (recommended Phase 1)

| Piece | Choice |
|-------|--------|
| Storage | **Cloudflare R2** (S3 API, **zero egress**, 10GB free tier) |
| Edge | Cloudflare CDN in front of `media.surfing.dasmlab.org` (or public R2 custom domain) |
| DNS | Zone you already control; `cf_credss` exists locally |
| Prior art | `/home/dasm/sailgp_cloudflare_exp` (`r2_stub/`, Workers) |

**Pros:** Free egress, you already have CF muscle memory, stays under free tier for lab-scale video.  
**Cons:** CF free-tier fair-use; need Workers or signed URLs if uploads stay private.  
**Fit:** Best “get off basement disk this month” path.

#### Option 2 — Hyperscaler storage + CDN

| Cloud | Storage | CDN | Notes |
|-------|---------|-----|-------|
| Azure | Blob (Hot/Cool) | Azure CDN / Front Door | cheapcloud already mounts Azure secrets — natural if we wake cheapcloud |
| AWS | S3 | CloudFront | familiar; egress expensive at scale |
| GCP | GCS | Cloud CDN | similar egress tax |

**Cheapest disks among big-3:** usually **Azure Cool / Archive** or **GCS Nearline** for cold; for **hot video origin**, R2/Bunny/Hetzner beat all three on egress.

**Pros:** Enterprise features, regional choice.  
**Cons:** Egress (~$0.08/GB) murders hobby video. Only win if traffic stays tiny or you terminate behind Cloudflare anyway.

#### Option 3 — Bunny.net Storage + Stream *(other method #1)*

- Storage ~$0.005–0.01/GB, CDN ~$0.01/GB (Volume lower)  
- **Bunny Stream**: free transcoding/player, pay bandwidth — purpose-built for VOD  
- Region toggles = cheap+dirty (turn off expensive geos)

**Pros:** Best $/GB for actual video streaming; HLS out of the box.  
**Cons:** New vendor; less overlap with existing CF experiments.

#### Option 4 — Hetzner Object Storage + Cloudflare pull-through *(other method #2)*

- ~€6/mo includes **1 TB storage + 1 TB egress** (EU)  
- Put **Cloudflare** (or Bunny) as pull-through cache in front — Hetzner docs say “not a CDN, put one in front”  
- Classic cheapcloud energy: EU region, cancel when done, migrate with `rclone` / `s3cli`

**Pros:** Predictable EU pricing; great cold+warm origin.  
**Cons:** Origin alone is not global; need CDN layer; newer product under load caveats.

#### Phase recommendation

```
Phase 0  Deprecate WhatsNew; homepage constellation scaffold
Phase 1  Surfing → R2 origin + CF custom domain (reuse sailgp r2_stub patterns)
         Manifest URLs: https://media.surfing.dasmlab.org/<id>
         API: upload to R2 (presign), still CRUD on OCP
Phase 2  cheapcloud grows a "MediaBroker" mode:
         - providers: r2 | azure-blob | hetzner | bunny
         - dry-run → live; cheap regions; migrate/rclone jobs
Phase 3  If traffic grows: Bunny Stream for HLS; keep R2/Hetzner as archive
```

**Anti-pattern:** Azure/AWS/GCP CDN as primary video edge for this lab — pay the egress tax only if something else already lives there.

---

### D. cheapcloud updates (as we go)

**Current:** `ghcr.io/dasmlab/cheapcloud` on both 2026 clusters, `CHEAPCLOUD_DRY_RUN=1`, Azure secrets mount, 2Gi PVC, Route on apps domain. **No local source tree** under `/home/dasm` — only GitOps.

**Proposed evolution (Surfing-aligned):**

1. **Find/check out source** (`dasmlab/cheapcloud` or equivalent) into workspace  
2. Add provider interface: `Put` / `GetURL` / `Delete` / `Migrate`  
3. Config YAML: active provider + region + cost caps  
4. Surfing-service calls cheapcloud (or embeds the same Go package) instead of local `storage.go` filesystem  
5. Flip `CHEAPCLOUD_DRY_RUN=0` once R2 (or Hetzner) credentials are in `cheapcloud-secrets`  
6. Diagnostic sidecar (`cheapdiag`) reports latency / $ estimate per provider

Cheap+dirty ops we explicitly want:
- Spin provider for a demo week, migrate off, delete bucket  
- Prefer cheap regions (EU Hetzner, CF R2 auto)  
- Cap monthly spend in config; fail closed

---

## CI/CD (unchanged pattern)

Keep the interview-me / mini-acm style pipeline:

```
push → self-hosted runner → buildah/kaniko → ghcr.io/lmcdasm|dasmlab/...
     → sed __VERSION__ into k8s_envelope
     → commit dasmlab-live-cicd clusters/2026-prod-*/<app>/live/
     → Argo sync
```

New pieces:
- `surfing-service` gains R2/S3 env secrets (or talks to cheapcloud)
- Optional workflow for cheapcloud when source is wired
- Homepage image bump after constellation lands

---

## Acceptance criteria

- [ ] No homepage call to `whatsnew.svc.dasmlab.org`
- [ ] WhatsNew Argo apps gone / pruned on 2026 clusters
- [ ] Home first viewport: brand + constellation map + side panel (no WhatsNew, no 6-card hero)
- [ ] Surfing media URL points at CDN/object store, not PVC `/serve` for new uploads
- [ ] Documented provider matrix + chosen Phase 1 path (R2+CF default)
- [ ] cheapcloud has a written MediaBroker spike or ticket with secrets layout

---

## Decisions locked

1. **Phase 1 origin:** Cloudflare **R2** + CF CDN (confirmed 2026-08-01).  
2. **WhatsNew cluster teardown:** PR to `dasmlab-live-cicd` empties live + archives manifests so Argo prune runs on next sync.  
3. **Homepage visual scope:** Lab map **and** architecture wiring ship **in parallel** for 0.7.0.  
4. **“Constellation”** was only the first sketch name — keep raising the bar (edges, contrast, motion) against interview-me / mini-acm / thao-rip as those repos keep moving. Prefer snappy labels in UI: **Lab map**, **How the lab is wired**.

## Still open

1. **Public vs signed media URLs** for Surfing on R2?
