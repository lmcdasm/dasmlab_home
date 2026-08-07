package api

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const sharesFileName = "shares.json"

var (
	shareStore = make(map[string]*ShareLink)
	shareMu    sync.RWMutex
)

// ShareLink is a DASMLAB-fronted pointer into album/media (tollgate).
// Socials get this URL — not a raw forever-upload. Hits are counted for cheapcloud.
type ShareLink struct {
	Token       string    `json:"token"`
	DayID       string    `json:"day_id"`
	MediaID     string    `json:"media_id,omitempty"`
	Scope       string    `json:"scope"` // album | media
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	CDNURL      string    `json:"cdn_url"`
	TargetURL   string    `json:"target_url"` // redirect destination (CDN or album page)
	ImageURL    string    `json:"image_url,omitempty"`
	Hits        int64     `json:"hits"`
	CreatedAt   time.Time `json:"created_at"`
	// UnlockCodeHash enables premium "only these people" later (email/SMS delivery).
	UnlockCodeHash string `json:"unlock_code_hash,omitempty"`
	Private        bool   `json:"private,omitempty"`
}

type createShareRequest struct {
	DayID      string `json:"day_id"`
	MediaID    string `json:"media_id"`
	UnlockCode string `json:"unlock_code"` // optional premium foreshadow
	Private    bool   `json:"private"`
	AlbumPage  string `json:"album_page"` // FE SPA URL for album-scope shares
}

