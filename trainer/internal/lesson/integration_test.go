package lesson

import (
	"path/filepath"
	"testing"

	"github.com/paul/totem-trainer/internal/keymap"
)

func TestIntegration_KeymapToLesson(t *testing.T) {
	// Parse the actual keymap
	path, err := filepath.Abs(filepath.Join("..", "..", "..", "config", "totem.keymap"))
	if err != nil {
		t.Fatal(err)
	}
	km, err := keymap.ParseFile(path)
	if err != nil {
		t.Skipf("keymap not available: %v", err)
	}

	// Extract typeable chars from BASE layer
	var baseChars []string
	for _, k := range km.Layers[0].Bindings {
		if k.Char != "" && k.Char != " " && k.Char != "\n" && k.Char != "\t" {
			baseChars = append(baseChars, k.Char)
		}
	}

	if len(baseChars) < 20 {
		t.Errorf("BASE layer should have 20+ typeable chars, got %d", len(baseChars))
	}

	// Verify we can generate an exercise from home row
	homeRow := HomeRow()
	allowed := make(map[string]bool)
	for _, k := range homeRow.Keys {
		allowed[k] = true
	}
	words := CommonWords()
	filtered := FilterWords(words, allowed)
	if len(filtered) == 0 {
		t.Error("no words match home row keys — wordlist may be wrong")
	}

	exercise := GenerateExercise(filtered, 10)
	if exercise == "" {
		t.Error("exercise should not be empty")
	}
	t.Logf("home row exercise: %q", exercise)

	// Verify NUM layer has numbers (layer index 2)
	num := km.Layers[2]
	var nums []string
	for _, k := range num.Bindings {
		if k.Char >= "0" && k.Char <= "9" {
			nums = append(nums, k.Char)
		}
	}
	if len(nums) != 10 {
		t.Errorf("NUM layer should have 10 numbers, got %d: %v", len(nums), nums)
	}
}
