# Project & topic hub inventory (2.0)

Living checklist for the Engineering Knowledge Network on dasmlab.org.

## Project hubs (first wave)

| Slug | Title | Lane | Live / demo | CE / source | Hub status |
|------|-------|------|-------------|-------------|------------|
| `dasmlab-home` | dasmlab-home | Frontend | https://dasmlab.org | github.com/lmcdasm/dasmlab_home | Ship in 2.0 |
| `surfing` | Surfing CDN | Frontend / Backend | /surfing | private + future CE notes | Ship in 2.0 |
| `mock-me` | Mock-Me | Cloud / Infra | demo facade | mock-me-ce (planned) | Ship in 2.0 |
| `interview-me` | Interview-Me | AI / Frontend | demo facade | interview-me-ce (planned) | Ship Wave 1+ |
| `cheapcloud` | CheapCloud | Cloud | demo readonly | cheapcloud-ce (planned) | Ship Wave 2+ |

## Project hubs (later)

| Slug | Title | Notes |
|------|-------|-------|
| `design-carousel-service` | Design Carousel | Backend lane |
| `activity-cdp` | Public Activity spine | Surfing engagement |
| `etcd-synth` | etcd synthetic load | Demo dry-run only |
| `camera-scrape` | Camera scrape | After identity boundary |
| `tutorial-generator` | Tutorial generator | Public repo exists |

## Topic hubs (first wave)

| Slug | Tech | Linked projects |
|------|------|-----------------|
| `vue` | Vue 3 | dasmlab-home, interview-me, mock-me |
| `gin` | Gin (Go) | surfing, mock-me, interview-me |
| `oidc` | OIDC / Keycloak | home, mock-me, interview-me |
| `metallb` | MetalLB | infra lab narratives |
| `quasar` | Quasar | dasmlab-home, siblings |
| `openshift` | OpenShift / GitOps | cheapcloud, live-cicd |

## Labs (Wave 3+)

| Slug | Focus |
|------|-------|
| `activity-anon-cdp` | Anonymous visitor CDP on public sites |
| `surfing-r2-origin` | Surfing bytes off PVC → object/CDN |
| `demo-visitor-facade` | Cross-product fake mode pattern |

## Lane fill status

| Lane | Route | Gap to close |
|------|-------|--------------|
| Frontend | /projects/frontend | Link cards → hubs |
| Backend | /projects/backend | Add surfing-service, activity |
| AI/ML | /projects/ai-ml | Link interview-me hub |
| Cloud | /projects/cloud | Replace Coming soon with CheapCloud / OpenShift / demo |
| Infrastructure | /projects/infrastructure | Link etcd-synth, GitOps |
| Security | /projects/security | Thicken + OIDC topic |

## SEO surface checklist

- [x] robots.txt
- [x] sitemap.xml
- [x] Organization + Person JSON-LD
- [x] Per-route title / description / OG
- [x] FAQPage on hubs with FAQ
- [x] Answer-first article template
- [x] History mode routing (nginx try_files)
