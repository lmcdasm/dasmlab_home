# Project & topic hub inventory (2.0)

Living checklist for the Engineering Knowledge Network on dasmlab.org.

## Public vanity hosts (CTAs)

Never list `dev-*` or `*.apps.2026-prod-1…` on visitor CTAs.

| Product | Vanity |
|---------|--------|
| Mock-Me | https://mock-me.dasmlab.org/demo |
| Interview-Me | https://interview-me.dasmlab.org/demo |
| Camera Scrape | https://camera-scrape.dasmlab.org/ |
| Surfing | https://dasmlab.org/surfing |
| CheapCloud | hub only (edge auth — no public Live yet) |
| mcp-explorer | hub/repo only (no live deploy) |

## Project hubs (first wave)

| Slug | Title | Lane | Live / demo | CE / source | Hub status |
|------|-------|------|-------------|-------------|------------|
| `dasmlab-home` | dasmlab-home | Frontend | https://dasmlab.org | github.com/lmcdasm/dasmlab_home | Live |
| `surfing` | Surfing CDN | Frontend / Backend | /surfing | private + future CE notes | Live |
| `mock-me` | Mock-Me | Cloud / Frontend | mock-me.dasmlab.org/demo | open-core/mock-me-ce (planned GH) | Live |
| `interview-me` | Interview-Me | AI / Frontend | interview-me.dasmlab.org/demo | open-core/interview-me-ce (planned GH) | Live |
| `cheapcloud` | CheapCloud | Cloud | hub only (operator auth) | open-core/cheapcloud-ce (planned GH) | Live |
| `camera-scrape` | Live Cams | Infra / Frontend / Backend | camera-scrape.dasmlab.org | open-core/camera-scrape-ce (planned GH) | Live |

## Project hubs (later)

| Slug | Title | Notes |
|------|-------|-------|
| `design-carousel-service` | Design Carousel | Backend lane — design-carousel.svc.dasmlab.org |
| `activity-cdp` | Public Activity spine | Surfing engagement |
| `etcd-synth` | etcd synthetic load | Demo dry-run only |
| `tutorial-generator` | Tutorial generator | Public repo exists |

## CE GitHub repos to create (when ready)

Push each `open-core/<name>/` folder — do not link until published:

- `dasmlab/mock-me-ce`, `dasmlab/interview-me-ce`, `dasmlab/cheapcloud-ce`
- `dasmlab/camera-scrape-ce`, `dasmlab/etcd-synth-ce`
- `dasmlab/dasmlab-go` (Apache 2.0)

## Topic hubs (first wave)

| Slug | Tech | Linked projects |
|------|------|-----------------|
| `vue` | Vue 3 | dasmlab-home, interview-me, mock-me |
| `gin` | Gin (Go) | surfing, mock-me, interview-me, camera-scrape |
| `oidc` | OIDC / Keycloak | home, mock-me, interview-me, camera-scrape |
| `metallb` | MetalLB | infra + camera-scrape |
| `quasar` | Quasar | dasmlab-home, siblings |
| `openshift` | OpenShift / GitOps | cheapcloud, camera-scrape, live-cicd |

## Labs

| Slug | Focus |
|------|-------|
| `activity-anon-cdp` | Anonymous visitor CDP on public sites |
| `surfing-r2-origin` | Surfing bytes off PVC → object/CDN |
| `demo-visitor-facade` | Cross-product fake mode pattern |

## SEO surface checklist

- [x] robots.txt
- [x] sitemap.xml
- [x] Organization + Person JSON-LD
- [x] Per-route title / description / OG
- [x] FAQPage on hubs with FAQ
- [x] Answer-first article template
- [x] History mode routing (nginx try_files)
- [x] Vanity Live/Demo CTAs (no apps.* / dev-*)
