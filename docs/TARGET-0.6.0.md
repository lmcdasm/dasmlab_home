# Target map: 0.6.0 — 2026 Spring FE update

**Branch:** `2026-spring-FE-update`  
**Baseline:** 0.5.x (latest tag: `0.5.1-alpha`)  
**Target release:** **0.6.0** (minor bump)  
**Vision:** Move DASMLAB.org from “clean but tutorial” to a portfolio that feels **purposefully crafted and human** — distinctive, ownable, and recognizably DASMLAB in a sea of algorithmic sameness.

---

## Guiding principle

> In a world of algorithmic sameness, human craft is becoming the differentiator.

We are not chasing tooling (Vue/Quasar stays). We are investing in **design and content**: proprietary effects, cohesive visual system, restraint in copy, and interactions that signal intention. The goal is something **unique, organic, and real** — not generic AI-generated animation.

---

## Current state (baseline)

- **Stack:** Vue 3, Quasar 2, Vite, Vue Router (no Pinia in use).
- **Brand tokens (today):** Primary `#3C6B5E`, Secondary `#647550`, Accent `#A46CAD`; dark theme in use.
- **Gaps:** Tutorial/default look; limited “jazz”; no unified visual language; copy and layout not yet aligned with 2026 trends below.

---

## 0.6.0 pillars (aligned with 2026 web design trends)

Use these as a checklist to iterate on. Each pillar has a **DASMLAB-specific** angle.

| # | Trend | DASMLAB goal | Status |
|---|--------|---------------|--------|
| 1 | **Proprietary effects & styles** | Define and implement a small set of **ownable** effects (e.g. signature motion, cursor/scroll reactions, custom filters or overlays) that feel recognizably DASMLAB and scale across the site. |
| 2 | **Art + advanced UI** | Where it fits, blend hand-crafted or illustrative touches with the app UI (e.g. hero, project cards, About) so the product feels like a “gallery piece,” not template output. |
| 3 | **Minimalism in copy** | Strip copy to essentials. Let layout and visuals carry meaning. Short, confident headlines and CTAs; cut everything that isn’t load-bearing. |
| 4 | **TL;DR experience** | Offer a quick-scan overview (pitch-deck style or structured summary) so visitors get “what we do” fast, then choose where to go deeper (e.g. nav or hero). |
| 5 | **Explosion of color** | Evolve from single accent to a **full DASMLAB color system** (palette + usage rules). Use multiple hues deliberately across backgrounds, type, and UI for energy and memorability. |
| 6 | **Dynamic text treatments** | Use typography and light motion (e.g. reveal, emphasis, subtle animation) to give key copy weight and make reading feel intentional. |
| 7 | **Guided scrolling** | Add scroll wayfinding: progress, cues, or light interactivity so users know where they are and what’s next. Combat “shrink attention” with clarity, not clutter. |
| 8 | **Infinite canvas** | Use metaphors of possibility (grids, nodes, open space, dot fields) where it fits the narrative — e.g. “where creation happens” for lab/portfolio. |

---


## Priority groupings (from Q&A)

| Priority | Pillars |
|----------|---------|
| **PRIO 1** | Proprietary effects, Color system |
| **PRIO 1/2** | Art / hand-crafted feel |
| **PRIO 2** | Copy + TL;DR |
| **PRIO 3** | Dynamic text, Guided scrolling, Infinite canvas |

First sprint: Color system (see [COLOR-SYSTEM.md](COLOR-SYSTEM.md)), Proprietary effects (1–2), start Art/hand-crafted. Then Copy + TL;DR. See [VOICE-AND-PRIORITIES.md](VOICE-AND-PRIORITIES.md) for one-liner and vibe.

---

## Approach switcher (live-switch UX)

Build a **reusable "approach" switcher** so we can live-switch between different presentation modes in the same frontend:

- **Mechanism:** Dropdown or toggle (toolbar or drawer) to select current approach.
- **Persistence:** e.g. `localStorage`.
- **First use case:** Home page — **Hero + block** (current: WhatsNew + Carousel + cards) vs **Nav-level TL;DR** (condensed summary at top or in nav; then dive deeper).
- **Later:** Reuse for other A/B-style layout or component variants.

---

## Sources & exploration (how to create “neat things”)

To avoid generic results, we need references and techniques that support **craft** and **custom** feel. Use this section as a living list.

### Generative / code-driven effects (to study and adapt)

- **Canvas / WebGL:** Custom shaders, particle systems, and procedural backgrounds (we already use tsparticles; consider extending or replacing with brand-specific behavior).
- **CSS:** `@keyframes`, `clip-path`, `mask`, `mix-blend-mode`, custom properties for theme-aware motion.
- **Vue/Quasar:** `@vueuse/motion`, transitions, and scroll-driven state (e.g. `useScroll`, `useInView`) to drive reveals and progress.
- **Typography:** Variable fonts, gradient text, staggered character/word animation, “typewriter” or reveal effects used sparingly on hero or key lines.

### References and inspiration (curate, don’t copy)

- **Sites mentioned in the article:** Springboards (custom hero animations), Anthropic (illustrative system), Pencil.dev (ASCII/classical + UI), Rootly (UI as gallery), Ruul / Sandbar (minimal copy), Flabbergast (TL;DR nav), Tesoro (full color system), Habito Studio / Purpose Talent (dynamic text), Dropbox Dash McLaren (scroll speedometer), Emons (numbered scroll steps), Flim (grid-as-canvas).
- **Communities:** Webflow showcases, Awwwards, CSS-Tricks, Codrops — for techniques we can reimplement in Vue/Quasar.
- **DASMLAB-specific:** Decide what “lab” and “portfolio” mean visually (nodes, grids, experiments, tools) and translate that into one or two recurring motifs.

### Tools and libraries (optional, only if they serve the system)

- **Vue/Quasar-first:** VueUse (motion, scroll, gesture), existing tsparticles — tune for brand.
- **CSS-first:** Custom SCSS variables and mixins for the new color system and proprietary effects.
- **Asset creation:** Illustration or texture work (even simple) that feels hand-made to pair with UI (trend #2).

---

## 0.6.0 scope (iterative)

- **In scope for 0.6.0:**  
  - Visual and copy pass aligned with the 8 pillars.  
  - DASMLAB color system (tokens + usage).  
  - At least one or two proprietary effects or motifs used consistently.  
  - TL;DR or overview experience (e.g. hero + nav).  
  - Guided scroll and/or dynamic text where it adds clarity.

- **Out of scope for 0.6.0:**  
  - Framework migration (e.g. Qwik).  
  - Full redesign of backend or APIs.  
  - New product features unrelated to “portfolio + presence.”

---

## Next steps

1. ~~**Color system:**~~ Draft done → [COLOR-SYSTEM.md](COLOR-SYSTEM.md). **Next:** Apply tokens in `quasar.variables.scss` and use across components.
2. **Approach switcher:** Implement composable + toolbar dropdown + IndexPage reacting to `hero-block` vs `nav-tldr` (see above).
3. **Copy audit:** List every user-facing string; apply minimalism and TL;DR (see [VOICE-AND-PRIORITIES.md](VOICE-AND-PRIORITIES.md)).
4. **Pick 2–3 "signature" effects:** e.g. hero motion, scroll progress, text treatment — implement in one place, then scale.
5. **Revisit this doc** after each sprint; update pillar status and "Sources & exploration" as we find what works.

---

## Doc info

- **Created:** 2026-03-11  
- **Branch:** `2026-spring-FE-update`  
- **Target release:** 0.6.0  
- **License:** Same as project (AAL)
