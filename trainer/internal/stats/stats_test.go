package stats

import (
	"path/filepath"
	"testing"
	"time"
)

func TestKeyStats_RecordHit(t *testing.T) {
	ks := &KeyStats{}
	ks.RecordHit(150 * time.Millisecond)
	ks.RecordHit(200 * time.Millisecond)
	ks.RecordMiss(300 * time.Millisecond)

	if ks.Hits != 2 {
		t.Errorf("want 2 hits, got %d", ks.Hits)
	}
	if ks.Misses != 1 {
		t.Errorf("want 1 miss, got %d", ks.Misses)
	}
	acc := ks.Accuracy()
	if acc < 0.66 || acc > 0.67 {
		t.Errorf("want accuracy ~0.667, got %f", acc)
	}
}

func TestSession_WPM(t *testing.T) {
	s := &Session{
		TotalChars: 250,
		Duration:   60 * time.Second,
	}
	// WPM = (chars / 5) / minutes = (250/5) / 1 = 50
	wpm := s.WPM()
	if wpm < 49.9 || wpm > 50.1 {
		t.Errorf("want WPM ~50, got %f", wpm)
	}
}

func TestHistory_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.jsonl")
	h := NewHistory(path)

	session := SessionRecord{
		Timestamp: time.Now(),
		Stage:     "home_row",
		Duration:  120.0,
		WPM:       34.5,
		Accuracy:  0.92,
		TotalKeys: 248,
		PerKey: map[string]KeyRecord{
			"a": {Hits: 30, Misses: 2, AvgMs: 180},
		},
	}

	err := h.Append(session)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}

	h2 := NewHistory(path)
	records, err := h2.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("want 1 record, got %d", len(records))
	}
	if records[0].Stage != "home_row" {
		t.Errorf("want stage 'home_row', got %q", records[0].Stage)
	}
	if records[0].PerKey["a"].Hits != 30 {
		t.Errorf("want 30 hits for 'a', got %d", records[0].PerKey["a"].Hits)
	}
}

func TestProgress_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	p := NewProgress(path)
	p.UnlockedStages = map[string]bool{"home_row": true, "top_row": true}
	p.KeyAccuracy = map[string]float64{"a": 0.97, "r": 0.88}

	err := p.Save()
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	p2 := NewProgress(path)
	err = p2.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !p2.UnlockedStages["top_row"] {
		t.Error("top_row should be unlocked")
	}
	if p2.KeyAccuracy["r"] != 0.88 {
		t.Errorf("want accuracy 0.88 for 'r', got %f", p2.KeyAccuracy["r"])
	}
}

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir should not be empty")
	}
	if !filepath.IsAbs(dir) {
		t.Error("ConfigDir should return absolute path")
	}
}
