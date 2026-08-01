package api

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	cookieState   = "surf_oauth_state"
	cookieSession = "surf_session"
	stateTTL      = 10 * time.Minute
	sessionTTL    = 24 * time.Hour
	roleAdmin     = "admin"
)

type oidcConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AppPublicURL string
	Enabled      bool
}

type AuthUser struct {
	Sub               string   `json:"sub"`
	PreferredUsername string   `json:"preferred_username"`
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	Roles             []string `json:"roles"`
	IsAdmin           bool     `json:"is_admin"`
}

type oidcService struct {
	cfg      oidcConfig
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth    oauth2.Config
	mu       sync.Mutex
	states   map[string]time.Time
}

var authSvc *oidcService

func oidcConfigFromEnv() oidcConfig {
	issuer := strings.TrimSpace(os.Getenv("OIDC_ISSUER"))
	if issuer == "" {
		base := strings.TrimRight(os.Getenv("KEYCLOAK_URL"), "/")
		realm := envOrDefault("KEYCLOAK_REALM", "dasmlab")
		if base != "" {
			issuer = base + "/realms/" + realm
		}
	}
	appURL := strings.TrimRight(os.Getenv("APP_PUBLIC_URL"), "/")
	redirect := strings.TrimSpace(os.Getenv("OIDC_REDIRECT_URI"))
	if redirect == "" && appURL != "" {
		redirect = appURL + "/api/surfing/auth/callback"
	}
	cfg := oidcConfig{
		Issuer:       issuer,
		ClientID:     envOrDefault("OIDC_CLIENT_ID", "dasmlab-home"),
		ClientSecret: strings.TrimSpace(os.Getenv("OIDC_CLIENT_SECRET")),
		RedirectURL:  redirect,
		AppPublicURL: appURL,
	}
	cfg.Enabled = cfg.Issuer != "" && cfg.ClientID != "" && cfg.ClientSecret != "" && cfg.RedirectURL != ""
	return cfg
}

func oidcHTTPClient() *http.Client {
	// OpenShift ingress uses a cluster-local CA not in public trust stores.
	skip := strings.EqualFold(os.Getenv("OIDC_INSECURE_SKIP_VERIFY"), "true") ||
		os.Getenv("OIDC_INSECURE_SKIP_VERIFY") == "1"
	if !skip {
		return http.DefaultClient
	}
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // lab/cluster ingress CA
		},
	}
}

func initAuth() error {
	cfg := oidcConfigFromEnv()
	svc := &oidcService{cfg: cfg, states: map[string]time.Time{}}
	if !cfg.Enabled {
		authSvc = svc
		log.Info("OIDC: disabled (set KEYCLOAK_URL + OIDC_CLIENT_SECRET + APP_PUBLIC_URL)")
		return nil
	}
	ctx := oidc.ClientContext(context.Background(), oidcHTTPClient())
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return fmt.Errorf("oidc provider: %w", err)
	}
	svc.provider = provider
	svc.verifier = provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})
	svc.oauth = oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	authSvc = svc
	log.Infof("OIDC: enabled issuer=%s client=%s", cfg.Issuer, cfg.ClientID)
	return nil
}

func (s *oidcService) enabled() bool { return s != nil && s.cfg.Enabled }

func AuthConfig(c *gin.Context) {
	if authSvc == nil {
		c.JSON(http.StatusOK, gin.H{"enabled": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"enabled":        authSvc.enabled(),
		"issuer":         authSvc.cfg.Issuer,
		"client_id":      authSvc.cfg.ClientID,
		"redirect_uri":   authSvc.cfg.RedirectURL,
		"app_public_url": authSvc.cfg.AppPublicURL,
		"login_path":     "/auth/login",
		"logout_path":    "/auth/logout",
	})
}

func AuthLogin(c *gin.Context) {
	if authSvc == nil || !authSvc.enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "OIDC not configured"})
		return
	}
	state, err := randomString(24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "state"})
		return
	}
	authSvc.putState(state)
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
	c.SetCookie(cookieState, state, int(stateTTL.Seconds()), "/", "", secure, true)
	c.Redirect(http.StatusFound, authSvc.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline))
}

