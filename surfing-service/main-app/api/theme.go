package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AI provider is OpenAI-compatible today (api.openai.com). Later point
// SURFING_AI_BASE_URL at cheapcloud's managed AI farm (vLLM on spots).

type aiConfig struct {
	Provider   string
	BaseURL    string
	APIKey     string
	ChatModel  string
	ImageModel string
}

func loadAIConfig() aiConfig {
	provider := strings.ToLower(strings.TrimSpace(firstNonEmpty(os.Getenv("SURFING_AI_PROVIDER"), "openai")))
	base := strings.TrimRight(strings.TrimSpace(firstNonEmpty(
		os.Getenv("SURFING_AI_BASE_URL"),
		"https://api.openai.com/v1",
	)), "/")
	return aiConfig{
		Provider:   provider,
		BaseURL:    base,
		APIKey:     strings.TrimSpace(firstNonEmpty(os.Getenv("SURFING_AI_API_KEY"), os.Getenv("OPENAI_API_KEY"))),
		ChatModel:  firstNonEmpty(os.Getenv("SURFING_AI_CHAT_MODEL"), "gpt-4o-mini"),
		ImageModel: firstNonEmpty(os.Getenv("SURFING_AI_IMAGE_MODEL"), "gpt-image-1"),
	}
}

func (c aiConfig) Enabled() bool {
	return c.APIKey != "" && c.BaseURL != ""
}

type generateThemeRequest struct {
	Prompt      string `json:"prompt"`
	SampleCount int    `json:"sample_count"`
	Sport       string `json:"sport"`
}

// GenerateTheme samples album photos, asks the AI farm for a style brief +
// banner/background art, stores art on R2, and attaches theme to the day.
func GenerateTheme(c *gin.Context) {
	dayID := c.Param("id")
	var req generateThemeRequest
	_ = c.ShouldBindJSON(&req)
	if req.SampleCount <= 0 {
		req.SampleCount = 3
	}
	if req.SampleCount > 6 {
		req.SampleCount = 6
	}

	storeMu.RLock()
	day, ok := dayStore[dayID]
	storeMu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	samples := sampleThemePhotos(day, req.SampleCount)
	sport := strings.TrimSpace(req.Sport)
	if sport == "" {
		sport = "Windsurfing Trips"
	}
	userPrompt := strings.TrimSpace(req.Prompt)
	if userPrompt == "" {
		userPrompt = fmt.Sprintf("%s — %s%s", day.Title, sport, locationSuffix(day.Location))
	}

	cfg := loadAIConfig()
	theme := DayTheme{
		Label:          day.Title,
		Prompt:         userPrompt,
		Primary:        "#0f8f7c",
		Secondary:      "#063642",
		Accent:         "#5eb4c8",
		Provider:       "local",
		SampleMediaIDs: sampleIDs(samples),
		GeneratedAt:    time.Now().UTC(),
	}

	if !cfg.Enabled() {
		theme.StyleBrief = "Local fallback theme (no SURFING_AI_API_KEY). Ocean teal wash tailored to album title."
		theme.Provider = "local"
		applyLocalThemeDefaults(&theme, day, sport)
		persistDayTheme(dayID, theme)
		c.JSON(http.StatusOK, gin.H{"theme": theme, "day": getDayCopy(dayID), "warning": "AI key unset — local CSS theme only"})
		return
	}

	brief, err := aiStyleBrief(c.Request.Context(), cfg, day, sport, userPrompt, samples)
	if err != nil {
		log.Warnf("GenerateTheme: style brief failed: %v", err)
		theme.StyleBrief = "Style brief unavailable — using prompt directly."
	} else {
		theme.StyleBrief = brief.Brief
		if brief.Primary != "" {
			theme.Primary = brief.Primary
		}
		if brief.Secondary != "" {
			theme.Secondary = brief.Secondary
		}
		if brief.Accent != "" {
			theme.Accent = brief.Accent
		}
		if brief.Label != "" {
			theme.Label = brief.Label
		}
	}

	bannerPrompt := buildBannerPrompt(day, sport, userPrompt, theme.StyleBrief)
	bannerBytes, bannerCT, err := aiGenerateImage(c.Request.Context(), cfg, bannerPrompt, imageSizeForModel(cfg.ImageModel))
	if err != nil {
		log.Warnf("GenerateTheme: banner image failed: %v", err)
		// Tailored fallback: use a sample photo from *this* album as the banner.
		if len(samples) > 0 && strings.HasPrefix(samples[0].URL, "http") {
			theme.BannerURL = samples[0].URL
			theme.BackgroundURL = samples[0].URL
			theme.StyleBrief += " | AI image unavailable — using album sample as banner"
		} else {
			applyLocalThemeDefaults(&theme, day, sport)
			theme.StyleBrief += " | image gen failed: " + err.Error()
		}
		theme.Provider = cfg.Provider
		persistDayTheme(dayID, theme)
		c.JSON(http.StatusOK, gin.H{"theme": theme, "day": getDayCopy(dayID), "warning": "image generation failed — sample/palette theme applied"})
		return
	}

	bannerKey, bannerURL, err := putThemeObject(dayID, "banner", extFromContentType(bannerCT), bannerCT, bannerBytes)
	if err != nil || bannerURL == "" {
		// Keep bytes on PVC as fallback serve path.
		localID := uuid.NewString()
		path := mediaFilePath(localID, extFromContentType(bannerCT))
		_ = os.WriteFile(path, bannerBytes, 0o664)
		bannerURL = "/serve?id=" + localID
		log.Warnf("GenerateTheme: R2 theme put failed (%v) — using PVC serve %s key=%s", err, bannerURL, bannerKey)
	}

	bgPrompt := buildBackgroundPrompt(day, sport, userPrompt, theme.StyleBrief)
	bgBytes, bgCT, bgErr := aiGenerateImage(c.Request.Context(), cfg, bgPrompt, imageSizeForModel(cfg.ImageModel))
	bgURL := bannerURL
	if bgErr == nil && len(bgBytes) > 0 {
		if _, u, err := putThemeObject(dayID, "background", extFromContentType(bgCT), bgCT, bgBytes); err == nil && u != "" {
			bgURL = u
		}
	} else if bgErr != nil {
		log.Warnf("GenerateTheme: background image failed (reusing banner): %v", bgErr)
	}

	theme.BannerURL = bannerURL
	theme.BackgroundURL = bgURL
	theme.Provider = cfg.Provider
	persistDayTheme(dayID, theme)

	c.JSON(http.StatusOK, gin.H{"theme": theme, "day": getDayCopy(dayID)})
}

