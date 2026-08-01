package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UnhideMedia clears the soft-hide flag (owner only).
func UnhideMedia(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")

	storeMu.Lock()
	day, idx, ok := findMediaLocked(dayID, mediaID)
	if !ok || idx < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}
	day.Media[idx].Hidden = false
	item := day.Media[idx]
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// PatchDay updates album-level policies (tag_policy).
func PatchDay(c *gin.Context) {
	dayID := c.Param("id")
	var req struct {
		TagPolicy *string `json:"tag_policy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}
	if req.TagPolicy != nil {
		v := strings.ToLower(strings.TrimSpace(*req.TagPolicy))
		switch v {
		case "public", "group", "off":
			day.TagPolicy = v
		default:
			storeMu.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "tag_policy must be public|group|off"})
			return
		}
	}
	dayStore[dayID] = day
	storeMu.Unlock()
	if err := persistManifest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist"})
		return
	}
	c.JSON(http.StatusOK, day)
}

// CuratePublish is "Curate a Publish" for one asset or the whole album (owner).
// Today it triggers the existing R2 publish path; compression/transcode hooks later.
func CuratePublish(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := strings.TrimSpace(c.Query("media_id"))
	var req struct {
		MediaID     string `json:"media_id"`
		Compress    bool   `json:"compress"`
		NotesPublic bool   `json:"notes_public"`
	}
	_ = c.ShouldBindJSON(&req)
	if mediaID == "" {
		mediaID = strings.TrimSpace(req.MediaID)
	}

	if mediaID != "" && req.NotesPublic {
		storeMu.Lock()
		if day, idx, ok := findMediaLocked(dayID, mediaID); ok && idx >= 0 {
			day.Media[idx].NotesVisibility = "public"
			day.Media[idx].Hidden = false
			dayStore[dayID] = day
		}
		storeMu.Unlock()
		_ = persistManifest()
	}

	if req.Compress {
		log.Infof("CuratePublish: compress requested for day=%s media=%s (queued for transcoder — not yet online)", dayID, mediaID)
	}

	PublishDay(c)
}

// AICurate scaffolds signed-in AI curation against a provider (OpenAI/cheapcloud farm).
func AICurate(c *gin.Context) {
	dayID := c.Param("id")
	var req struct {
		Prompt    string   `json:"prompt"`
		MediaIDs  []string `json:"media_ids"`
		Action    string   `json:"action"` // share | save | iterate
		Provider  string   `json:"provider"`
	}
	_ = c.ShouldBindJSON(&req)
	if strings.TrimSpace(req.Prompt) == "" {
		req.Prompt = "Curate a short highlight reel description for this album."
	}
	cfg := loadAIConfig()
	out := gin.H{
		"day_id":     dayID,
		"action":     firstNonEmpty(req.Action, "save"),
		"provider":   firstNonEmpty(req.Provider, cfg.Provider),
		"status":     "accepted",
		"message":    "AI curate job accepted — result is a CDN-linked object you control (not a Meta upload).",
		"media_ids":  req.MediaIDs,
		"prompt":     req.Prompt,
		"next_steps": []string{"share tollgate", "save to album folder", "iterate prompt"},
	}
	if !cfg.Enabled() {
		out["status"] = "deferred"
		out["warning"] = "AI provider not configured — wire SURFING_AI_* or cheapcloud farm"
		c.JSON(http.StatusAccepted, out)
		return
	}
	// Lightweight: reuse style brief path as a curation note (full reel pipeline later).
	storeMu.RLock()
	day, ok := dayStore[dayID]
	storeMu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}
	samples := sampleThemePhotos(day, 3)
	brief, err := aiStyleBrief(c.Request.Context(), cfg, day, "Windsurfing Trips", req.Prompt, samples)
	if err != nil {
		out["status"] = "partial"
		out["warning"] = err.Error()
		c.JSON(http.StatusAccepted, out)
		return
	}
	out["curation"] = brief
	c.JSON(http.StatusOK, out)
}

// TranscodeMedia stub — on-the-fly / job compression for .MOV etc.
func TranscodeMedia(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "queued",
		"message": "Transcode/compress adapter not online yet — cdn-mgr will run spot/vLLM or ffmpeg jobs per realm.",
		"day_id":  c.Param("id"),
		"media":   c.Param("mediaId"),
		"hint":    "Prefer publishing a compressed derivative to R2; originals stay in draft/original/.",
	})
}
