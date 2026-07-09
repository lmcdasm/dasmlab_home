package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const preloadIDNamespace = "preload:"

func preloadDayID(slug string) string {
	sum := sha256.Sum256([]byte(preloadIDNamespace + "day:" + strings.ToLower(slug)))
	return hex.EncodeToString(sum[:16])
}

func preloadMediaID(daySlug, filename string) string {
	sum := sha256.Sum256([]byte(preloadIDNamespace + "media:" + strings.ToLower(daySlug) + ":" + strings.ToLower(filename)))
	return hex.EncodeToString(sum[:16])
}

type preloadDayMeta struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	Location string `json:"location"`
}

func PreloadDaysFromDir(dir string) (int, error) {
	if dir == "" {
		return 0, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	imported := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		slug := entry.Name()
		dayDir := filepath.Join(dir, slug)
		count, err := preloadDayFromDir(slug, dayDir)
		if err != nil {
			log.Warnf("Preload: day %s failed: %v", slug, err)
			continue
		}
		imported += count
	}

	return imported, nil
}

func preloadDayFromDir(slug, dayDir string) (int, error) {
	meta := preloadDayMeta{
		Title: slug,
		Date:  time.Now().UTC().Format("2006-01-02"),
	}

	metaPath := filepath.Join(dayDir, "day.json")
	if raw, err := os.ReadFile(metaPath); err == nil {
		_ = json.Unmarshal(raw, &meta)
	}
	if strings.TrimSpace(meta.Title) == "" {
		meta.Title = slug
	}

	dayID := preloadDayID(slug)
	storeMu.Lock()
	day, exists := dayStore[dayID]
	if !exists {
		day = DayEntry{
			ID:        dayID,
			Title:     meta.Title,
			Date:      meta.Date,
			Location:  meta.Location,
			CreatedAt: time.Now().UTC(),
			Media:     []MediaItem{},
		}
	}
	existing := make(map[string]struct{}, len(day.Media))
	for _, item := range day.Media {
		existing[item.ID] = struct{}{}
	}
	storeMu.Unlock()

	files, err := os.ReadDir(dayDir)
	if err != nil {
		return 0, err
	}

	added := 0
	for _, file := range files {
		if file.IsDir() || file.Name() == "day.json" {
			continue
		}

		mediaID := preloadMediaID(slug, file.Name())
		if _, ok := existing[mediaID]; ok {
			continue
		}

		srcPath := filepath.Join(dayDir, file.Name())
		mediaType, ext, err := preloadMediaType(srcPath, file.Name())
		if err != nil || mediaType == "" {
			log.Warnf("Preload: skip %s: %v", file.Name(), err)
			continue
		}

		dstPath := mediaFilePath(mediaID, ext)
		if err := copyFile(srcPath, dstPath); err != nil {
			log.Warnf("Preload: copy %s failed: %v", file.Name(), err)
			continue
		}

		item := MediaItem{
			ID:        mediaID,
			Filename:  file.Name(),
			MediaType: mediaType,
			URL:       "/serve?id=" + mediaID,
			CreatedAt: time.Now().UTC(),
		}

		storeMu.Lock()
		day.Media = append(day.Media, item)
		dayStore[dayID] = day
		storeMu.Unlock()
		added++
	}

	if added > 0 || !exists {
		storeMu.Lock()
		dayStore[dayID] = day
		storeMu.Unlock()
		if err := persistManifest(); err != nil {
			log.Warnf("Preload: persist manifest failed: %v", err)
		}
	}

	log.Infof("Preload: day %s (%s) added %d media file(s)", slug, meta.Title, added)
	return added, nil
}

func preloadMediaType(path, filename string) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic", ".heif":
		return "image", ext, nil
	case ".mp4", ".mov", ".m4v", ".webm", ".avi", ".mkv":
		return "video", ext, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	mediaType, extFromDetect := detectMediaType(buf[:n], filename)
	if mediaType == "" {
		return "", "", os.ErrInvalid
	}
	if ext == "" {
		ext = extFromDetect
	}
	return mediaType, ext, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o775); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
