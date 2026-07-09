package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	defaultDataDir   = "/data"
	manifestFileName = "manifest.json"
	mediaSubdirName  = "media"
)

var (
	dataDir      string
	manifestPath string
	mediaBasePath string
	dayStore     = make(map[string]DayEntry)
	storeMu      sync.RWMutex
)

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func ensureStorageDirs() error {
	if err := os.MkdirAll(mediaBasePath, 0o775); err != nil {
		return err
	}
	return os.MkdirAll(dataDir, 0o775)
}

func loadManifest() error {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var days []DayEntry
	if err := json.Unmarshal(raw, &days); err != nil {
		return err
	}

	storeMu.Lock()
	defer storeMu.Unlock()
	for _, day := range days {
		if day.ID == "" {
			continue
		}
		if day.Media == nil {
			day.Media = []MediaItem{}
		}
		dayStore[day.ID] = day
	}
	return nil
}

func persistManifest() error {
	storeMu.RLock()
	days := make([]DayEntry, 0, len(dayStore))
	for _, day := range dayStore {
		days = append(days, day)
	}
	storeMu.RUnlock()

	payload, err := json.MarshalIndent(days, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := manifestPath + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o664); err != nil {
		return err
	}
	return os.Rename(tmpPath, manifestPath)
}

func mediaFilePath(mediaID, ext string) string {
	return filepath.Join(mediaBasePath, mediaID+ext)
}

func removeMediaFile(mediaID, ext string) {
	_ = os.Remove(mediaFilePath(mediaID, ext))
}

func Initialize() error {
	dataDir = envOrDefault("SURFING_DATA_DIR", defaultDataDir)
	mediaBasePath = filepath.Join(dataDir, mediaSubdirName)
	manifestPath = filepath.Join(dataDir, manifestFileName)

	if err := ensureStorageDirs(); err != nil {
		return err
	}

	if err := loadManifest(); err != nil {
		log.Warnf("Storage: could not load manifest %s: %v", manifestPath, err)
	}

	log.Infof(
		"Storage: ready data_dir=%s media=%s manifest=%s days=%d",
		dataDir,
		mediaBasePath,
		manifestPath,
		dayCount(),
	)
	return nil
}

func dayCount() int {
	storeMu.RLock()
	defer storeMu.RUnlock()
	return len(dayStore)
}
