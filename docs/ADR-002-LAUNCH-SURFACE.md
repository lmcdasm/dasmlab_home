# ADR-002: DASMLAB 2.0 launch surface (`/launch`)

**Status:** Accepted  
**Date:** 2026-08-08  
**Site:** dasmlab.org (dasmlab_home)  
**Pairs with:** Observatory ADR-0402 (Campaign orchestration)

## Context

The first DPO-orchestrated campaign (`dasmlab-2.0-launch`) needs a public slash page that is the **web_slash** channel artifact: brand-first, SEO/GEO friendly, measurable via Activity UTMs and DPO collectors.

## Decision

1. Canonical route: **`/launch`**. Alias: **`/2.0` → `/launch`** (redirect).
2. One composition hero: brand signal, one headline, one supporting sentence, CTA group (Explore / Observatory / Contact). No card soup in the first viewport.
3. Default UTM when shared from campaign: `utm_source=<channel>&utm_medium=campaign&utm_campaign=dasmlab-2.0-launch`.
4. Campaign id cross-link: page meta / footer references DPO campaign `dasmlab-2.0-launch`.
5. After ship: DPO TRACK-HOME baselines `pre-launch` → `post-launch`.

## Non-goals

- Surfing album page or Share tollgate as the launch surface
- Replacing IndexPage as the site home

## Consequences

- New page + routes in `src/router/routes.js`
- ShareSheet / campaign renderers point at `https://dasmlab.org/launch`
- Observatory story doc cites this ADR
