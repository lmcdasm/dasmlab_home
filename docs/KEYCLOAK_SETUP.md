# Keycloak SSO — dasmlab-home / Surfing (v1.0.0)

| Item | Value |
|------|-------|
| Keycloak | `https://keycloak.apps.2026-prod-1.ocp.dasmlab.org` |
| Realm | `dasmlab` |
| Client ID | `dasmlab-home` |
| App | `https://dasmlab.org` (+ **explicit** preview hosts) |
| Callback (prod) | `https://dasmlab.org/api/surfing/auth/callback` |

## Client setup (same shape as mock-me / interview-me)

1. Clients → Create → OpenID Connect → **`dasmlab-home`**
2. Client authentication **ON**, Standard flow **ON**
3. Redirect URIs — Keycloak does **not** honor `https://dev-*-…` host wildcards.
   Add an **explicit** host per developer preview (same as mock-me):
   - `https://dasmlab.org/api/surfing/auth/callback`
   - `https://dasmlab.org/*`
   - `https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org/api/surfing/auth/callback`
   - `https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org/*`
   - (repeat for each `dev-<owner>-dasmlab-home…` host)
4. Web origins: matching prod + each preview origin
5. Post-logout redirect URIs: `https://dasmlab.org/*`, `https://dasmlab.org`, and the same for each preview origin.
   Logout uses `post_logout_redirect_uri={origin}/` (no `#/…` fragment — Keycloak rejects those).
6. Copy client secret → K8s secret `surfing-service-system/dasmlab-home-oidc` key `client-secret`
7. Client role **`admin`** → assign to your user (Filter by clients → dasmlab-home)
8. Optional client role **`group`** (or `member`) → test “group can download”
9. Ensure **roles** client scope / User Client Role mapper so `resource_access["dasmlab-home"].roles` includes `admin` (and `group` when testing)
10. **Add to ID token** AND **Add to access token** must both be ON for that mapper
11. After role changes: **Sign out + Sign in** (roles are snapshotted into `surf_session` at login)

Automation: `scripts/ci/ensure-keycloak-preview-uris.sh` (called from `deploy-preview.sh`)
adds the preview host’s redirect + web-origin using the cluster Keycloak initial-admin secret.

If the ID-token mapper is incomplete, surfing-service still grants owner when:
- `admin` appears on the **access** token or **realm_access**, or
- `preferred_username` is listed in env `OIDC_OWNER_USERNAMES` (default / prod: `dasm`)

## Preview login (shared surfing)

Surfing is one shared API. Login from a preview FE uses that host as OIDC
`redirect_uri` / cookie scope (`Host` / `X-Forwarded-Host`), not a hard-coded
`dasmlab.org` callback. Env `OIDC_REDIRECT_URI` remains the **default** for
prod/`APP_PUBLIC_URL`; allowlisted origins:

- `https://dasmlab.org`
- `https://dev-<owner>-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org`
- optional extras via `OIDC_ALLOWED_ORIGINS`

## Behaviour

| Actor | Sees |
|-------|------|
| Public (not signed in) | Non-hidden media only; public notes; approved tags; **public** downloads |
| Owner (`admin` role) | Hidden items **grayed**; private notes; all download levels; approve tags; publish / curate / AI curate |
| Group (`group` / `member` role) | Same as public + **group** downloads |

### Download access

| `download_visibility` | Who gets Download |
|-----------------------|-------------------|
| `public` | Anyone via gated `GET /days/:id/media/:mediaId/download` |
| `private` | Owner (`admin`) only |
| `group` | Owner + users with client role `group` / `member` |

Raw `object_key` is stripped from list responses. Playback CDN sniffing is acknowledged until signed URLs / bucket+infra encrypt (cdn-mgr + cheapcloud).

Mutating APIs require owner when OIDC is enabled. Without OIDC secret, service stays open for bootstrap (dev).

## Related

- Share tollgate: `docs/SHARE-TOLLGATE.md`
- Value prop: `docs/VALUE-PROP-PERSONAL-CDN.md`
- interview-me twin: `interview-me/docs/KEYCLOAK_SETUP.md`
