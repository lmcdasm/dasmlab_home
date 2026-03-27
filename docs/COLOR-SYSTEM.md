# DASMLAB color system (0.6.0)

Derived from the three existing brand colors plus neutrals and semantics. We can expand to a full palette later; this is the working set for the 2026 Spring FE update.

---

## Seed colors (existing)

| Token    | Hex       | Role                          |
|----------|-----------|-------------------------------|
| Primary  | `#3C6B5E` | Teal — main brand, headers, CTAs |
| Secondary| `#647550` | Olive — support, secondary UI |
| Accent   | `#A46CAD` | Purple — highlights, links, emphasis |

---

## Derived palette

### Brand (expand usage across the site)

- **Primary** `#3C6B5E` — Keep. Use for: toolbar, key buttons, primary headings, links.
- **Primary light** `#5A8B7E` — Lighter teal for hovers, backgrounds, badges.
- **Primary dark** `#2D5248` — Darker teal for active states, contrast on light.
- **Secondary** `#647550` — Keep. Use for: secondary buttons, cards, borders.
- **Secondary light** `#7D8F6A` — Softer olive for backgrounds.
- **Accent** `#A46CAD` — Keep. Use for: accents, special highlights, one key word per section.
- **Accent light** `#C090C9` — Softer purple for hovers, gradients.

### Neutrals (dark theme baseline)

- **Dark** `#1D1D1D` — Main dark surface (existing).
- **Dark page** `#121212` — Page background (existing).
- **Dark elevated** `#2A2A2A` — Cards, elevated surfaces.
- **Grey 8** — Muted text (Quasar default).
- **Grey 6** — Borders, dividers.

### Neutrals (light theme, if used)

- **Light bg** `#FAFAFA`
- **Light surface** `#FFFFFF`
- **Grey 8** — Body text.

### Semantic (keep, minor tweaks optional)

- **Positive** `#56bA6D` — Success, go.
- **Negative** `#CC4757` — Error, destructive.
- **Info** `#B8CDD1` — Info, subtle.
- **Warning** `#F0D792` — Warning, caution.

---

## Usage rules (Pillar 5 — “explosion of color”)

1. **Use more than one brand color per view** — e.g. hero with primary, cards with secondary borders, one accent highlight.
2. **Reserve accent for emphasis** — One key word, one CTA, or one decorative element per section.
3. **Backgrounds** — Use primary-light or secondary-light as tinted backgrounds in sections, not only grey.
4. **Consistency** — Same role (e.g. “primary CTA”) uses the same token everywhere.

---

## Implementation

- **Quasar:** Update `src/css/quasar.variables.scss` with the new tokens (primary-light, primary-dark, secondary-light, accent-light, dark-elevated). Quasar’s CSS vars will expose them.
- **CSS custom properties:** Optionally add a small `_color-tokens.scss` or section in `app.scss` that maps these to `--dasm-*` for use in custom components.
- **Next:** When we add proprietary effects, use these tokens so motion and overlays stay on-brand.

---

## Doc info

- **Created:** 2026-03-11  
- **Branch:** `2026-spring-FE-update`  
- **Target:** 0.6.0
