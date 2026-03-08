package lesson

import (
	"testing"
)

func TestStageKeys(t *testing.T) {
	s := HomeRow()
	keys := s.Keys
	// Colemak DH home row
	want := []string{"a", "r", "s", "t", "d", "h", "n", "e", "i", "o"}
	if len(keys) != len(want) {
		t.Fatalf("home row: want %d keys, got %d", len(want), len(keys))
	}
	for i, k := range want {
		if keys[i] != k {
			t.Errorf("home row[%d]: want %q, got %q", i, k, keys[i])
		}
	}
}

func TestAllStages(t *testing.T) {
	stages := AllStages()
	if len(stages) != 7 {
		t.Fatalf("want 7 stages, got %d", len(stages))
	}
	wantNames := []string{"home_row", "top_row", "bottom_row", "full_alpha", "numbers", "symbols", "mixed"}
	for i, name := range wantNames {
		if stages[i].Name != name {
			t.Errorf("stage[%d]: want name %q, got %q", i, name, stages[i].Name)
		}
	}
}

func TestFilterWords(t *testing.T) {
	words := []string{"the", "and", "hello", "test", "art", "rate"}
	allowed := map[string]bool{"t": true, "h": true, "e": true, "a": true, "r": true}

	got := FilterWords(words, allowed)

	// Only words where every char is in allowed set
	// "the" -> t,h,e -> yes
	// "and" -> a,n,d -> no (n,d not allowed)
	// "hello" -> no
	// "test" -> no (s not allowed)
	// "art" -> a,r,t -> yes
	// "rate" -> r,a,t,e -> yes
	if len(got) != 3 {
		t.Fatalf("want 3 filtered words, got %d: %v", len(got), got)
	}
}

func TestFilterWords_Empty(t *testing.T) {
	got := FilterWords(nil, map[string]bool{"a": true})
	if len(got) != 0 {
		t.Errorf("want 0 filtered words from nil input, got %d", len(got))
	}
}

func TestGenerateExercise_OnlyUsesAllowedKeys(t *testing.T) {
	words := []string{"the", "and", "art", "rate", "ear", "tear", "heat"}
	allowed := map[string]bool{"t": true, "h": true, "e": true, "a": true, "r": true}

	exercise := GenerateExercise(FilterWords(words, allowed), 5)

	for _, ch := range exercise {
		if ch == ' ' {
			continue
		}
		if !allowed[string(ch)] {
			t.Errorf("exercise contains disallowed char %q in %q", string(ch), exercise)
		}
	}
}

func TestGenerateExercise_Empty(t *testing.T) {
	result := GenerateExercise(nil, 5)
	if result != "" {
		t.Errorf("want empty string for nil words, got %q", result)
	}
}

func TestGenerateWeightedExercise(t *testing.T) {
	words := []string{"art", "the", "rate"}
	weakKeys := map[string]float64{"r": 5.0}

	exercise := GenerateWeightedExercise(words, 10, weakKeys)
	if exercise == "" {
		t.Error("weighted exercise should not be empty")
	}
	// Should produce space-separated words
	parts := splitWords(exercise)
	if len(parts) != 10 {
		t.Errorf("want 10 words, got %d", len(parts))
	}
}

func TestGenerateWeightedExercise_Empty(t *testing.T) {
	result := GenerateWeightedExercise(nil, 5, nil)
	if result != "" {
		t.Errorf("want empty string for nil words, got %q", result)
	}
}

func TestBigramDrill(t *testing.T) {
	drill := BigramDrill("st", 10)
	if len(drill) == 0 {
		t.Error("bigram drill should not be empty")
	}
	// Should contain the target bigram
	if !containsBigram(drill, "st") {
		t.Errorf("drill %q should contain bigram 'st'", drill)
	}
}

func TestTerminalSnippets(t *testing.T) {
	snippets := TerminalSnippets()
	if len(snippets) == 0 {
		t.Error("terminal snippets should not be empty")
	}
}

func TestCommonWords(t *testing.T) {
	words := CommonWords()
	if len(words) < 100 {
		t.Errorf("want at least 100 common words, got %d", len(words))
	}
	// All words should be lowercase and non-empty
	for _, w := range words {
		if w == "" {
			t.Error("found empty word in common words")
		}
		for _, ch := range w {
			if ch >= 'A' && ch <= 'Z' {
				t.Errorf("word %q contains uppercase letter", w)
			}
		}
	}
}

func TestCommonWords_FilterableByHomeRow(t *testing.T) {
	words := CommonWords()
	homeKeys := HomeRow().Keys
	allowed := make(map[string]bool)
	for _, k := range homeKeys {
		allowed[k] = true
	}
	filtered := FilterWords(words, allowed)
	// Should have at least some words typeable with home row only
	if len(filtered) < 5 {
		t.Errorf("want at least 5 home-row-only words, got %d: %v", len(filtered), filtered)
	}
}

func containsBigram(s, bg string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i:i+2] == bg {
			return true
		}
	}
	return false
}

func splitWords(s string) []string {
	var words []string
	current := ""
	for _, ch := range s {
		if ch == ' ' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
