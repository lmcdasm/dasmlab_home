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
	storeMu.RLock()
	defer storeMu.RUnlock()

	days := make([]DayEntry, 0, len(dayStore))
	for _, day := range dayStore {
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

	caption := strings.TrimSpace(c.PostForm("caption"))
	item := MediaItem{
		ID:        mediaID,
		Filename:  header.Filename,
		MediaType: mediaType,
		Caption:   caption,
		URL:       "/serve?id=" + mediaID,
		CreatedAt: time.Now().UTC(),
	}

	storeMu.Lock()
	day := dayStore[dayID]
	day.Media = append(day.Media, item)
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
	var filename string
	for i, item := range day.Media {
		if item.ID == mediaID {
			found = i
			filename = item.Filename
			break
		}
	}
	if found < 0 {
		storeMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	day.Media = append(day.Media[:found], day.Media[found+1:]...)
	dayStore[dayID] = day
	storeMu.Unlock()

	removeMediaFile(mediaID, filepath.Ext(filename))

	if err := persistManifest(); err != nil {
		log.Warnf("DeleteMedia: failed to persist manifest: %v", err)
	}

	c.Status(http.StatusNoContent)
}

func ServeMedia(c *gin.Context) {
	mediaID := c.Query("id")
	if mediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var filename string
	storeMu.RLock()
	for _, day := range dayStore {
		for _, item := range day.Media {
			if item.ID == mediaID {
				filename = item.Filename
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

func normalizeExt(ext string) string {
	if ext == "" {
		return ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		return "." + ext
	}
	return ext
}
