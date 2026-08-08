// Package activity stores append-only engagement events (anon + known).
package activity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	TypeLogin    = "login"
	TypeNavigate = "navigate"
	TypeEngaged  = "engaged"
	TypePage     = "page"
	TypeIdentify = "identify"
	TypeAlias    = "alias"
	TypeTrack    = "track"

	DefaultListLimit = 200
	MaxListLimit     = 500
)

// Event is one activity log line (JSONL). First-party mini-CDP shape.
type Event struct {
	TS          time.Time `json:"ts"`
	Type        string    `json:"type"` // login|navigate|engaged|page|identify|alias|track
	Event       string    `json:"event,omitempty"`
	User        string    `json:"user,omitempty"` // preferred_username when known
	Sub         string    `json:"sub,omitempty"`
	Email       string    `json:"email,omitempty"`
	AnonymousID string    `json:"anonymousId,omitempty"`
	SessionID   string    `json:"sessionId,omitempty"`
	PreviousID  string    `json:"previousId,omitempty"` // alias: anon id being linked
	Path        string    `json:"path,omitempty"`
	Title       string    `json:"title,omitempty"`
	Referrer    string    `json:"referrer,omitempty"`
	UTMSource   string    `json:"utmSource,omitempty"`
	UTMMedium   string    `json:"utmMedium,omitempty"`
	UTMCampaign string    `json:"utmCampaign,omitempty"`
	DwellMs     int64     `json:"dwellMs,omitempty"`
	VisibleMs   int64     `json:"visibleMs,omitempty"`
	EngagedMs   int64     `json:"engagedMs,omitempty"`
	ScrollMax   int       `json:"scrollMaxPct,omitempty"`
	UA          string    `json:"ua,omitempty"`
	Locale      string    `json:"locale,omitempty"`
	Country     string    `json:"country,omitempty"` // coarse (CF-IPCountry etc.)
	Bot         bool      `json:"bot,omitempty"`
	MessageID   string    `json:"messageId,omitempty"`
}

// Store is an append-only JSONL file under dataDir/activity/events.jsonl.
type Store struct {
	mu   sync.Mutex
	path string
}

func NewStore(dataDir string) (*Store, error) {
	root := filepath.Join(dataDir, "activity")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(root, "events.jsonl")}, nil
}

func validType(t string) bool {
	switch t {
	case TypeLogin, TypeNavigate, TypeEngaged, TypePage, TypeIdentify, TypeAlias, TypeTrack:
		return true
	default:
		return false
	}
}

// Append writes one event (server stamps TS if zero).
func (s *Store) Append(ev Event) error {
	if s == nil {
		return nil
	}
	if ev.TS.IsZero() {
		ev.TS = time.Now().UTC()
	}
	ev.Type = strings.TrimSpace(ev.Type)
	if !validType(ev.Type) {
		return fmt.Errorf("invalid activity type %q", ev.Type)
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	b = append(b, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

// List returns the most recent events (newest first), capped at MaxListLimit.
func (s *Store) List(limit int) ([]Event, error) {
	if s == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Event{}, nil
		}
		return nil, err
	}
	defer f.Close()

	buf := make([]string, 0, limit)
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if len(buf) < limit {
			buf = append(buf, line)
			continue
		}
		copy(buf[0:], buf[1:])
		buf[len(buf)-1] = line
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	out := make([]Event, 0, len(buf))
	for i := len(buf) - 1; i >= 0; i-- {
		var ev Event
		if err := json.Unmarshal([]byte(buf[i]), &ev); err != nil {
			continue
		}
		out = append(out, ev)
	}
	return out, nil
}
