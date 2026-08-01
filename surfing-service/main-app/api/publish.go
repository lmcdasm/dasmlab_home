package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type publishResult struct {
	DayID     string `json:"day_id"`
	Published int    `json:"published"`
	Skipped   int    `json:"skipped"`
	Failed    int    `json:"failed"`
	Cleanup   int    `json:"pvc_cleaned,omitempty"`
	Errors    []string `json:"errors,omitempty"`
	Day       DayEntry `json:"day"`
}

// PublishDay copies unpublished PVC media for an album into R2 under
// surfing/albums/{dayId}/original/… and rewrites CDN URLs in the manifest.
// Query: cleanup_pvc=true removes local bytes after a successful put.
func PublishDay(c *gin.Context) {
	dayID := c.Param("id")
	cleanup := strings.EqualFold(c.Query("cleanup_pvc"), "true") || c.Query("cleanup_pvc") == "1"

	if !mediaObjectStore.Enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object store not configured (set SURFING_OBJECT_STORE=r2)"})
		return
	}

	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	result := publishResult{DayID: dayID}
	media := make([]MediaItem, len(day.Media))
	copy(media, day.Media)

	for i := range media {
		item := &media[i]
		if item.Hidden {
			result.Skipped++
			continue
		}
		// Link-only shares already live on the open web — nothing to put on R2.
		if item.Origin == "link" || (item.ExternalURL != "" && item.ObjectKey == "" && item.MediaType == "other") {
			if !item.Published {
				item.Published = true
				item.URL = item.ExternalURL
				result.Published++
			} else {
				result.Skipped++
			}
			continue
		}
		if item.Published && strings.HasPrefix(item.URL, "http") {
			result.Skipped++
			continue
		}
		ext := normalizeExt(filepath.Ext(item.Filename))
		path := resolveMediaPath(item.ID, ext)
		data, err := os.ReadFile(path)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, item.ID+": read "+err.Error())
			continue
		}
		key, publicURL, err := putMediaObject(dayID, item.ID, ext, mimeFromExt(ext, item.MediaType), data)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, item.ID+": put "+err.Error())
			continue
		}
		if publicURL == "" {
			result.Failed++
			result.Errors = append(result.Errors, item.ID+": empty public URL (set R2_PUBLIC_BASE_URL)")
			continue
		}
		item.URL = publicURL
		item.Published = true
		item.Origin = "r2"
		item.ObjectKey = key
		result.Published++

		if cleanup {
			if err := os.Remove(path); err == nil {
				result.Cleanup++
			}
		}
	}

	day.Media = media
	day.Published = dayPublished(day)
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("PublishDay: persist failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "published bytes but failed to persist manifest", "result": result})
		return
	}

	result.Day = day
	status := http.StatusOK
	if result.Failed > 0 && result.Published == 0 {
		status = http.StatusBadGateway
	} else if result.Failed > 0 {
		status = http.StatusMultiStatus
	}
	c.JSON(status, result)
}

func resolveMediaPath(mediaID, ext string) string {
	path := mediaFilePath(mediaID, ext)
	if _, err := os.Stat(path); err == nil {
		return path
	}
	// Case variants from imports (.JPG / .MOV).
	base := filepath.Join(mediaBasePath, mediaID)
	candidates := []string{
		base + strings.ToUpper(ext),
		base + strings.ToLower(ext),
		base + ".JPG", base + ".jpg",
		base + ".JPEG", base + ".jpeg",
		base + ".MOV", base + ".mov",
		base + ".MP4", base + ".mp4",
		base + ".HEIC", base + ".heic",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return path
}
