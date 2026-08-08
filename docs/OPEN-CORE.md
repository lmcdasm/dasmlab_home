# Open-core & dual licensing — Technologies DASMLAB Inc.

**GitHub organization (locked for 2.0):** `github.com/dasmlab`  
(Personal mirrors may remain under `lmcdasm` during migration.)

## License matrix

| Asset | License | Location |
|-------|---------|----------|
| Shared libraries (`dasmlab-go`, `dasmlab-ui`) | Apache 2.0 | CE / public |
| App Community Editions (`*-ce`) | MPL 2.0 + [COMMERCIAL.md](./licenses/COMMERCIAL.md) | CE / public |
| Docs / tutorials on dasmlab.org | CC BY 4.0 | Site |
| Hosted SaaS, secrets, prod GitOps, Keycloak ops | Proprietary | Never published |

Templates live in [docs/licenses/](./licenses/).

## First CE extracts (scaffolded in-repo)

Until separate GitHub repos are created, scaffolds live under:

```text
open-core/
  README.md
  dasmlab-go/          # Apache 2.0 stub (demo middleware notes)
  mock-me-ce/          # MPL + commercial dual offer
  interview-me-ce/     # MPL + commercial dual offer
  cheapcloud-ce/       # MPL + commercial dual offer
```

Publish by pushing each folder to `github.com/dasmlab/<name>` with the included LICENSE / COMMERCIAL / CITATION files. Do **not** include production secrets or full SaaS trees.

## Activity Phase A (home)

Anonymous public `POST /activity` on surfing-service is part of the 2.0 home engagement spine. Promote the preview FE that includes Phase A tracker to production when cutting the 2.0 release tag (`dasmlab-home` **2.0.0**). Surfing API Phase A should already be on the shared prod API if preview verified anonymous writes.
