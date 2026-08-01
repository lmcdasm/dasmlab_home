# Keycloak SSO — dasmlab-home / Surfing (v1.0.0)

| Item | Value |
|------|-------|
| Keycloak | `https://keycloak.apps.2026-prod-1.ocp.dasmlab.org` |
| Realm | `dasmlab` |
| Client ID | `dasmlab-home` |
| App | `https://dasmlab.org` (+ preview hosts) |
| Callback | `https://dasmlab.org/api/surfing/auth/callback` |

## Client setup (same shape as interview-me)

1. Clients → Create → OpenID Connect → **`dasmlab-home`**
2. Client authentication **ON**, Standard flow **ON**
3. Redirect URIs:
   - `https://dasmlab.org/api/surfing/auth/callback`
   - `https://dasmlab.org/*`
   - `https://dev-*-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org/api/surfing/auth/callback`
   - `https://dev-*-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org/*`
4. Web origins: matching prod + preview
5. Copy client secret → K8s secret `surfing-service-system/dasmlab-home-oidc` key `client-secret`
6. Client role **`admin`** → assign to your user (Filter by clients → dasmlab-home)
7. Ensure **roles** client scope / User Client Role mapper so `resource_access["dasmlab-home"].roles` includes `admin`

## Behaviour

| Actor | Sees |
|-------|------|
| Public (not signed in) | Non-hidden media only; public notes; approved tags |
| Owner (`admin` role) | Hidden items **grayed**; private notes; approve tags; publish / curate / AI curate |

Mutating APIs require owner when OIDC is enabled. Without OIDC secret, service stays open for bootstrap (dev).

## Related

- Share tollgate: `docs/SHARE-TOLLGATE.md`
- Value prop: `docs/VALUE-PROP-PERSONAL-CDN.md`
- interview-me twin: `interview-me/docs/KEYCLOAK_SETUP.md`
