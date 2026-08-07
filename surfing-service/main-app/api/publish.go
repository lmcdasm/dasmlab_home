package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type publishResult struct {
	DayID     string   `json:"day_id"`
	Published int      `json:"published"`
	Skipped   int      `json:"skipped"`
	Failed    int      `json:"failed"`
	Remaining int      `json:"remaining,omitempty"`
	Cleanup   int      `json:"pvc_cleaned,omitempty"`
	Errors    []string `json:"errors,omitempty"`
	Day       DayEntry `json:"day"`
}

// PublishDay promotes draft R2 objects or copies legacy PVC bytes to R2.
// Does NOT hold storeMu across network I/O — lock per item, persist every few successes.
// Query: cleanup_pvc=true, limit=N (chunk size, default 25, 0=all remaining).
func PublishDay(c *gin.Context) {
	dayID := c.Param("id")
	cleanup := strings.EqualFold(c.Query("cleanup_pvc"), "true") || c.Query("cleanup_pvc") == "1"
	limit := 25
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}

	if !mediaObjectStore.Enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "object store not configured (set SURFING_OBJECT_STORE=r2)"})
		return
	}

	storeMu.RLock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.RUnlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}
	// Snapshot IDs that still need publish work.
	pending := make([]string, 0, len(day.Media))
	for _, item := range day.Media {
		if needsPublish(item) {
			pending = append(pending, item.ID)
		}
	}
	storeMu.RUnlock()

	result := publishResult{DayID: dayID}
	done := 0
	for _, mediaID := range pending {
		if limit > 0 && done >= limit {
			break
		}

		storeMu.RLock()
		day, ok = dayStore[dayID]
		if !ok {
			storeMu.RUnlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
			return
		}
		idx := -1
		var item MediaItem
		for i := range day.Media {
			if day.Media[i].ID == mediaID {
				idx = i
				item = day.Media[i]
				break
			}
		}
		storeMu.RUnlock()
		if idx < 0 {
			result.Skipped++
			continue
		}
		if !needsPublish(item) {
			result.Skipped++
			continue
		}
		if item.Hidden {
			result.Skipped++
			continue
		}

		ext := normalizeExt(filepath.Ext(item.Filename))
		ct := mimeFromExt(ext, item.MediaType)
		var (
			key       string
			publicURL string
			err       error
			pvcPath   string
		)

		switch {
		case item.Origin == "link" || (item.ExternalURL != "" && item.ObjectKey == "" && item.MediaType == "other"):
			storeMu.Lock()
			if d, i, ok := findMediaLocked(dayID, mediaID); ok && i >= 0 {
				d.Media[i].Published = true
				d.Media[i].URL = d.Media[i].ExternalURL
				d.Published = dayPublished(d)
				dayStore[dayID] = d
				result.Published++
				done++
			}
			storeMu.Unlock()
			continue

		case item.Origin == "r2-draft" || strings.Contains(item.ObjectKey, "/draft/"):
			key, publicURL, err = promoteDraftObject(dayID, item.ID, ext, ct)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, item.ID+": promote "+err.Error())
				continue
			}

		default:
			pvcPath = resolveMediaPath(item.ID, ext)
			data, readErr := os.ReadFile(pvcPath)
			if readErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, item.ID+": read "+readErr.Error())
				continue
			}
			key, publicURL, err = putMediaObject(dayID, item.ID, ext, ct, data)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, item.ID+": put "+err.Error())
				continue
			}
		}

		if publicURL == "" {
			result.Failed++
			result.Errors = append(result.Errors, item.ID+": empty public URL (set R2_PUBLIC_BASE_URL)")
			continue
		}

		storeMu.Lock()
		if d, i, ok := findMediaLocked(dayID, mediaID); ok && i >= 0 {
			d.Media[i].URL = publicURL
			d.Media[i].Published = true
			d.Media[i].Origin = "r2"
			d.Media[i].ObjectKey = key
			d.Published = dayPublished(d)
			dayStore[dayID] = d
			result.Published++
			done++
		}
		storeMu.Unlock()

		if cleanup && pvcPath != "" {
			if err := os.Remove(pvcPath); err == nil {
				result.Cleanup++
			}
		}

		// Persist often so a timeout/restart never loses a large batch of R2 puts.
		if result.Published%5 == 0 {
			if err := persistManifest(); err != nil {
				log.Warnf("PublishDay: mid-persist failed: %v", err)
			}
		}
	}

	if err := persistManifest(); err != nil {
		log.Warnf("PublishDay: persist failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "published bytes but failed to persist manifest", "result": result})
		return
	}

	storeMu.RLock()
	day = dayStore[dayID]
	remaining := 0
	for _, item := range day.Media {
		if needsPublish(item) {
			remaining++
		}
	}
	result.Day = day
	result.Remaining = remaining
	storeMu.RUnlock()

	status := http.StatusOK
	if result.Failed > 0 && result.Published == 0 {
		status = http.StatusBadGateway
	} else if result.Failed > 0 || remaining > 0 {
		status = http.StatusMultiStatus
	}
	c.JSON(status, result)
}

func needsPublish(item MediaItem) bool {
	if item.Hidden {
		return false
	}
	if item.Origin == "link" || (item.ExternalURL != "" && item.ObjectKey == "" && item.MediaType == "other") {
		return !item.Published
	}
	if item.Origin == "r2-draft" || strings.Contains(item.ObjectKey, "/draft/") {
		return true
	}
	if item.Published && strings.HasPrefix(item.URL, "http") {
		return false
	}
	return true
}

func resolveMediaPath(mediaID, ext string) string {
	path := mediaFilePath(mediaID, ext)
	if _, err := os.Stat(path); err == nil {
		return path
	}
	alt := mediaFilePath(mediaID, strings.ToLower(ext))
	if _, err := os.Stat(alt); err == nil {
		return alt
	}
	return path
}
