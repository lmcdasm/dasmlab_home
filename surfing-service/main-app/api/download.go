package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// DownloadMedia is the DASMLAB-gated download tollgate.
// Visibility: public | private (owner) | group (owner or client role group/member).
// Always prefer this over raw CDN links so cheapcloud/cdn-mgr can meter + encrypt later.
func DownloadMedia(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")

	storeMu.RLock()
	day, idx, ok := findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.RUnlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	item := day.Media[idx]
	storeMu.RUnlock()

	if item.Hidden && !isOwner(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	vis := item.DownloadVisibility
	if vis == "" {
		vis = "public"
	}
	owner := isOwner(c)
	group := false
	if u, ok := currentUser(c); ok {
		group = hasGroupRole(u)
	}
	if !downloadAllowed(vis, owner, group) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":               "download not allowed for your access level",
			"download_visibility": vis,
		})
		return
	}

	// Prefer absolute CDN / external URL when present.
	target := strings.TrimSpace(item.ExternalURL)
	if target == "" {
		target = strings.TrimSpace(item.URL)
	}
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		c.Header("Content-Disposition", `attachment; filename="`+sanitizeFilename(item.Filename)+`"`)
		c.Redirect(http.StatusFound, target)
		return
	}

	// PVC /serve fallback
	ext := filepath.Ext(item.Filename)
	path := mediaFilePath(item.ID, ext)
	if _, err := os.Stat(path); err != nil {
		alt := mediaFilePath(item.ID, strings.ToLower(ext))
		if _, err2 := os.Stat(alt); err2 != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		path = alt
	}
	c.Header("Content-Disposition", `attachment; filename="`+sanitizeFilename(item.Filename)+`"`)
	c.File(path)
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "download"
	}
	name = strings.ReplaceAll(name, `"`, "")
	name = strings.ReplaceAll(name, "/", "_")
	return name
}
