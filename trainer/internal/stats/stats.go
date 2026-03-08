package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// KeyStats tracks live per-key performance during a session.
type KeyStats struct {
	Hits    int
	Misses  int
	TotalMs float64
}

func (ks *KeyStats) RecordHit(latency time.Duration) {
	ks.Hits++
	ks.TotalMs += float64(latency.Milliseconds())
}

func (ks *KeyStats) RecordMiss(latency time.Duration) {
	ks.Misses++
	ks.TotalMs += float64(latency.Milliseconds())
}

func (ks *KeyStats) Accuracy() float64 {
	total := ks.Hits + ks.Misses
	if total == 0 {
		return 0
	}
	return float64(ks.Hits) / float64(total)
}

func (ks *KeyStats) AvgLatency() float64 {
	total := ks.Hits + ks.Misses
	if total == 0 {
		return 0
	}
	return ks.TotalMs / float64(total)
}

// Session tracks live session state.
type Session struct {
	TotalChars int
	Correct    int
	Duration   time.Duration
	PerKey     map[string]*KeyStats
}

func NewSession() *Session {
	return &Session{PerKey: make(map[string]*KeyStats)}
}

func (s *Session) WPM() float64 {
	minutes := s.Duration.Minutes()
	if minutes == 0 {
		return 0
	}
	return (float64(s.TotalChars) / 5.0) / minutes
}

func (s *Session) Accuracy() float64 {
	if s.TotalChars == 0 {
		return 0
	}
	return float64(s.Correct) / float64(s.TotalChars)
}

// Persistence types (JSON serialization).

type KeyRecord struct {
	Hits   int     `json:"hits"`
	Misses int     `json:"misses"`
	AvgMs  float64 `json:"avg_ms"`
}

type SessionRecord struct {
	Timestamp time.Time            `json:"timestamp"`
	Stage     string               `json:"stage"`
	Duration  float64              `json:"duration_secs"`
	WPM       float64              `json:"wpm"`
	Accuracy  float64              `json:"accuracy"`
	TotalKeys int                  `json:"total_keys"`
	PerKey    map[string]KeyRecord `json:"per_key"`
}

// History manages the append-only session log.
type History struct {
	path string
}

func NewHistory(path string) *History {
	return &History{path: path}
}

func (h *History) Append(record SessionRecord) error {
	if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(h.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = f.Write(data)
	return err
}

func (h *History) Load() ([]SessionRecord, error) {
	data, err := os.ReadFile(h.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read history: %w", err)
	}
	var records []SessionRecord
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var r SessionRecord
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			continue // skip corrupt lines
		}
		records = append(records, r)
	}
	return records, nil
}

// Progress tracks cumulative state across sessions.
type Progress struct {
	path           string
	UnlockedStages map[string]bool    `json:"unlocked_stages"`
	KeyAccuracy    map[string]float64 `json:"key_accuracy"`
}

func NewProgress(path string) *Progress {
	return &Progress{
		path:           path,
		UnlockedStages: make(map[string]bool),
		KeyAccuracy:    make(map[string]float64),
	}
}

func (p *Progress) Save() error {
	if err := os.MkdirAll(filepath.Dir(p.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.path, data, 0o644)
}

func (p *Progress) Load() error {
	data, err := os.ReadFile(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read progress: %w", err)
	}
	return json.Unmarshal(data, p)
}

// ConfigDir returns the path to the totem-trainer config directory.
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "totem-trainer")
}
