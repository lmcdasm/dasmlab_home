package api

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"surfing-service/logutil"
)

var (
	componentName = "surfing-api"
	log           = logutil.InitLogger(componentName)
	maxUploadBytes int64 = 500 << 20 // 500 MiB
)

func init() {
	if raw := strings.TrimSpace(os.Getenv("SURFING_MAX_UPLOAD_MB")); raw != "" {
		if mb, err := strconv.ParseInt(raw, 10, 64); err == nil && mb > 0 {
			maxUploadBytes = mb << 20
		}
	}
}

func IsAlive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func ListDays(c *gin.Context) {
	if err := reloadManifestFromDisk(); err != nil {
		log.Warnf("ListDays: reload manifest failed: %v", err)
	}

	owner := isOwner(c)
	group := false
	if u, ok := currentUser(c); ok {
		group = hasGroupRole(u)
	}

	storeMu.RLock()
	defer storeMu.RUnlock()

	days := make([]DayEntry, 0, len(dayStore))
	for _, day := range dayStore {
		if day.TagPolicy == "" {
			day.TagPolicy = "public"
		}
		day.Media = projectMediaFor(day.Media, owner, group)
		for i := range day.Media {
			if day.Media[i].CanDownload {
				day.Media[i].DownloadPath = "/days/" + day.ID + "/media/" + day.Media[i].ID + "/download"
			}
		}
		if day.DatePrecision == "" {
			day.DatePrecision = "day"
		}
		day.Published = dayPublished(DayEntry{Media: day.Media})
		days = append(days, day)
	}
	c.JSON(http.StatusOK, days)
}

func CreateDay(c *gin.Context) {
	var req struct {
		Title    string `json:"title"`
		Date     string `json:"date"`
		Location string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	date := strings.TrimSpace(req.Date)
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}

	id := uuid.NewString()
	day := DayEntry{
		ID:        id,
		Title:     title,
		Date:      date,
		Location:  strings.TrimSpace(req.Location),
		CreatedAt: time.Now().UTC(),
		Media:     []MediaItem{},
	}

	storeMu.Lock()
	dayStore[id] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("CreateDay: failed to persist manifest: %v", err)
	}

	c.JSON(http.StatusCreated, day)
}

func DeleteDay(c *gin.Context) {
	id := c.Param("id")

	storeMu.Lock()
	day, ok := dayStore[id]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	for _, item := range day.Media {
		removeMediaFile(item.ID, filepath.Ext(item.Filename))
		deleteMediaObject(id, item.ID, filepath.Ext(item.Filename))
	}
	delete(dayStore, id)
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("DeleteDay: failed to persist manifest: %v", err)
	}

	c.Status(http.StatusNoContent)
}

func UploadMedia(c *gin.Context) {
	dayID := c.Param("id")

	storeMu.RLock()
	_, ok := dayStore[dayID]
	storeMu.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read file"})
		return
	}

	mediaType, ext := detectMediaType(data, header.Filename)
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type (photos and videos only)"})
		return
	}

	mediaID := uuid.NewString()
	targetPath := mediaFilePath(mediaID, ext)
	if err := os.WriteFile(targetPath, data, 0o664); err != nil {
		log.Errorf("UploadMedia: failed to write %s: %v", targetPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store file"})
		return
	}

	// Draft-first: PVC staging. Publish pushes to R2/CDN.
	// SURFING_AUTO_PUBLISH=1 dual-writes on upload (optional).
	serveURL := "/serve?id=" + mediaID
	item := MediaItem{
		ID:        mediaID,
		Filename:  header.Filename,
		MediaType: mediaType,
		Kind:      strings.TrimSpace(c.PostForm("kind")),
		Caption:   strings.TrimSpace(c.PostForm("caption")),
		Notes:     strings.TrimSpace(c.PostForm("notes")),
		URL:       serveURL,
		CreatedAt: time.Now().UTC(),
		Published: false,
		Origin:    "pvc",
	}
	if item.Kind == "" {
		if mediaType == "video" {
			item.Kind = KindVideo
		} else {
			item.Kind = KindPhoto
		}
	}
	normalizeMediaKind(&item)

	autoPublish := strings.EqualFold(strings.TrimSpace(os.Getenv("SURFING_AUTO_PUBLISH")), "1") ||
		strings.EqualFold(strings.TrimSpace(os.Getenv("SURFING_AUTO_PUBLISH")), "true")
	if autoPublish && mediaObjectStore.Enabled() {
		key, publicURL, err := putMediaObject(dayID, mediaID, ext, mimeFromExt(ext, mediaType), data)
		if err != nil {
			log.Errorf("UploadMedia: auto-publish put failed: %v", err)
			_ = os.Remove(targetPath)
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to store file on CDN origin"})
			return
		}
		if publicURL != "" {
			item.URL = publicURL
			item.Published = true
			item.Origin = "r2"
			item.ObjectKey = key
		}
	}

	storeMu.Lock()
	day := dayStore[dayID]
	day.Media = append(day.Media, item)
	day.Published = dayPublished(day)
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("UploadMedia: failed to persist manifest: %v", err)
	}

	c.JSON(http.StatusCreated, item)
}

