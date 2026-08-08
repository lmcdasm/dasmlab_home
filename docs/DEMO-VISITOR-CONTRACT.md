# Demo / Unknown Visitor Contract

**Canonical for:** Technologies DASMLAB Inc. products (mock-me, interview-me, cheapcloud, etcd-synth, camera-scrape, and future apps).  
**Home site:** https://dasmlab.org — already public; links to product demos with UTM + Activity `track` events.

---

## Personas

| Persona | Entry | Can do | Cannot do |
|---------|-------|--------|-----------|
| **Unknown / demo** | Public “Try demo” CTA; cookie `*_demo` (no Keycloak) | Browse UI, run **scripted fake flows**, read synthetic fixtures | Mutate prod, deploy live nodes, stress real etcd, touch customer/PII data, admin APIs |
| **Authenticated** | Keycloak (or product invite guest) | Product capabilities per role | — |
| **Owner / admin** | Keycloak + roles | Full control + Activity panel | — |

---

## Hard rules

1. Demo is **opt-in** and **loudly labeled**: banner text like `Demo / fake mode — not a live system`.
2. Backend **rejects mutate** routes unless the session is real auth (or an explicit invite guest scoped to one resource).
3. Demo handlers return **fixtures / simulated step timelines only** — never call live deploy, cloud credentials, or production stress paths.
4. Activity / engagement events from demo sessions SHOULD set `demo: true` (or equivalent props) so operators can filter.
5. Cookie names: `{productPrefix}_demo` (e.g. `mm_demo`, `im_demo`). HttpOnly where set by API; `Secure` + `SameSite=Lax` in prod.

---

## Middleware pattern (Go / Gin or chi)

```text
RequireAuth          → real OIDC session required
AllowDemoRead        → auth OR valid demo cookie for GET/read
DenyDemoMutate       → if demo cookie only → 403 on POST/PUT/PATCH/DELETE (except /demo/* simulate)
```

Demo simulate endpoints live under `/api/.../demo/` (or `/demo/`) and are the **only** write-shaped routes demos may hit.

---

## FE pattern (Vue)

- Composable `useDemoMode()` → `{ isDemo, enterDemo, exitDemo }`
- Persistent banner component while `isDemo`
- Disable or rewire destructive CTAs to call `/demo/simulate`
- Link from dasmlab.org project hubs: `?utm_source=dasmlab&utm_medium=hub&utm_campaign=demo`

---

## Product notes

| Product | Demo shape |
|---------|------------|
| **interview-me** | Public synthetic interview session (no PII write); reuse guest cookie patterns |
| **mock-me** | Scripted orchestration timeline; zero live node deploys |
| **cheapcloud** | Readonly recommend / envelope UI over fixtures; no cloud credentials |
| **etcd-synth** | UI browse + dry-run / sandbox only; real stress operator-gated |
| **camera-scrape** | Identity boundary + public sample gallery vs private scrapes |

---

## Rollout order

1. interview-me (reference)  
2. mock-me  
3. cheapcloud  
4. etcd-synth  
5. camera-scrape  

Optional later: extract shared `dasmlab-go-auth` (Apache 2.0) once two+ apps share middleware — not a Wave 0 blocker.
