package api

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PresignMediaUpload reserves a draft MediaItem in the album index and returns a
// short-lived R2 PUT URL. Browser uploads bytes direct to CDN — cluster only keeps the index.
// POST /days/:id/media/presign  JSON: { filename, content_type, size? }
func PresignMediaUpload(c *gin.Context) {
	dayID := c.Param("id")
	if !mediaObjectStore.Enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "direct upload unavailable (object store not configured)",
			"fallback": "multipart",
		})
		return
	}

	var req struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Size        int64  `json:"size"`
		Caption     string `json:"caption"`
		Kind        string `json:"kind"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	filename := strings.TrimSpace(req.Filename)
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename required"})
		return
	}
	contentType := strings.TrimSpace(req.ContentType)
	mediaType, ext := detectMediaTypeFromName(filename, contentType)
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type (photos and videos only)"})
		return
	}
	if contentType == "" {
		contentType = mimeFromExt(ext, mediaType)
	}
	if req.Size > maxUploadBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file exceeds upload limit"})
		return
	}

	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	mediaID := uuid.NewString()
	key := mediaDraftObjectKey(dayID, mediaID, ext)
	item := MediaItem{
		ID:        mediaID,
		Filename:  filename,
		MediaType: mediaType,
		Kind:      strings.TrimSpace(req.Kind),
		Caption:   strings.TrimSpace(req.Caption),
		URL:       mediaObjectStore.PublicURL(key), // may 404 until PUT completes
		CreatedAt: time.Now().UTC(),
		Published: false,
		Origin:    "r2-draft",
		ObjectKey: key,
	}
	if item.Kind == "" {
		if mediaType == "video" {
			item.Kind = KindVideo
		} else {
			item.Kind = KindPhoto
		}
	}
	normalizeMediaKind(&item)
	day.Media = append(day.Media, item)
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("PresignMediaUpload: persist failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist draft index"})
		return
	}

	ctx := c.Request.Context()
	putURL, err := mediaObjectStore.PresignPut(ctx, key, contentType, 20*time.Minute)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not mint upload URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"media":         item,
		"upload_url":    putURL,
		"upload_method": "PUT",
		"headers": gin.H{
			"Content-Type": contentType,
		},
		"object_key": key,
		"public_url": item.URL,
		"expires_in": 1200,
		"mode":       "direct-r2-draft",
	})
}

// CompleteMediaUpload confirms bytes landed on R2 for a draft item.
// POST /days/:id/media/:mediaId/complete
func CompleteMediaUpload(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")

	storeMu.Lock()
	day, idx, ok := findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	item := day.Media[idx]
	key := item.ObjectKey
	if key == "" {
		ext := normalizeExt(filepath.Ext(item.Filename))
		key = mediaDraftObjectKey(dayID, mediaID, ext)
	}
	storeMu.Unlock()

	ctx := c.Request.Context()
	size, exists, err := mediaObjectStore.Head(ctx, key)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not verify upload"})
		return
	}
	if !exists {
		c.JSON(http.StatusConflict, gin.H{"error": "object not found on CDN yet — retry complete shortly"})
		return
	}

	storeMu.Lock()
	day, idx, ok = findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	item = day.Media[idx]
	item.ObjectKey = key
	item.Origin = "r2-draft"
	item.Published = false
	if pub := mediaObjectStore.PublicURL(key); pub != "" {
		item.URL = pub
	}
	day.Media[idx] = item
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"media": item,
		"size":  size,
		"ready": true,
	})
}

func detectMediaTypeFromName(filename, contentType string) (mediaType, ext string) {
	ext = normalizeExt(filepath.Ext(filename))
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "image/"):
		return "image", extOr(ext, ".jpg")
	case strings.HasPrefix(ct, "video/"):
		return "video", extOr(ext, ".mp4")
	}
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic", ".heif":
		return "image", ext
	case ".mp4", ".mov", ".m4v", ".webm", ".avi", ".mkv":
		return "video", ext
	}
	return "", ext
}

func extOr(ext, fallback string) string {
	if ext == "" || ext == "." {
		return fallback
	}
	return ext
}
