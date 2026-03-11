package lesson

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPack_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	data := `{
		"name": "TestPack",
		"description": "A test pack",
		"words": ["hello", "world"],
		"stages": [
			{
				"name": "basics",
				"keys": ["a", "b", "c"],
				"words": ["abc", "cab"]
			},
			{
				"name": "snippets",
				"keys": ["a", "b"],
				"snippets": ["a = b;"]
			}
		]
	}`
	os.WriteFile(path, []byte(data), 0644)

	p, err := LoadPack(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "TestPack" {
		t.Errorf("name = %q, want %q", p.Name, "TestPack")
	}
	if len(p.RawStages) != 2 {
		t.Fatalf("got %d stages, want 2", len(p.RawStages))
	}
	if p.RawStages[0].Name != "basics" {
		t.Errorf("stage 0 name = %q, want %q", p.RawStages[0].Name, "basics")
	}
	if len(p.RawStages[1].Snippets) != 1 {
		t.Errorf("stage 1 snippets = %d, want 1", len(p.RawStages[1].Snippets))
	}
}

func TestLoadPack_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte(`{not json`), 0644)

	_, err := LoadPack(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadPack_MissingName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "noname.json")
	data := `{"stages": [{"name": "s1", "keys": ["a"]}]}`
	os.WriteFile(path, []byte(data), 0644)

	_, err := LoadPack(path)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoadPack_StageNoKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nokeys.json")
	data := `{"name": "Test", "stages": [{"name": "empty"}]}`
	os.WriteFile(path, []byte(data), 0644)

	_, err := LoadPack(path)
	if err == nil {
		t.Fatal("expected error for stage with no keys")
	}
}

func TestLoadPack_StageInheritsPackWords(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inherit.json")
	data := `{
		"name": "Inherit",
		"words": ["pack", "level"],
		"stages": [
			{"name": "s1", "keys": ["a"]},
			{"name": "s2", "keys": ["b"], "words": ["stage", "level"]}
		]
	}`
	os.WriteFile(path, []byte(data), 0644)

	p, err := LoadPack(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stages := p.Stages()
	// s1 should inherit pack words
	if len(stages[0].Words) != 2 || stages[0].Words[0] != "pack" {
		t.Errorf("stage 0 words = %v, want [pack level]", stages[0].Words)
	}
	// s2 should use its own words
	if len(stages[1].Words) != 2 || stages[1].Words[0] != "stage" {
		t.Errorf("stage 1 words = %v, want [stage level]", stages[1].Words)
	}
	// Both should have Pack set
	if stages[0].Pack != "Inherit" {
		t.Errorf("stage 0 pack = %q, want %q", stages[0].Pack, "Inherit")
	}
}

func TestLoadPacks_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	packs, errs := LoadPacks(dir)
	if len(packs) != 0 {
		t.Errorf("expected no packs, got %d", len(packs))
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestLoadPacks_NonexistentDir(t *testing.T) {
	packs, errs := LoadPacks("/nonexistent/path/that/does/not/exist")
	if len(packs) != 0 {
		t.Errorf("expected no packs, got %d", len(packs))
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}
