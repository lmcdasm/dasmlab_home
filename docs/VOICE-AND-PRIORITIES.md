# DASMLAB voice & priorities (0.6.0)

Captured from Q&A for the 2026 Spring FE update. Use this for copy passes and to keep the vibe consistent across multi-sites and content types.

---

## One-liner / positioning

**DASMLAB** is a **living lab and portfolio** — a place where technology gets poked, prodded, and brought to life. From raw prototypes to polished tools, every project is part experiment, part education: demystifying the build process and sharing “how I did it” so others can learn, remix, and iterate.

*(Source: AboutPage, refined for reuse.)*

---

## Vibe words

- **Studio** · **Portfolio** · **Workshop** · **Playground**

Multi-sites and many project types live under one umbrella. The challenge: **one consistent vibe and theme**, with different “styles of view” by content type:

| Content type        | Example style              |
|---------------------|----------------------------|
| Reference           | Clear, scannable, durable  |
| POC / MVP           | Experimental, “under build”|
| Knowledge sharing  | Educational, how-I-did-it  |
| Fun / playground    | Playful, exposition        |

---

## Voice notes (from existing pages)

- **About:** Living lab, poked/prodded, brought to life, experiment + education, demystify, share, learn/remix/iterate. Personal: decades in tech, telecom → automotive → AI, 60+ countries, Outaouais/Ottawa, Fernie, Laurentians. Beyond code: snowboard, windsurf, garden, house music, pottery.
- **Contact:** Short. “Reach out,” email CTA.
- **Home:** Cards by category (Frontend, Backend, AI/ML, Cloud, Infra, Security); Design Carousel as live content. (WhatsNew retired — see TARGET-0.7.0.)

---

## Priority groupings (0.6.0)

| Priority | Pillars |
|----------|---------|
| **PRIO 1** | Proprietary effects, Color system |
| **PRIO 1/2** | Art / hand-crafted feel |
| **PRIO 2** | Copy + TL;DR |
| **PRIO 3** | Dynamic text, Guided scrolling, Infinite canvas |

First sprint: Color system (done in doc; next: code), Proprietary effects (pick 1–2), start Art/hand-crafted. Then Copy + TL;DR pass. Then the rest.

---

## Approach switcher (live-switch UX)

We want to **live-switch between different “approaches”** in the same frontend (e.g. hero + block vs nav-level TL;DR). Build a simple, reusable pattern:

- **Mechanism:** Dropdown or toggle (e.g. in toolbar or drawer) that selects the current “approach.”
- **Persistence:** e.g. `localStorage` so the choice survives refresh.
- **First use case:** Two home approaches:
  - **Hero + block** — Current: Carousel + project cards (constellation map planned in 0.7.0).
  - **Nav-level TL;DR** — Condensed summary (pitch-deck style or compact strip) at top or in nav; user can then dive into full home or sections.
- **Later:** Reuse the same switcher for other A/B-style component or layout variants.

---

## Doc info

- **Created:** 2026-03-11  
- **Branch:** `2026-spring-FE-update`  
- **Target:** 0.6.0
