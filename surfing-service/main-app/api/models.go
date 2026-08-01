package api

import (
	"strings"
	"time"
)

// Media kinds for themed gallery sections (Videos / Photos / Other).
const (
	KindPhoto = "photo"
	KindVideo = "video"
	KindOther = "other"
)

// MediaItem is one photo/video/share in an album (DayEntry).
// Draft file items live on PVC (url=/serve?id=…). After Publish, url is the CDN/R2 public URL.
// Link-only "other" items (Garmin/iPhone shares) use ExternalURL and never touch object storage.
type MediaItem struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	MediaType string    `json:"media_type"` // image | video | other
	Kind      string    `json:"kind,omitempty"` // photo | video | other
	Caption   string    `json:"caption,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	URL       string    `json:"url"`
	// ExternalURL is an outbound share (activity link, video open-out, etc.).
	ExternalURL string    `json:"external_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	// Published means bytes are on R2/CDN and URL is absolute (or link-only origin).
	Published bool   `json:"published,omitempty"`
	Origin    string `json:"origin,omitempty"` // pvc | r2 | link
	ObjectKey string `json:"object_key,omitempty"`
	// Hidden soft-removes from gallery listings without dropping bytes (hard drop later).
	// Keeping the row also stops preload from re-seeding the same file on restart.
	Hidden bool `json:"hidden,omitempty"`
}

func normalizeMediaKind(item *MediaItem) {
	if item == nil {
		return
	}
	k := strings.ToLower(strings.TrimSpace(item.Kind))
	switch k {
	case KindPhoto, KindVideo, KindOther:
		item.Kind = k
		return
	}
	switch strings.ToLower(strings.TrimSpace(item.MediaType)) {
	case "video":
		item.Kind = KindVideo
	case "image", "photo":
		item.Kind = KindPhoto
	default:
		if item.ExternalURL != "" {
			item.Kind = KindOther
		} else {
			item.Kind = KindPhoto
		}
	}
}

// DayEntry is an album / trip day (Surfing "day").
type DayEntry struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Date      string      `json:"date"`
	Location  string      `json:"location,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	Media     []MediaItem `json:"media"`
	// Published is true when every media item is published (convenience for UI).
	Published bool `json:"published,omitempty"`
	// Theme is an AI-tailored banner/background derived from album samples.
	Theme *DayTheme `json:"theme,omitempty"`
}

// DayTheme is a sport/event visual skin for one album (banner + page wash).
type DayTheme struct {
	Label          string    `json:"label,omitempty"`
	Prompt         string    `json:"prompt,omitempty"`
	StyleBrief     string    `json:"style_brief,omitempty"`
	BannerURL      string    `json:"banner_url,omitempty"`
	BackgroundURL  string    `json:"background_url,omitempty"`
	Primary        string    `json:"primary,omitempty"`
	Secondary      string    `json:"secondary,omitempty"`
	Accent         string    `json:"accent,omitempty"`
	Provider       string    `json:"provider,omitempty"` // openai | cheapcloud | local
	SampleMediaIDs []string  `json:"sample_media_ids,omitempty"`
	GeneratedAt    time.Time `json:"generated_at,omitempty"`
}

func dayPublished(day DayEntry) bool {
	visible := 0
	for _, m := range day.Media {
		if m.Hidden {
			continue
		}
		visible++
		if !m.Published {
			return false
		}
	}
	return visible > 0
}

func visibleMedia(items []MediaItem) []MediaItem {
	out := make([]MediaItem, 0, len(items))
	for _, m := range items {
		if m.Hidden {
			continue
		}
		normalizeMediaKind(&m)
		out = append(out, m)
	}
	return out
}
