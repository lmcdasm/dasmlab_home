package api

import (
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var tagNameRe = regexp.MustCompile(`^[\p{L}\p{M}][\p{L}\p{M}\s'.-]{0,39}$`)

func normalizeTagName(raw string) (string, error) {
	name := strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
	if name == "" {
		return "", errTagEmpty
	}
	if utf8.RuneCountInString(name) > 40 {
		return "", errTagTooLong
	}
	// Hard ban URLs / @handles / emails — name text only.
	lower := strings.ToLower(name)
	if strings.Contains(lower, "http://") || strings.Contains(lower, "https://") ||
		strings.Contains(lower, "www.") || strings.Contains(name, "/") ||
		strings.Contains(name, "@") || strings.Contains(name, ".com") {
		return "", errTagNotName
	}
	if !tagNameRe.MatchString(name) {
		return "", errTagNotName
	}
	return name, nil
}

var (
	errTagEmpty   = errString("name is required")
	errTagTooLong = errString("name must be ≤ 40 characters")
	errTagNotName = errString("tags are plain names only — no links, handles, or URLs")
)

type errString string

func (e errString) Error() string { return string(e) }

func findMediaLocked(dayID, mediaID string) (DayEntry, int, bool) {
	day, ok := dayStore[dayID]
	if !ok {
		return DayEntry{}, -1, false
	}
	for i := range day.Media {
		if day.Media[i].ID == mediaID {
			return day, i, true
		}
	}
	return day, -1, false
}

// RecordMediaPlay increments play_count (in-app muted play or tollgate resolve).
func RecordMediaPlay(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")

	storeMu.Lock()
	day, idx, ok := findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	day.Media[idx].PlayCount++
	count := day.Media[idx].PlayCount
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("RecordMediaPlay: persist failed: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"play_count": count})
}

// ProposeMediaTag lets allowed viewers submit a plain name tag (owner must approve).
func ProposeMediaTag(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	name, err := normalizeTagName(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeMu.Lock()
	day, idx, ok := findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	item := day.Media[idx]
	for _, t := range item.Tags {
		if strings.EqualFold(t.Name, name) && t.Status != "rejected" {
			storeMu.Unlock()
			c.JSON(http.StatusConflict, gin.H{"error": "name already tagged on this media", "tag": t})
			return
		}
	}
	tag := MediaTag{
		ID:        uuid.NewString(),
		Name:      name,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
	}
	item.Tags = append(item.Tags, tag)
	day.Media[idx] = item
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist tag"})
		return
	}
	c.JSON(http.StatusCreated, tag)
}

// ModerateMediaTag owner-approves or rejects a pending name tag.
func ModerateMediaTag(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")
	tagID := c.Param("tagId")
	action := strings.ToLower(strings.TrimSpace(c.Param("action")))
	if action != "approve" && action != "reject" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action must be approve or reject"})
		return
	}

	storeMu.Lock()
	day, idx, ok := findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	item := day.Media[idx]
	found := -1
	for i := range item.Tags {
		if item.Tags[i].ID == tagID {
			found = i
			break
		}
	}
	if found < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	if action == "approve" {
		item.Tags[found].Status = "approved"
	} else {
		item.Tags[found].Status = "rejected"
	}
	tag := item.Tags[found]
	day.Media[idx] = item
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist tag"})
		return
	}
	c.JSON(http.StatusOK, tag)
}