func locationSuffix(loc string) string {
	loc = strings.TrimSpace(loc)
	if loc == "" {
		return ""
	}
	return " in " + loc
}

func sampleIDs(items []MediaItem) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.ID)
	}
	return out
}

func sampleThemePhotos(day DayEntry, n int) []MediaItem {
	out := make([]MediaItem, 0, n)
	for _, m := range day.Media {
		if m.Hidden {
			continue
		}
		normalizeMediaKind(&m)
		if m.Kind != KindPhoto {
			continue
		}
		if m.URL == "" {
			continue
		}
		out = append(out, m)
		if len(out) >= n {
			break
		}
	}
	return out
}

func applyLocalThemeDefaults(theme *DayTheme, day DayEntry, sport string) {
	theme.Label = firstNonEmpty(theme.Label, day.Title, sport)
	hay := strings.ToLower(day.Title + " " + day.Location + " " + sport)
	switch {
	case strings.Contains(hay, "bonaire"), strings.Contains(hay, "caribbean"), strings.Contains(hay, "ocean"):
		theme.Primary, theme.Secondary, theme.Accent = "#0aa3a0", "#04343f", "#7ed0e0"
	case strings.Contains(hay, "camp"), strings.Contains(hay, "forest"):
		theme.Primary, theme.Secondary, theme.Accent = "#5d8a4a", "#1e2f1a", "#c4a35a"
	case strings.Contains(hay, "school"), strings.Contains(hay, "first day"):
		theme.Primary, theme.Secondary, theme.Accent = "#3d6ea8", "#1a2740", "#f0c14b"
	default:
		theme.Primary, theme.Secondary, theme.Accent = "#0f8f7c", "#063642", "#5eb4c8"
	}
}

func persistDayTheme(dayID string, theme DayTheme) {
	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		return
	}
	cp := theme
	day.Theme = &cp
	dayStore[dayID] = day
	storeMu.Unlock()
	if err := persistManifest(); err != nil {
		log.Warnf("GenerateTheme: persist failed: %v", err)
	}
}

func getDayCopy(dayID string) DayEntry {
	storeMu.RLock()
	defer storeMu.RUnlock()
	day := dayStore[dayID]
	day.Media = visibleMedia(day.Media)
	day.Published = dayPublished(DayEntry{Media: day.Media})
	return day
}