func DeleteMedia(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")

	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	found := -1
	for i, item := range day.Media {
		if item.ID == mediaID {
			found = i
			break
		}
	}
	if found < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	// Soft-hide only: keep bytes on PVC/R2; hard drop comes later via cdn-mgr.
	// Leaving the row (Hidden=true) also blocks preload from re-importing the file.
	day.Media[found].Hidden = true
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("DeleteMedia: failed to persist manifest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist hide"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateMedia patches caption, notes, kind, and/or external_url (soft metadata only).
func UpdateMedia(c *gin.Context) {
	dayID := c.Param("id")
	mediaID := c.Param("mediaId")

	var req struct {
		Caption            *string `json:"caption"`
		Notes              *string `json:"notes"`
		Kind               *string `json:"kind"`
		ExternalURL        *string `json:"external_url"`
		NotesVisibility    *string `json:"notes_visibility"`
		DownloadVisibility *string `json:"download_visibility"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}

	found := -1
	for i, item := range day.Media {
		if item.ID == mediaID {
			found = i
			break
		}
	}
	if found < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	item := day.Media[found]
	if req.Caption != nil {
		item.Caption = strings.TrimSpace(*req.Caption)
	}
	if req.Notes != nil {
		item.Notes = strings.TrimSpace(*req.Notes)
	}
	if req.Kind != nil {
		item.Kind = strings.TrimSpace(*req.Kind)
	}
	if req.ExternalURL != nil {
		item.ExternalURL = strings.TrimSpace(*req.ExternalURL)
	}
	if req.NotesVisibility != nil {
		v := strings.ToLower(strings.TrimSpace(*req.NotesVisibility))
		switch v {
		case "public", "private", "group":
			item.NotesVisibility = v
		default:
			storeMu.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "notes_visibility must be public|private|group"})
			return
		}
	}
	if req.DownloadVisibility != nil {
		v := strings.ToLower(strings.TrimSpace(*req.DownloadVisibility))
		switch v {
		case "public", "private", "group":
			item.DownloadVisibility = v
		default:
			storeMu.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "download_visibility must be public|private|group"})
			return
		}
	}
	normalizeMediaKind(&item)
	day.Media[found] = item
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("UpdateMedia: failed to persist manifest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist media update"})
		return
	}

	owner := isOwner(c)
	group := false
	if u, ok := currentUser(c); ok {
		group = hasGroupRole(u)
	}
	projected := projectMediaFor([]MediaItem{item}, owner, group)
	resp := projected[0]
	if resp.CanDownload {
		resp.DownloadPath = "/days/" + dayID + "/media/" + resp.ID + "/download"
	}
	c.JSON(http.StatusOK, resp)
}

// AddMediaLink creates a link-only "other" item (Garmin / iPhone / activity shares).
func AddMediaLink(c *gin.Context) {
	dayID := c.Param("id")

	var req struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Notes       string `json:"notes"`
		Caption     string `json:"caption"`
		Kind        string `json:"kind"`
		SourceLabel string `json:"source_label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	link := strings.TrimSpace(req.URL)
	if link == "" || !(strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url must be http(s)"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = strings.TrimSpace(req.Caption)
	}
	if title == "" {
		title = "Shared activity"
	}
	caption := strings.TrimSpace(req.Caption)
	if caption == "" {
		caption = title
	}
	if label := strings.TrimSpace(req.SourceLabel); label != "" && caption == title {
		caption = label + " · " + title
	}

	item := MediaItem{
		ID:          uuid.NewString(),
		Filename:    title,
		MediaType:   "other",
		Kind:        KindOther,
		Caption:     caption,
		Notes:       strings.TrimSpace(req.Notes),
		URL:         link,
		ExternalURL: link,
		CreatedAt:   time.Now().UTC(),
		Published:   true,
		Origin:      "link",
	}
	if k := strings.TrimSpace(req.Kind); k != "" {
		item.Kind = k
	}
	normalizeMediaKind(&item)

	storeMu.Lock()
	day, ok := dayStore[dayID]
	if !ok {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "day not found"})
		return
	}
	day.Media = append(day.Media, item)
	day.Published = dayPublished(day)
	dayStore[dayID] = day
	storeMu.Unlock()

	if err := persistManifest(); err != nil {
		log.Warnf("AddMediaLink: failed to persist manifest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist link"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func ServeMedia(c *gin.Context) {
	mediaID := c.Query("id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var filename string
	var itemURL string
	storeMu.RLock()
	for _, day := range dayStore {
		for _, item := range day.Media {
			if item.ID == mediaID {
				filename = item.Filename
				itemURL = item.URL
				break
			}
		}
		if filename != "" {
			break
		}
	}
	storeMu.RUnlock()

	if filename == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	// CDN / R2 public URLs: redirect so browsers hit the edge, not basement PVC.
	if strings.HasPrefix(itemURL, "http://") || strings.HasPrefix(itemURL, "https://") {
		c.Redirect(http.StatusFound, itemURL)
		return
	}

	ext := filepath.Ext(filename)
	path := mediaFilePath(mediaID, ext)
	if _, err := os.Stat(path); err != nil {
		// Uploaded/imported files may use a different extension case (.jpg vs .JPG).
		alt := mediaFilePath(mediaID, strings.ToLower(ext))
		if _, err2 := os.Stat(alt); err2 != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		path = alt
	}

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Type", contentType)
	if c.Query("download") == "1" || c.Query("download") == "true" {
		c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	}
	c.File(path)
}

func detectMediaType(data []byte, filename string) (string, string) {
	ext := strings.ToLower(filepath.Ext(filename))
	mtype := mimetype.Detect(data)

	switch {
	case strings.HasPrefix(mtype.String(), "image/"):
		if ext == "" {
			ext = extFromMime(mtype.Extension())
		}
		return "image", normalizeExt(ext)
	case strings.HasPrefix(mtype.String(), "video/"):
		if ext == "" {
			ext = extFromMime(mtype.Extension())
		}
		return "video", normalizeExt(ext)
	default:
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic", ".heif":
			return "image", ext
		case ".mp4", ".mov", ".m4v", ".webm", ".avi", ".mkv":
			return "video", ext
		default:
			return "", ""
		}
	}
}

func extFromMime(ext string) string {
	if ext == "" {
		return ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		return "." + ext
	}
	return ext
}

func mimeFromExt(ext, mediaType string) string {
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	switch mediaType {
	case "image":
		return "image/jpeg"
	case "video":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

func normalizeExt(ext string) string {
	if ext == "" {
		return ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		return "." + ext
	}
	return ext
}
