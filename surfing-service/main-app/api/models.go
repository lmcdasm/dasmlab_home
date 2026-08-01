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
	// Hidden soft-removes from public gallery. Owners still see grayed items when signed in.
	Hidden bool `json:"hidden,omitempty"`
	// PlayCount is simple click/play metering (tollgate + in-app viewer).
	PlayCount int64 `json:"play_count,omitempty"`
	// Tags are people-name associations only (no links). Owner must approve.
	Tags []MediaTag `json:"tags,omitempty"`
	// NotesVisibility: public | private | group (group = signed-up members later).
	NotesVisibility string `json:"notes_visibility,omitempty"`
	// DownloadVisibility: who may use the Download control / download API.
	// public = anyone; private = owner; group = Keycloak group role (+ owner).
	DownloadVisibility string `json:"download_visibility,omitempty"`
	// CanDownload is projected per-request (not persisted). Always serialized (false must not omit).
	CanDownload bool `json:"can_download"`
	// DownloadPath is the DASMLAB-gated path (obfuscated day/media ids — not a bare CDN guess).
	DownloadPath string `json:"download_path,omitempty"`
}

// MediaTag is a plain name on content ("I'm in this"). Not a social graph.
type MediaTag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"` // pending | approved | rejected
	CreatedAt time.Time `json:"created_at"`
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
	// TagPolicy: public | group | off — who may propose plain-name tags.
	TagPolicy string `json:"tag_policy,omitempty"`
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
	return projectMedia(items, false)
}

// projectMedia applies public vs owner visibility for gallery listings.
func projectMedia(items []MediaItem, owner bool) []MediaItem {
	return projectMediaFor(items, owner, false)
}

func projectMediaFor(items []MediaItem, owner bool, groupMember bool) []MediaItem {
	out := make([]MediaItem, 0, len(items))
	for _, m := range items {
		normalizeMediaKind(&m)
		if m.Hidden && !owner {
			continue
		}
		if m.NotesVisibility == "" {
			m.NotesVisibility = "public"
		}
		if m.DownloadVisibility == "" {
			m.DownloadVisibility = "public"
		}
		if !owner {
			switch strings.ToLower(m.NotesVisibility) {
			case "private", "group":
				m.Notes = ""
			}
			approved := make([]MediaTag, 0, len(m.Tags))
			for _, t := range m.Tags {
				if t.Status == "approved" {
					approved = append(approved, t)
				}
			}
			m.Tags = approved
		}
		m.CanDownload = downloadAllowed(m.DownloadVisibility, owner, groupMember)
		m.DownloadPath = ""
		m.ObjectKey = "" // never leak storage keys to clients
		out = append(out, m)
	}
	return out
}

func downloadAllowed(vis string, owner, groupMember bool) bool {
	switch strings.ToLower(strings.TrimSpace(vis)) {
	case "", "public":
		return true
	case "private":
		return owner
	case "group":
		return owner || groupMember
	default:
		return owner
	}
}

func hasGroupRole(u AuthUser) bool {
	if u.IsAdmin {
		return true
	}
	for _, r := range u.Roles {
		switch strings.ToLower(r) {
		case "group", "member", "album-group", "cdn-group":
			return true
		}
	}
	return false
}
