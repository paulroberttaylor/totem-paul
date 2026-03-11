package lesson

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed wordlists/common.txt
var commonWordsRaw string

// Stage defines a set of keys to practice.
type Stage struct {
	Name     string
	Keys     []string // typeable characters in this stage
	Pack     string   // empty for built-in stages
	Words    []string // nil = use CommonWords()
	Snippets []string // nil = use TerminalSnippets()
}

func HomeRow() Stage {
	return Stage{Name: "home_row", Keys: []string{"a", "r", "s", "t", "g", "m", "n", "e", "i", "o"}}
}

func TopRow() Stage {
	return Stage{Name: "top_row", Keys: []string{"q", "w", "f", "p", "b", "j", "l", "u", "y", ";"}}
}

func BottomRow() Stage {
	return Stage{Name: "bottom_row", Keys: []string{"z", "x", "c", "d", "v", "k", "h", ",", ".", "/"}}
}

func FullAlpha() Stage {
	s := Stage{Name: "full_alpha"}
	s.Keys = append(s.Keys, HomeRow().Keys...)
	s.Keys = append(s.Keys, TopRow().Keys...)
	s.Keys = append(s.Keys, BottomRow().Keys...)
	return s
}

func Numbers() Stage {
	return Stage{Name: "numbers", Keys: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}}
}

func Symbols() Stage {
	return Stage{Name: "symbols", Keys: []string{
		"!", "@", "#", "$", "%", "^", "&", "*", "~", "`", "_", "|",
		"'", "\"", "{", "}", "[", "]", "(", ")", "+", "-", "=",
	}}
}

func Mixed() Stage {
	s := Stage{Name: "mixed"}
	s.Keys = append(s.Keys, FullAlpha().Keys...)
	s.Keys = append(s.Keys, Numbers().Keys...)
	s.Keys = append(s.Keys, Symbols().Keys...)
	s.Keys = append(s.Keys, " ")
	return s
}

// AllStages returns all available stages in progression order.
func AllStages() []Stage {
	return []Stage{HomeRow(), TopRow(), BottomRow(), FullAlpha(), Numbers(), Symbols(), Mixed()}
}

// FilterWords returns only words where every character is in the allowed set.
func FilterWords(words []string, allowed map[string]bool) []string {
	var result []string
	for _, w := range words {
		ok := true
		for _, ch := range w {
			if !allowed[string(ch)] {
				ok = false
				break
			}
		}
		if ok && len(w) > 0 {
			result = append(result, w)
		}
	}
	return result
}

// GenerateExercise picks random words and joins them with spaces.
// wordCount is the target number of words.
func GenerateExercise(words []string, wordCount int) string {
	if len(words) == 0 {
		return ""
	}
	picked := make([]string, wordCount)
	for i := range picked {
		picked[i] = words[rand.Intn(len(words))]
	}
	return strings.Join(picked, " ")
}

// GenerateWeightedExercise biases word selection toward words containing weak keys.
func GenerateWeightedExercise(words []string, wordCount int, weakKeys map[string]float64) string {
	if len(words) == 0 {
		return ""
	}

	// Score each word by sum of weakness scores for its characters
	type scored struct {
		word   string
		weight float64
	}
	var pool []scored
	var totalWeight float64
	for _, w := range words {
		weight := 1.0
		for _, ch := range w {
			if wk, ok := weakKeys[string(ch)]; ok {
				weight += wk
			}
		}
		pool = append(pool, scored{word: w, weight: weight})
		totalWeight += weight
	}

	picked := make([]string, wordCount)
	for i := range picked {
		r := rand.Float64() * totalWeight
		cumulative := 0.0
		for _, s := range pool {
			cumulative += s.weight
			if r <= cumulative {
				picked[i] = s.word
				break
			}
		}
	}

	return strings.Join(picked, " ")
}

// BigramDrill generates a string that repeats the given bigram in short words.
func BigramDrill(bigram string, count int) string {
	parts := make([]string, count)
	for i := range parts {
		parts[i] = bigram
	}
	return strings.Join(parts, " ")
}

// TerminalSnippets returns common terminal commands for mixed-layer practice.
func TerminalSnippets() []string {
	return []string{
		"ls -la",
		"cd ~/projects",
		"git status",
		"git push origin main",
		"grep -r 'TODO' .",
		"cat ~/.bashrc",
		"echo $HOME",
		"mkdir -p src/lib",
		"rm -rf build/",
		"curl -s https://api",
		"docker ps -a",
		"ssh user@host",
		"tar -xzf archive.tar.gz",
		"find . -name '*.go'",
		"export PATH=$PATH:/usr/local/bin",
		"tmux new -s dev",
		"ctrl+c",
		"pip install -r requirements.txt",
		"npm run build && npm test",
		"chmod 755 script.sh",
	}
}

// SymbolSnippets returns exercises focused on pure SYM layer characters.
// Every character here is on the SYM layer — no base-layer letters or punctuation.
func SymbolSnippets() []string {
	return []string{
		// Quoting pairs
		"''", "\"\"", "' '", "\" \"",
		// Bracket pairs
		"{}", "[]", "()", "(())", "[[]]", "{{}}",
		"[()]", "{[]}", "({[]})",
		// Operators
		"+ -", "= +", "- =", "+ = -",
		"< >", "<= >=", "< = >",
		"+ - = < >",
		// Tilde and backtick
		"~ ~", "` `", "~ ` ~", "`` ~ ``",
		// Underscore and pipe
		"_ _", "| |", "_ | _", "| _ |",
		// Mixed drills
		"! @ # $ %", "^ & * ' \"",
		"~ ` _ | {", "} + - = <",
		"! ' \" ( )", "[ ] { } \\",
		// Realistic patterns (SYM chars only)
		"!= ==", "|| &&", "<< >>",
		"#{}", "${}", "!()", "@[]",
		"'_'", "\"_\"", "|>", "<|",
		"(*)", "[*]", "{*}",
		"+ = - < > ! @ # $ %",
		"^ & * ~ ` _ | \\ ' \"",
		"( ) [ ] { } + - = <",
	}
}

// AllStagesWithPacks returns built-in stages plus any stages from pack files in dir.
func AllStagesWithPacks(dir string) ([]Stage, []error) {
	stages := AllStages()
	packs, errs := LoadPacks(dir)
	for _, p := range packs {
		stages = append(stages, p.Stages()...)
	}
	return stages, errs
}

// CommonWords returns the embedded list of common English words.
func CommonWords() []string {
	var words []string
	for _, line := range strings.Split(commonWordsRaw, "\n") {
		w := strings.TrimSpace(line)
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}
