package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"surfing-service/activity"

	"github.com/gin-gonic/gin"
)

const (
	cookieActLogin = "surf_act_login"
	cookieAnonID   = "surf_aid"
	cookieSessID   = "surf_sid"

	anonIDTTL    = 365 * 24 * time.Hour
	sessionIdle  = 30 * time.Minute
	activityRPM  = 90 // per client IP
)

var activityStore *activity.Store

var (
	rateMu   sync.Mutex
	rateHits = map[string][]time.Time{}
)

func initActivityStore(dataDir string) error {
	s, err := activity.NewStore(dataDir)
	if err != nil {
		return err
	}
	activityStore = s
	log.Infof("Activity: log at %s/activity/events.jsonl viewers=%v public_write=true", dataDir, activityViewers())
	return nil
}

func activityViewers() []string {
	raw := strings.TrimSpace(os.Getenv("ACTIVITY_VIEWERS"))
	if raw == "" {
		return []string{"dasm"}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"dasm"}
	}
	return out
}

func canViewActivity(u AuthUser) bool {
	if !authSvc.enabled() {
		return true
	}
	uname := strings.TrimSpace(u.PreferredUsername)
	for _, v := range activityViewers() {
		if uname == v {
			return true
		}
	}
	return false
}

// RequireAuth requires a signed-in session (any user).
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authSvc == nil || !authSvc.enabled() {
			c.Next()
			return
		}
		if _, ok := currentUser(c); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sign in required"})
			return
		}
		c.Next()
	}
}

// RequireActivityViewer requires owner/admin AND preferred_username on ACTIVITY_VIEWERS (default dasm).
func RequireActivityViewer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authSvc == nil || !authSvc.enabled() {
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
		if !canViewActivity(u) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  "forbidden",
				"detail": "activity log is restricted to ACTIVITY_VIEWERS (default: dasm)",
			})
			return
		}
		c.Next()
	}
}

