
# API Security Rationale (`API_SECURITY.md`)

## Overview

As we transition the `whatsnew-service` from a private, cluster-only API (`whatsnew-service.whatsnew-service-system.svc.cluster.local`) to a publicly accessible endpoint (`whatsnew-service.dasmlab.org`), it becomes imperative to secure the service against unauthorized access, scraping, and abuse—while supporting a modern, flexible client experience.

This document examines candidate API security approaches and explains our rationale for selecting **OAuth2/OpenID Connect (OIDC)** as our authentication and authorization model.

---

## Security Options Considered

### 1. **API Keys**
- **Description:** Simple shared secret (key) issued to each client; sent in header or query.
- **Pros:** Very easy to implement and use.
- **Cons:** Poor for real user authentication, hard to rotate securely, easy to leak (especially in browser apps), no built-in user context or permissioning.

### 2. **Basic HTTP Auth**
- **Description:** Username/password pair in each request.
- **Pros:** Built-in to HTTP; trivial setup.
- **Cons:** Not user-friendly for browsers or modern APIs, insecure if not over HTTPS, not scalable, not suited for delegated access or SSO.

### 3. **JWT Bearer Tokens**
- **Description:** Stateless signed tokens (often used with OAuth2/OIDC).
- **Pros:** Portable, can encode claims/roles, widely supported in Go.
- **Cons:** Need a secure token issuance and revocation process; typically a component of a bigger system.

### 4. **Third-party API Gateways**
- **Description:** External proxies (Kong, Ambassador, etc.) that sit in front of services.
- **Pros:** Built-in auth, rate limiting, and logging; drop-in for multi-service environments.
- **Cons:** More infrastructure/operational overhead, less flexibility for bespoke auth logic.

### 5. **OAuth2 / OpenID Connect (OIDC)**
- **Description:** Industry standard for delegated, federated, and user- or machine-authenticated access. Integrates with SSO, Google, Auth0, Okta, and self-hosted IdPs (e.g., Keycloak, ORY Hydra).
- **Pros:** 
  - Modern, highly secure, widely adopted.
  - Supports both user and machine auth.
  - Enables fine-grained permissions (scopes/claims).
  - Revocable, rotatable, supports SSO.
  - Developer ecosystem is mature, with well-maintained libraries for Go (Gin, etc.).
- **Cons:** 
  - More complex to set up (requires IdP configuration and token validation logic).
  - Clients need to handle OAuth2 flows.

---

## Why We Chose OAuth2/OIDC

After evaluating the above options, **OAuth2 with OpenID Connect** stands out as the best-fit solution for our requirements:

- **Scalability:** Supports growth from a small internal API to a public or partner-facing service without rearchitecture.
- **Security:** Access tokens are short-lived, rotatable, and can encode permissions, preventing key leakage and privilege escalation.
- **User and Service Access:** Works for both user-facing (browser, mobile) and machine-to-machine integrations.
- **Standards Compliance:** Industry standard; supported by all major identity providers.
- **Ecosystem:** Strong middleware and library support for Go/Gin, including token validation and middleware enforcement.
- **SSO/Identity Federation:** Lets us leverage existing identity systems (Google, GitHub, enterprise IdP) for authentication—no need for "new passwords."
- **Future-Proofing:** Eases the addition of new clients, external partners, and federated login as we grow.

---

## Design Pattern

### High-level Flow:

1. **User or Service Requests Access:**  
   - User: Initiates OAuth2 login (Authorization Code flow).
   - Service: Obtains a token via the Client Credentials flow.

2. **Token Issued by Identity Provider (IdP):**  
   - Could be Google, Auth0, Keycloak, etc.

3. **Client Calls our Public API:**  
   - Presents access token (JWT) in the `Authorization: Bearer ...` header.

4. **Gin Middleware Verifies Token:**  
   - Token is checked for validity, expiry, and permissions.
   - Access granted or denied accordingly.

5. **Optional:** Rate limiting, logging, and further fine-grained authorization applied at the application or gateway layer.

---

## Example Gin Middleware (Token Validation)

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/coreos/go-oidc/v3/oidc"
    "net/http"
    "strings"
)

func AuthMiddleware(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        rawIDToken := strings.TrimPrefix(authHeader, "Bearer ")
        _, err := verifier.Verify(c.Request.Context(), rawIDToken)
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Next()
    }
}
```

---

## Implementation Plan

1. **Select/Configure Identity Provider:**  
   Use Google, Auth0, Okta, or a self-hosted IdP (Keycloak, ORY Hydra).

2. **Protect All Non-public Endpoints:**  
   Require a valid OAuth2 token for every sensitive or modifying endpoint (read-only GETs may be made public, as needed).

3. **Frontend Integration:**  
   - For browser-based clients: Use Authorization Code (PKCE) flow.
   - For headless clients/services: Use Client Credentials flow.
   - Document how to obtain and refresh tokens.

4. **Rate Limiting & Abuse Detection:**  
   Add as middleware or at ingress if needed.

5. **Documentation:**  
   Publicly document OAuth2/OIDC flows, required scopes, and how to obtain credentials.

---

## Comparison Table

| Option      | UX            | Security | Effort | Rotatable | Suits Public APIs? | Future-proof | Notes                   |
|-------------|---------------|----------|--------|-----------|--------------------|--------------|-------------------------|
| API Key     | Manual        | Low-Med  | Low    | No        | No                 | Poor         | Ok for internal/dev     |
| Basic Auth  | Browser Pop   | Low      | Low    | No        | No                 | Poor         | Not modern              |
| JWT Only    | Token-based   | Med      | Med    | Yes       | Sometimes          | Medium       | Use with IdP            |
| Gateway     | Any           | High     | High   | Yes       | Yes                | Good         | For large orgs          |
| **OAuth2/OIDC** | SSO/Modern| High     | Med    | Yes       | Yes                | Excellent    | Industry best practice  |

---

## Conclusion

**OAuth2 with OIDC is the preferred, modern, and secure pattern for opening up our APIs to public access.**  
This ensures we can serve both web and machine clients, leverage SSO, and scale securely as we grow. The Go ecosystem (Gin, coreos/go-oidc) makes implementation straightforward and maintainable.

**See `README.md` for next steps and developer onboarding.**

