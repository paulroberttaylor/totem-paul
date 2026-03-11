package lesson

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Pack represents a JSON stage pack loaded from disk.
type Pack struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Words       []string    `json:"words"`
	RawStages   []PackStage `json:"stages"`
}

// PackStage is a single stage within a pack.
type PackStage struct {
	Name     string   `json:"name"`
	Keys     []string `json:"keys"`
	Words    []string `json:"words"`
	Snippets []string `json:"snippets"`
}

// LoadPack reads and validates a single pack JSON file.
func LoadPack(path string) (*Pack, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading pack %s: %w", path, err)
	}

	var p Pack
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing pack %s: %w", path, err)
	}

	if p.Name == "" {
		return nil, fmt.Errorf("pack %s: name is required", path)
	}

	for i, s := range p.RawStages {
		if len(s.Keys) == 0 {
			return nil, fmt.Errorf("pack %s: stage %d (%q) has no keys", path, i, s.Name)
		}
	}

	return &p, nil
}

// LoadPacks scans dir for *.json files and loads each as a pack.
// If dir does not exist, it returns nil, nil.
// Individual pack errors are collected but do not stop loading.
func LoadPacks(dir string) ([]Pack, []error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, []error{fmt.Errorf("reading pack directory: %w", err)}
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)

	var packs []Pack
	var errs []error
	for _, f := range files {
		p, err := LoadPack(f)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		packs = append(packs, *p)
	}
	return packs, errs
}

// Stages converts a Pack into a slice of Stage values.
// Stage words inherit from the pack if not set on the stage.
func (p *Pack) Stages() []Stage {
	stages := make([]Stage, len(p.RawStages))
	for i, ps := range p.RawStages {
		words := ps.Words
		if len(words) == 0 {
			words = p.Words
		}
		stages[i] = Stage{
			Name:     ps.Name,
			Keys:     ps.Keys,
			Pack:     p.Name,
			Words:    words,
			Snippets: ps.Snippets,
		}
	}
	return stages
}