func putThemeObject(dayID, role, ext, contentType string, data []byte) (string, string, error) {
	if !mediaObjectStore.Enabled() {
		return "", "", fmt.Errorf("object store disabled")
	}
	if ext == "" {
		ext = ".png"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	key := fmt.Sprintf("surfing/albums/%s/theme/%s-%d%s", dayID, role, time.Now().Unix(), ext)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if err := mediaObjectStore.Put(ctx, key, data, contentType); err != nil {
		return key, "", err
	}
	return key, mediaObjectStore.PublicURL(key), nil
}

func extFromContentType(ct string) string {
	ct = strings.ToLower(ct)
	switch {
	case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
		return ".jpg"
	case strings.Contains(ct, "webp"):
		return ".webp"
	default:
		return ".png"
	}
}

func buildBannerPrompt(day DayEntry, sport, userPrompt, brief string) string {
	return fmt.Sprintf(
		"Wide cinematic website hero banner for a personal media gallery. Sport/theme: %s. Album: %s%s. User story: %s. Style brief: %s. "+
			"Abstract-to-photographic blend inspired by wind, water, travel photography — NO text, NO logos, NO watermarks, NO people faces. "+
			"Rich atmosphere suitable as a full-bleed web banner.",
		sport, day.Title, locationSuffix(day.Location), userPrompt, brief,
	)
}

func buildBackgroundPrompt(day DayEntry, sport, userPrompt, brief string) string {
	return fmt.Sprintf(
		"Soft seamless website background wash for a photo gallery. Sport/theme: %s. Album: %s. Story: %s. Brief: %s. "+
			"Subtle texture, low-contrast, elegant, no text, no logos, suitable behind UI content.",
		sport, day.Title, userPrompt, brief,
	)
}

type styleBriefResult struct {
	Label     string `json:"label"`
	Brief     string `json:"brief"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Accent    string `json:"accent"`
}

func aiStyleBrief(ctx context.Context, cfg aiConfig, day DayEntry, sport, userPrompt string, samples []MediaItem) (styleBriefResult, error) {
	type contentPart map[string]any
	parts := []contentPart{
		{"type": "text", "text": fmt.Sprintf(
			"You design visual themes for a personal CDN gallery (like a sport trip page). "+
				"Album title=%q location=%q sport=%q user_prompt=%q. "+
				"Return STRICT JSON only with keys: label, brief, primary, secondary, accent. "+
				"Colors must be hex (#RRGGBB). brief is 1-2 sentences describing atmosphere for image generation. "+
				"Tailor to the sample photos if provided.",
			day.Title, day.Location, sport, userPrompt,
		)},
	}
	for _, s := range samples {
		url := s.URL
		if !strings.HasPrefix(url, "http") {
			continue // vision needs absolute URLs (CDN)
		}
		parts = append(parts, contentPart{
			"type": "image_url",
			"image_url": map[string]any{
				"url": url,
			},
		})
	}

	body := map[string]any{
		"model": cfg.ChatModel,
		"messages": []map[string]any{
			{"role": "user", "content": parts},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.7,
	}
	raw, err := aiPOST(ctx, cfg, "/chat/completions", body)
	if err != nil {
		return styleBriefResult{}, err
	}
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return styleBriefResult{}, err
	}
	if len(resp.Choices) == 0 {
		return styleBriefResult{}, fmt.Errorf("empty chat response")
	}
	var out styleBriefResult
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &out); err != nil {
		return styleBriefResult{}, fmt.Errorf("parse brief json: %w", err)
	}
	return out, nil
}

func aiGenerateImage(ctx context.Context, cfg aiConfig, prompt, size string) ([]byte, string, error) {
	body := map[string]any{
		"model":  cfg.ImageModel,
		"prompt": prompt,
		"n":      1,
		"size":   size,
	}
	model := strings.ToLower(cfg.ImageModel)
	if strings.Contains(model, "dall-e") {
		body["quality"] = "standard"
	}

	raw, err := aiPOST(ctx, cfg, "/images/generations", body)
	if err != nil {
		return nil, "", err
	}
	var resp struct {
		Data []struct {
			URL     string `json:"url"`
			B64JSON string `json:"b64_json"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, "", err
	}
	if len(resp.Data) == 0 {
		return nil, "", fmt.Errorf("empty image response")
	}
	item := resp.Data[0]
	if item.B64JSON != "" {
		decoded, err := base64.StdEncoding.DecodeString(item.B64JSON)
		if err != nil {
			return nil, "", err
		}
		return decoded, "image/png", nil
	}
	if item.URL == "" {
		return nil, "", fmt.Errorf("no image url/b64")
	}
	return downloadBytes(ctx, item.URL)
}

func imageSizeForModel(model string) string {
	m := strings.ToLower(model)
	if strings.Contains(m, "gpt-image") || strings.Contains(m, "chatgpt-image") {
		return "1536x1024"
	}
	if strings.Contains(m, "dall-e-3") {
		return "1792x1024"
	}
	return "1024x1024"
}

func downloadBytes(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	client := &http.Client{Timeout: 120 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return nil, "", fmt.Errorf("download %s: %s", res.Status, string(b))
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	ct := res.Header.Get("Content-Type")
	if ct == "" {
		ct = "image/png"
	}
	return data, ct, nil
}

func aiPOST(ctx context.Context, cfg aiConfig, path string, payload any) ([]byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BaseURL+path, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	client := &http.Client{Timeout: 300 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 8<<20))
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("ai %s: %s — %s", path, res.Status, strings.TrimSpace(string(raw)))
	}
	return raw, nil
}
