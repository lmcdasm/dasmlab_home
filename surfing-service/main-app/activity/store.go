// Package activity stores append-only login/navigation engagement events.
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

	DefaultListLimit = 200
	MaxListLimit     = 500
)

// Event is one activity log line (JSONL).
type Event struct {
	TS        time.Time `json:"ts"`
	Type      string    `json:"type"` // login | navigate | engaged
	User      string    `json:"user"` // preferred_username
	Sub       string    `json:"sub,omitempty"`
	Email     string    `json:"email,omitempty"`
	Path      string    `json:"path,omitempty"`
	DwellMs   int64     `json:"dwellMs,omitempty"`
	VisibleMs int64     `json:"visibleMs,omitempty"`
	EngagedMs int64     `json:"engagedMs,omitempty"`
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

// Append writes one event (server stamps TS if zero).
func (s *Store) Append(ev Event) error {
	if s == nil {
		return nil
	}
	if ev.TS.IsZero() {
		ev.TS = time.Now().UTC()
	}
	ev.Type = strings.TrimSpace(ev.Type)
	switch ev.Type {
	case TypeLogin, TypeNavigate, TypeEngaged:
	default:
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