type createShareResponse struct {
	Token       string `json:"token"`
	TollgateURL string `json:"tollgate_url"`
	CDNURL      string `json:"cdn_url"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	ImageURL    string `json:"image_url,omitempty"`
	Hits        int64  `json:"hits"`
	Private     bool   `json:"private"`
	Scope       string `json:"scope"`
	Channels    []shareChannel `json:"channels"`
}

type shareChannel struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Kind        string `json:"kind"` // intent | copy | native
	URL         string `json:"url,omitempty"`
	Hint        string `json:"hint,omitempty"`
	Hot         bool   `json:"hot,omitempty"`
}

func sharesPath() string {
	return filepath.Join(dataDir, sharesFileName)
}

func loadShares() error {
	raw, err := os.ReadFile(sharesPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var list []*ShareLink
	if err := json.Unmarshal(raw, &list); err != nil {
		return err
	}
	shareMu.Lock()
	defer shareMu.Unlock()
	shareStore = make(map[string]*ShareLink, len(list))
	for _, s := range list {
		if s == nil || s.Token == "" {
			continue
		}
		shareStore[s.Token] = s
	}
	return nil
}

func persistShares() error {
	shareMu.RLock()
	list := make([]*ShareLink, 0, len(shareStore))
	for _, s := range shareStore {
		cp := *s
		list = append(list, &cp)
	}
	shareMu.RUnlock()

	payload, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := sharesPath() + ".tmp"
	if err := os.WriteFile(tmp, payload, 0o664); err != nil {
		return err
	}
	return os.Rename(tmp, sharesPath())
}

func newShareToken() (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func hashUnlock(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	sum := sha256.Sum256([]byte("surfing-share:" + code))
	return hex.EncodeToString(sum[:])
}

func sharePublicBase(c *gin.Context) string {
	if v := strings.TrimSpace(os.Getenv("SURFING_SHARE_PUBLIC_BASE")); v != "" {
		return strings.TrimRight(v, "/")
	}
	proto := c.GetHeader("X-Forwarded-Proto")
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	// FE nginx proxies /api/surfing/ → service /. Direct service hosts (surfing.svc.*)
	// must mint /s/:token without that prefix.
	basePath := strings.TrimSpace(os.Getenv("SURFING_SHARE_PATH_PREFIX"))
	if basePath == "" {
		basePath = "/api/surfing"
	}
	h := strings.ToLower(host)
	if strings.Contains(h, "surfing.svc.") || strings.HasPrefix(h, "surfing.") {
		basePath = ""
	}
	return fmt.Sprintf("%s://%s%s", proto, host, strings.TrimRight(basePath, "/"))
}

// CreateShare mints a tollgate link for an album or a single media item.
func CreateShare(c *gin.Context) {
	var req createShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}
	dayID := strings.TrimSpace(req.DayID)
	if dayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "day_id is required"})
		return
	}

	storeMu.RLock()
	day, ok := dayStore[dayID]
	storeMu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	title := day.Title
	text := buildAlbumShareText(day)
	cdnURL := ""
	imageURL := ""
	targetURL := strings.TrimSpace(req.AlbumPage)
	scope := "album"
	mediaID := strings.TrimSpace(req.MediaID)

	if mediaID != "" {
		var item *MediaItem
		for i := range day.Media {
			if day.Media[i].ID == mediaID && !day.Media[i].Hidden {
				item = &day.Media[i]
				break
			}
		}
		if item == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
			return
		}
		scope = "media"
		title = firstNonEmpty(item.Caption, item.Filename, day.Title)
		text = buildMediaShareText(day, *item)
		cdnURL = absoluteMediaURL(item.URL)
		if item.ExternalURL != "" {
			cdnURL = item.ExternalURL
		}
		targetURL = cdnURL
		if item.Kind == KindPhoto || item.MediaType == "image" {
			imageURL = cdnURL
		} else if day.Theme != nil && day.Theme.BannerURL != "" {
			imageURL = day.Theme.BannerURL
		}
	} else {
		if day.Theme != nil && day.Theme.BannerURL != "" {
			imageURL = day.Theme.BannerURL
		}
		for _, m := range day.Media {
			if m.Hidden {
				continue
			}
			normalizeMediaKind(&m)
			if m.Kind == KindPhoto && strings.HasPrefix(m.URL, "http") {
				if imageURL == "" {
					imageURL = m.URL
				}
				if cdnURL == "" {
					cdnURL = m.URL
				}
				break
			}
		}
		if targetURL == "" {
			// Fallback: first CDN asset so the link always resolves somewhere useful.
			targetURL = cdnURL
			if targetURL == "" {
				targetURL = sharePublicBase(c) + "/days"
			}
		}
	}

	token, err := newShareToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not mint share token"})
		return
	}

	link := &ShareLink{
		Token:          token,
		DayID:          dayID,
		MediaID:        mediaID,
		Scope:          scope,
		Title:          title,
		Text:           text,
		CDNURL:         cdnURL,
		TargetURL:      targetURL,
		ImageURL:       imageURL,
		Hits:           0,
		CreatedAt:      time.Now().UTC(),
		UnlockCodeHash: hashUnlock(req.UnlockCode),
		Private:        req.Private || strings.TrimSpace(req.UnlockCode) != "",
	}

	shareMu.Lock()
	shareStore[token] = link
	shareMu.Unlock()
	if err := persistShares(); err != nil {
		log.Warnf("CreateShare: persist failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist share"})
		return
	}

	tollgate := sharePublicBase(c) + "/s/" + token
	resp := createShareResponse{
		Token:       token,
		TollgateURL: tollgate,
		CDNURL:      cdnURL,
		Title:       title,
		Text:        text,
		ImageURL:    imageURL,
		Hits:        0,
		Private:     link.Private,
		Scope:       scope,
		Channels:    buildShareChannels(tollgate, title, text),
	}
	c.JSON(http.StatusCreated, resp)
}

func absoluteMediaURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return ""
	}
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("SURFING_SHARE_PUBLIC_BASE")), "/")
	if base == "" {
		return u
	}
	if strings.HasPrefix(u, "/") {
		return base + u
	}
	return base + "/" + u
}

func buildAlbumShareText(day DayEntry) string {
	parts := []string{day.Title}
	if day.Location != "" {
		parts = append(parts, "— "+day.Location)
	}
	if day.Date != "" {
		parts = append(parts, "("+day.Date+")")
	}
	parts = append(parts, "Private gallery link (you control the source — not a Facebook upload).")
	return strings.Join(parts, " ")
}

func buildMediaShareText(day DayEntry, item MediaItem) string {
	cap := firstNonEmpty(item.Caption, item.Filename)
	notes := strings.TrimSpace(item.Notes)
	if notes != "" {
		if len(notes) > 180 {
			notes = notes[:177] + "…"
		}
		return fmt.Sprintf("%s — %s · from %s", cap, notes, day.Title)
	}
	return fmt.Sprintf("%s · from %s%s", cap, day.Title, locationSuffix(day.Location))
}

func buildShareChannels(tollgate, title, text string) []shareChannel {
	encURL := url.QueryEscape(tollgate)
	encText := url.QueryEscape(text)
	encTitle := url.QueryEscape(title)
	combo := url.QueryEscape(text + "\n" + tollgate)
	return []shareChannel{
		{ID: "copy_tollgate", Label: "Copy share link", Kind: "copy", Hint: "DASMLAB tollgate — hits counted for FinOps"},
		{ID: "copy_cdn", Label: "Copy CDN link", Kind: "copy", Hint: "Direct edge URL (bypass tollgate)"},
		{ID: "whatsapp", Label: "WhatsApp", Kind: "intent", Hot: true, URL: "https://api.whatsapp.com/send?text=" + combo},
		{ID: "facebook", Label: "Facebook", Kind: "intent", Hot: true, URL: "https://www.facebook.com/sharer/sharer.php?u=" + encURL},
		{ID: "threads", Label: "Threads", Kind: "intent", Hot: true, URL: "https://www.threads.net/intent/post?url=" + encURL + "&text=" + encText},
		{ID: "x", Label: "X", Kind: "intent", Hot: true, URL: "https://twitter.com/intent/tweet?url=" + encURL + "&text=" + encTitle},
		{ID: "linkedin", Label: "LinkedIn", Kind: "intent", URL: "https://www.linkedin.com/sharing/share-offsite/?url=" + encURL},
		{ID: "instagram", Label: "Instagram", Kind: "copy", Hint: "No web intent — copy link, paste into Instagram"},
		{ID: "tiktok", Label: "TikTok", Kind: "copy", Hint: "No web intent — copy link, paste into TikTok"},
		{ID: "native", Label: "Device share", Kind: "native", Hint: "Uses the OS share sheet when available"},
	}
}

// ResolveShare is the tollgate: allow-check → meter hit → OG for bots or redirect for humans.
func ResolveShare(c *gin.Context) {
	token := c.Param("token")
	unlock := strings.TrimSpace(c.Query("unlock"))

	shareMu.Lock()
	link, ok := shareStore[token]
	if !ok || link == nil {
		shareMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "share not found"})
		return
	}
	if link.UnlockCodeHash != "" {
		if hashUnlock(unlock) != link.UnlockCodeHash {
			shareMu.Unlock()
			c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(unlockGateHTML(link.Title)))
			return
		}
	}
	link.Hits++
	hits := link.Hits
	snapshot := *link
	shareMu.Unlock()

	_ = persistShares()
	go emitShareHit(snapshot)

	if isSocialCrawler(c.GetHeader("User-Agent")) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(ogHTML(snapshot)))
		return
	}

	log.Infof("share_hit token=%s day=%s media=%s hits=%d target=%s", snapshot.Token, snapshot.DayID, snapshot.MediaID, hits, snapshot.TargetURL)
	c.Redirect(http.StatusFound, snapshot.TargetURL)
}

// ShareMeta returns JSON metadata for a share (debugging / FE).
func ShareMeta(c *gin.Context) {
	token := c.Param("token")
	shareMu.RLock()
	link, ok := shareStore[token]
	shareMu.RUnlock()
	if !ok || link == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "share not found"})
		return
	}
	c.JSON(http.StatusOK, link)
}

func isSocialCrawler(ua string) bool {
	ua = strings.ToLower(ua)
	needles := []string{
		"facebookexternalhit", "facebot", "twitterbot", "linkedinbot",
		"slackbot", "discordbot", "whatsapp", "telegrambot", "threads",
	}
	for _, n := range needles {
		if strings.Contains(ua, n) {
			return true
		}
	}
	return false
}

func ogHTML(link ShareLink) string {
	title := html.EscapeString(link.Title)
	desc := html.EscapeString(link.Text)
	img := html.EscapeString(link.ImageURL)
	canon := html.EscapeString(link.TargetURL)
	var b bytes.Buffer
	b.WriteString("<!doctype html><html><head><meta charset=\"utf-8\">")
	b.WriteString("<title>" + title + "</title>")
	b.WriteString("<meta property=\"og:title\" content=\"" + title + "\">")
	b.WriteString("<meta property=\"og:description\" content=\"" + desc + "\">")
	b.WriteString("<meta property=\"og:type\" content=\"website\">")
	if canon != "" {
		b.WriteString("<meta property=\"og:url\" content=\"" + canon + "\">")
	}
	if img != "" {
		b.WriteString("<meta property=\"og:image\" content=\"" + img + "\">")
	}
	b.WriteString("<meta name=\"twitter:card\" content=\"summary_large_image\">")
	b.WriteString("<meta name=\"twitter:title\" content=\"" + title + "\">")
	b.WriteString("<meta name=\"twitter:description\" content=\"" + desc + "\">")
	if img != "" {
		b.WriteString("<meta name=\"twitter:image\" content=\"" + img + "\">")
	}
	b.WriteString("<meta http-equiv=\"refresh\" content=\"0;url=" + html.EscapeString(link.TargetURL) + "\">")
	b.WriteString("</head><body><p><a href=\"" + html.EscapeString(link.TargetURL) + "\">Continue</a></p></body></html>")
	return b.String()
}

func unlockGateHTML(title string) string {
	t := html.EscapeString(title)
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>Unlock required</title></head>" +
		"<body style=\"font-family:system-ui;padding:2rem;max-width:28rem\">" +
		"<h1>Private share</h1><p>" + t + " needs an unlock code (premium foreshadow).</p>" +
		"<p>Append <code>?unlock=YOUR_CODE</code> to the link.</p></body></html>"
}

func emitShareHit(link ShareLink) {
	endpoint := strings.TrimSpace(os.Getenv("SURFING_CHEAPCLOUD_HIT_URL"))
	if endpoint == "" {
		return
	}
	payload, _ := json.Marshal(map[string]any{
		"product_id":  "dasmlab-surfing",
		"event":       "share_hit",
		"token":       link.Token,
		"day_id":      link.DayID,
		"media_id":    link.MediaID,
		"scope":       link.Scope,
		"hits":        link.Hits,
		"cdn_url":     link.CDNURL,
		"recorded_at": time.Now().UTC().Format(time.RFC3339),
	})
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 8 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Warnf("share_hit cheapcloud emit failed: %v", err)
		return
	}
	res.Body.Close()
}
