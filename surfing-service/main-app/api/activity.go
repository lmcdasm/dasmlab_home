package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"surfing-service/activity"

	"github.com/gin-gonic/gin"
)

const cookieActLogin = "surf_act_login"

var activityStore *activity.Store

func initActivityStore(dataDir string) error {
	s, err := activity.NewStore(dataDir)
	if err != nil {
		return err
	}
	activityStore = s
	log.Infof("Activity: log at %s/activity/events.jsonl viewers=%v", dataDir, activityViewers())
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
// Machine access: Authorization: Bearer $ACTIVITY_MACHINE_TOKEN (for DPO / internal collectors).
func RequireActivityViewer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tok := strings.TrimSpace(os.Getenv("ACTIVITY_MACHINE_TOKEN")); tok != "" {
			authz := strings.TrimSpace(c.GetHeader("Authorization"))
			if strings.HasPrefix(authz, "Bearer ") && strings.TrimSpace(strings.TrimPrefix(authz, "Bearer ")) == tok {
				c.Next()
				return
			}
		}
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

func recordLoginEvent(u AuthUser) {
	if activityStore == nil || u.Sub == "" {
		return
	}
	_ = activityStore.Append(activity.Event{
		Type:  activity.TypeLogin,
		User:  u.PreferredUsername,
		Sub:   u.Sub,
		Email: u.Email,
	})
}

// emitLoginOnce records a login at most once per browser SSO session.
func emitLoginOnce(c *gin.Context, u AuthUser) {
	if _, err := c.Cookie(cookieActLogin); err == nil {
		return
	}
	recordLoginEvent(u)
	setAuthCookie(c, cookieActLogin, "1", int(sessionTTL.Seconds()))
}

// PostActivity appends navigate/engaged events for the signed-in user.
func PostActivity(c *gin.Context) {
	if activityStore == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "activity not configured"})
		return
	}
	u, ok := currentUser(c)
	if !ok {
		// When OIDC off, accept anonymous lab posts under "dev".
		if authSvc != nil && authSvc.enabled() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "sign in required"})
			return
		}
		u = AuthUser{PreferredUsername: "dev"}
	}
	var body struct {
		Type      string `json:"type"`
		Path      string `json:"path"`
		DwellMs   int64  `json:"dwellMs"`
		VisibleMs int64  `json:"visibleMs"`
		EngagedMs int64  `json:"engagedMs"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	typ := strings.TrimSpace(body.Type)
	if typ == "" {
		typ = activity.TypeNavigate
	}
	if typ != activity.TypeNavigate && typ != activity.TypeEngaged {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be navigate or engaged"})
		return
	}
	if err := activityStore.Append(activity.Event{
		Type:      typ,
		User:      u.PreferredUsername,
		Sub:       u.Sub,
		Email:     u.Email,
		Path:      strings.TrimSpace(body.Path),
		DwellMs:   body.DwellMs,
		VisibleMs: body.VisibleMs,
		EngagedMs: body.EngagedMs,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
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