func AuthCallback(c *gin.Context) {
	if authSvc == nil || !authSvc.enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "OIDC not configured"})
		return
	}
	state := c.Query("state")
	cookie, _ := c.Cookie(cookieState)
	if state == "" || cookie == "" || state != cookie || !authSvc.takeState(state) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}
	tok, err := authSvc.oauth.Exchange(oidc.ClientContext(c.Request.Context(), oidcHTTPClient()), code)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "token exchange failed"})
		return
	}
	rawID, ok := tok.Extra("id_token").(string)
	if !ok || rawID == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "missing id_token"})
		return
	}
	idTok, err := authSvc.verifier.Verify(c.Request.Context(), rawID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid id_token"})
		return
	}
	var claims map[string]any
	_ = idTok.Claims(&claims)
	user := userFromClaims(claims, authSvc.cfg.ClientID)
	payload, _ := json.Marshal(map[string]any{
		"user":         user,
		"access_token": tok.AccessToken,
		"expiry":       tok.Expiry.UTC().Format(time.RFC3339),
	})
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
	c.SetCookie(cookieSession, base64.RawURLEncoding.EncodeToString(payload), int(sessionTTL.Seconds()), "/", "", secure, true)
	c.SetCookie(cookieState, "", -1, "/", "", secure, true)
	dest := authSvc.cfg.AppPublicURL
	if dest == "" {
		dest = "/"
	}
	c.Redirect(http.StatusFound, dest+"/#/surfing")
}

func AuthLogout(c *gin.Context) {
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
	c.SetCookie(cookieSession, "", -1, "/", "", secure, true)
	dest := "/"
	if authSvc != nil && authSvc.cfg.AppPublicURL != "" {
		dest = authSvc.cfg.AppPublicURL + "/#/surfing"
	}
	if authSvc != nil && authSvc.enabled() {
		end := strings.TrimRight(authSvc.cfg.Issuer, "/") + "/protocol/openid-connect/logout"
		c.Redirect(http.StatusFound, end+"?post_logout_redirect_uri="+authSvc.cfg.AppPublicURL+"%2F%23%2Fsurfing&client_id="+authSvc.cfg.ClientID)
		return
	}
	c.Redirect(http.StatusFound, dest)
}

func AuthMe(c *gin.Context) {
	user, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"authenticated": false, "oidc_enabled": authSvc != nil && authSvc.enabled()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authenticated": true, "oidc_enabled": true, "user": user})
}

func currentUser(c *gin.Context) (AuthUser, bool) {
	raw, err := c.Cookie(cookieSession)
	if err != nil || raw == "" {
		return AuthUser{}, false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return AuthUser{}, false
	}
	var wrap struct {
		User AuthUser `json:"user"`
	}
	if err := json.Unmarshal(decoded, &wrap); err != nil {
		return AuthUser{}, false
	}
	if wrap.User.Sub == "" {
		return AuthUser{}, false
	}
	return wrap.User, true
}

func isOwner(c *gin.Context) bool {
	u, ok := currentUser(c)
	return ok && u.IsAdmin
}

func userFromClaims(claims map[string]any, clientID string) AuthUser {
	u := AuthUser{
		Sub:               strClaim(claims, "sub"),
		PreferredUsername: strClaim(claims, "preferred_username"),
		Email:             strClaim(claims, "email"),
		Name:              strClaim(claims, "name"),
	}
	if ra, ok := claims["resource_access"].(map[string]any); ok {
		if client, ok := ra[clientID].(map[string]any); ok {
			if roles, ok := client["roles"].([]any); ok {
				for _, r := range roles {
					if s, ok := r.(string); ok {
						u.Roles = append(u.Roles, s)
						if s == roleAdmin {
							u.IsAdmin = true
						}
					}
				}
			}
		}
	}
	return u
}

func strClaim(m map[string]any, k string) string {
	if v, ok := m[k].(string); ok {
		return v
	}
	return ""
}

func (s *oidcService) putState(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = time.Now().Add(stateTTL)
}

func (s *oidcService) takeState(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.states[state]
	delete(s.states, state)
	return ok && time.Now().Before(exp)
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// RequireOwner returns 401/403 for mutating owner-only endpoints when OIDC is on.
func RequireOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authSvc == nil || !authSvc.enabled() {
			// Dev mode without OIDC — allow mutations (local/cluster bootstrap).
			c.Next()
			return
		}
		u, ok := currentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sign in required"})
			return
		}
		if !u.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role required"})
			return
		}
		c.Next()
	}
}
