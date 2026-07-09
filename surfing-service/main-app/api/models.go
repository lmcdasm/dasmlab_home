package api

import "time"

type MediaItem struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	MediaType string    `json:"media_type"`
	Caption   string    `json:"caption,omitempty"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type DayEntry struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Date      string      `json:"date"`
	Location  string      `json:"location,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	Media     []MediaItem `json:"media"`
}