func newOpaqueID(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func lookLikeBot(ua string) bool {
	u := strings.ToLower(ua)
	for _, needle := range []string{
		"bot", "crawler", "spider", "slurp", "curl/", "wget", "python-requests",
		"go-http-client", "httpclient", "scrapy", "headless", "phantom", "selenium",
	} {
		if strings.Contains(u, needle) {
			return true
		}
	}
	return ua == ""
}

func activityRateLimited(c *gin.Context) bool {
	ip := c.ClientIP()
	if ip == "" {
		ip = "unknown"
	}
	now := time.Now()
	cutoff := now.Add(-time.Minute)
	rateMu.Lock()
	defer rateMu.Unlock()
	hits := rateHits[ip]
	kept := hits[:0]
	for _, t := range hits {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= activityRPM {
		rateHits[ip] = kept
		return true
	}
	rateHits[ip] = append(kept, now)
	return false
}

// ensureVisitorIDs mints/refreshes first-party anonymousId + sessionId cookies.
func ensureVisitorIDs(c *gin.Context) (anonID, sessID string) {
	anonID, _ = c.Cookie(cookieAnonID)
	if anonID == "" || len(anonID) < 8 {
		anonID = newOpaqueID(16)
	}
	setAuthCookie(c, cookieAnonID, anonID, int(anonIDTTL.Seconds()))

	sessID, _ = c.Cookie(cookieSessID)
	if sessID == "" || len(sessID) < 8 {
		sessID = newOpaqueID(12)
	}
	// Sliding 30m idle window (RudderStack-style default).
	setAuthCookie(c, cookieSessID, sessID, int(sessionIdle.Seconds()))
	return anonID, sessID
}

func enrichFromRequest(c *gin.Context) (ua, locale, country string, bot bool) {
	ua = strings.TrimSpace(c.GetHeader("User-Agent"))
	if len(ua) > 240 {
		ua = ua[:240]
	}
	locale = strings.TrimSpace(c.GetHeader("Accept-Language"))
	if i := strings.IndexByte(locale, ','); i >= 0 {
		locale = strings.TrimSpace(locale[:i])
	}
	if len(locale) > 32 {
		locale = locale[:32]
	}
	country = strings.TrimSpace(c.GetHeader("CF-IPCountry"))
	if country == "" {
		country = strings.TrimSpace(c.GetHeader("X-Country-Code"))
	}
	country = strings.ToUpper(country)
	if country == "XX" || country == "T1" {
		country = ""
	}
	bot = lookLikeBot(ua)
	return
}

func recordLoginEvent(c *gin.Context, u AuthUser) {
	if activityStore == nil || u.Sub == "" {
		return
	}
	anonID, sessID := ensureVisitorIDs(c)
	ua, locale, country, bot := enrichFromRequest(c)
	_ = activityStore.Append(activity.Event{
		Type:        activity.TypeLogin,
		User:        u.PreferredUsername,
		Sub:         u.Sub,
		Email:       u.Email,
		AnonymousID: anonID,
		SessionID:   sessID,
		UA:          ua,
		Locale:      locale,
		Country:     country,
		Bot:         bot,
	})
	_ = activityStore.Append(activity.Event{
		Type:        activity.TypeIdentify,
		User:        u.PreferredUsername,
		Sub:         u.Sub,
		Email:       u.Email,
		AnonymousID: anonID,
		SessionID:   sessID,
		UA:          ua,
		Locale:      locale,
		Country:     country,
	})
	if anonID != "" {
		_ = activityStore.Append(activity.Event{
			Type:        activity.TypeAlias,
			User:        u.PreferredUsername,
			Sub:         u.Sub,
			AnonymousID: anonID,
			PreviousID:  anonID,
			SessionID:   sessID,
		})
	}
}

// emitLoginOnce records login/identify/alias at most once per browser SSO session.
func emitLoginOnce(c *gin.Context, u AuthUser) {
	if _, err := c.Cookie(cookieActLogin); err == nil {
		return
	}
	recordLoginEvent(c, u)
	setAuthCookie(c, cookieActLogin, "1", int(sessionTTL.Seconds()))
}

// PostActivity appends page/navigate/engaged/track for any visitor (anon or known).
func PostActivity(c *gin.Context) {
	if activityStore == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "activity not configured"})
		return
	}
	if activityRateLimited(c) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limited"})
		return
	}

	var body struct {
		Type        string `json:"type"`
		Event       string `json:"event"`
		Path        string `json:"path"`
		Title       string `json:"title"`
		Referrer    string `json:"referrer"`
		UTMSource   string `json:"utmSource"`
		UTMMedium   string `json:"utmMedium"`
		UTMCampaign string `json:"utmCampaign"`
		DwellMs     int64  `json:"dwellMs"`
		VisibleMs   int64  `json:"visibleMs"`
		EngagedMs   int64  `json:"engagedMs"`
		ScrollMax   int    `json:"scrollMaxPct"`
		AnonymousID string `json:"anonymousId"`
		SessionID   string `json:"sessionId"`
		MessageID   string `json:"messageId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	typ := strings.TrimSpace(body.Type)
	if typ == "" {
		typ = activity.TypePage
	}
	// Normalize page ↔ navigate for older clients.
	if typ == activity.TypeNavigate {
		typ = activity.TypePage
	}
	switch typ {
	case activity.TypePage, activity.TypeEngaged, activity.TypeTrack:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be page, engaged, or track"})
		return
	}

	anonID, sessID := ensureVisitorIDs(c)
	if hint := strings.TrimSpace(body.AnonymousID); hint != "" && len(hint) >= 8 && anonID == "" {
		anonID = hint
	}
	if hint := strings.TrimSpace(body.SessionID); hint != "" && len(hint) >= 8 {
		// Prefer server cookie; accept body only to help first beacon before Set-Cookie round-trip.
		if _, err := c.Cookie(cookieSessID); err != nil {
			sessID = hint
			setAuthCookie(c, cookieSessID, sessID, int(sessionIdle.Seconds()))
		}
	}

	ua, locale, country, bot := enrichFromRequest(c)
	ev := activity.Event{
		Type:        typ,
		Event:       strings.TrimSpace(body.Event),
		AnonymousID: anonID,
		SessionID:   sessID,
		Path:        strings.TrimSpace(body.Path),
		Title:       strings.TrimSpace(body.Title),
		Referrer:    strings.TrimSpace(body.Referrer),
		UTMSource:   strings.TrimSpace(body.UTMSource),
		UTMMedium:   strings.TrimSpace(body.UTMMedium),
		UTMCampaign: strings.TrimSpace(body.UTMCampaign),
		DwellMs:     body.DwellMs,
		VisibleMs:   body.VisibleMs,
		EngagedMs:   body.EngagedMs,
		ScrollMax:   body.ScrollMax,
		UA:          ua,
		Locale:      locale,
		Country:     country,
		Bot:         bot,
		MessageID:   strings.TrimSpace(body.MessageID),
	}
	if u, ok := currentUser(c); ok {
		ev.User = u.PreferredUsername
		ev.Sub = u.Sub
		ev.Email = u.Email
	}

	if err := activityStore.Append(ev); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"ok":          true,
		"anonymousId": anonID,
		"sessionId":   sessID,
	})
}

// ListActivity returns recent events (newest first) — dual-gated to ACTIVITY_VIEWERS.
func ListActivity(c *gin.Context) {
	if activityStore == nil {
		c.JSON(http.StatusOK, gin.H{"events": []any{}})
		return
	}
	limit := 200
	if q := c.Query("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			limit = n
		}
	}
	events, err := activityStore.List(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if events == nil {
		events = []activity.Event{}
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}
