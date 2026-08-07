# Target map: 2.0.0 — Engineering Discovery Platform

**Branch (proposed):** `2026-dasmlab-2.0`  
**Baseline:** 1.0.x (Lab Constellation + Surfing + Activity Phase A)  
**Target release:** **2.0.0**  
**Brand:** Technologies DASMLAB Inc. · https://dasmlab.org

---

## Vision

Not a blog. An **Engineering Knowledge Network** with two audiences:

1. **Visitors** — arrive via search / AI answers; get answer-first pages.
2. **Engineers** — follow every answer into *how we built it*, source/CE, and related tech.

Private products stay IP-protected behind Keycloak (or edge auth) for live ops, while **demo / unknown visitors** can explore labeled fake modes that never mutate production.

Open core: libraries Apache 2.0; app CE MPL 2.0 + commercial dual license; hosted SaaS stays private.

---

## Three tracks

| Track | Outcome |
|-------|---------|
| **A — Home EKN** | Project hubs, topic hubs, labs, SEO/GEO spine, filled lanes, real About media |
| **B — Demo visitors** | Shared contract; interview-me → mock-me → cheapcloud → etcd → camera |
| **C — Open core** | Dual-license boilerplate; first CE extract; org layout |

See also:

- [DEMO-VISITOR-CONTRACT.md](./DEMO-VISITOR-CONTRACT.md)
- [HUB-INVENTORY.md](./HUB-INVENTORY.md)
- [licenses/](./licenses/) — Apache / MPL / CC BY / commercial templates
- [ADR-000 Observatory Platform](./ADR-000-OBSERVATORY-PLATFORM.md) — company DOP/DPO architecture (sibling repo `dasmlab-observatory-platform`)
- [ADR-001 CDN-mgr GEO](./ADR-001-CDN-MGR-GEO-ENGAGEMENT.md)

Supersedes constellation notes in TARGET-0.7.0 for product direction; keep 0.7 for Surfing CDN ops history.

---

## Information architecture

```text
dasmlab.org/
  about/
  projects/{lane}/
  projects/{slug}/          # hub: overview | architecture | how-we-built | tech | source | faq
  topics/{tech}/
  labs/{slug}/
  surfing/
  activity/                 # owner-only
```

Every project hub ends with **Behind the Design** (stack → topic hubs).

---

## Version

- `package.json` → **2.0.0**
- Sibling apps bump major when they ship demo facade and/or CE split.

---

## Success criteria

- Labeled demos on interview-me and mock-me with zero live side effects
- No primary-lane “Coming soon”; About uses intentional media; ≥5 project hubs + FAQ schema
- Sitemap + Organization/Person JSON-LD; answer-first template on hubs
- ≥1 CE repo under dual MPL + commercial docs
- Home tagged 2.0.0

## Non-goals

- Publishing production platforms / secrets / Keycloak ops
- Replacing Keycloak for operators
- Next/React rewrite
- Blocking on Activity Phases B–D
