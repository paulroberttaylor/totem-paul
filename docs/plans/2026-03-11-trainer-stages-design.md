# Trainer Stages: Apex Packs Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make Apex pack snippets work as first-class "type this code" exercises and ensure packs are discoverable without manual copying.

**Architecture:** The pack system (`pack.go`) and `apex.json` already exist. Two fixes needed: (1) snippet-first exercise generation, (2) embed repo packs so they're always available.

**Tech Stack:** Go, Bubble Tea, embed

---

### Task 1: Fix snippet-first exercise generation

Currently `app.go:104-134` only uses snippets as a fallback. Stages with snippets should use them directly.

**Files:**
- Modify: `trainer/internal/ui/app.go:104-134`
- Test: `trainer/internal/lesson/integration_test.go`

**Step 1: Write the failing test**

Add to `trainer/internal/lesson/integration_test.go`:

```go
func TestSnippetStageUsesSnippets(t *testing.T) {
	stage := lesson.Stage{
		Name:     "test_snippets",
		Keys:     []string{"a", "b"},
		Snippets: []string{"abc def", "ghi jkl"},
	}
	// A stage with snippets should pick from snippets, not generate from keys/words
	found := false
	for i := 0; i < 20; i++ {
		ex := lesson.GenerateSnippetExercise(stage.Snippets)
		if ex == "abc def" || ex == "ghi jkl" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected exercise to be one of the snippets")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd trainer && go test ./internal/lesson/ -run TestSnippetStage -v`
Expected: FAIL — `GenerateSnippetExercise` not defined

**Step 3: Add GenerateSnippetExercise to lesson.go**

Add to `trainer/internal/lesson/lesson.go`:

```go
// GenerateSnippetExercise picks a random snippet from the list.
func GenerateSnippetExercise(snippets []string) string {
	if len(snippets) == 0 {
		return ""
	}
	return snippets[rand.Intn(len(snippets))]
}
```

**Step 4: Run test to verify it passes**

Run: `cd trainer && go test ./internal/lesson/ -run TestSnippetStage -v`
Expected: PASS

**Step 5: Update app.go exercise generation**

In `trainer/internal/ui/app.go`, replace the `startLessonMsg` handler (lines ~104-134) so stages with snippets use them first:

```go
case startLessonMsg:
	stage := msg.stage
	var exercise string

	if len(stage.Snippets) > 0 {
		// Snippet-first: pick a random snippet to type verbatim
		exercise = lesson.GenerateSnippetExercise(stage.Snippets)
	} else {
		switch stage.Name {
		case "symbols":
			snippets := lesson.SymbolSnippets()
			exercise = lesson.GenerateExercise(snippets, 10)
		case "numbers":
			exercise = lesson.GenerateExercise(stage.Keys, 20)
		default:
			allowed := make(map[string]bool)
			for _, k := range stage.Keys {
				allowed[k] = true
			}
			allowed[" "] = true
			var wordPool []string
			if len(stage.Words) > 0 {
				wordPool = stage.Words
			} else {
				wordPool = lesson.CommonWords()
			}
			words := lesson.FilterWords(wordPool, allowed)
			exercise = lesson.GenerateExercise(words, 15)
			if exercise == "" {
				exercise = lesson.GenerateExercise(stage.Keys, 20)
			}
		}
	}
	a.typing = newTypingModel(exercise, stage.Name)
	a.current = screenTyping
	return a, nil
```

**Step 6: Run all tests**

Run: `cd trainer && go test ./... -v`
Expected: All PASS

**Step 7: Commit**

```bash
git add trainer/internal/lesson/lesson.go trainer/internal/lesson/integration_test.go trainer/internal/ui/app.go
git commit -m "feat: snippet-first exercise generation for pack stages"
```

---

### Task 2: Embed repo packs so they're always available

Currently packs only load from `~/.config/totem-trainer/packs/`. Embed the repo's packs as defaults.

**Files:**
- Modify: `trainer/internal/lesson/pack.go`
- Modify: `trainer/internal/lesson/lesson.go` (AllStagesWithPacks)
- Test: `trainer/internal/lesson/pack_test.go`

**Step 1: Write the failing test**

Add to `trainer/internal/lesson/pack_test.go`:

```go
func TestEmbeddedPacksLoad(t *testing.T) {
	packs := lesson.EmbeddedPacks()
	if len(packs) == 0 {
		t.Fatal("expected at least one embedded pack")
	}
	found := false
	for _, p := range packs {
		if p.Name == "Apex" {
			found = true
			if len(p.RawStages) < 2 {
				t.Errorf("Apex pack has %d stages, want >= 2", len(p.RawStages))
			}
		}
	}
	if !found {
		t.Fatal("expected Apex embedded pack")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd trainer && go test ./internal/lesson/ -run TestEmbeddedPacks -v`
Expected: FAIL — `EmbeddedPacks` not defined

**Step 3: Add embedded packs to pack.go**

Add to `trainer/internal/lesson/pack.go`:

```go
import (
	"embed"
	// ... existing imports
)

//go:embed packs/*.json
var embeddedPacksFS embed.FS

// EmbeddedPacks loads all packs embedded in the binary.
func EmbeddedPacks() []Pack {
	entries, err := embeddedPacksFS.ReadDir("packs")
	if err != nil {
		return nil
	}
	var packs []Pack
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := embeddedPacksFS.ReadFile("packs/" + e.Name())
		if err != nil {
			continue
		}
		var p Pack
		if err := json.Unmarshal(data, &p); err != nil || p.Name == "" {
			continue
		}
		packs = append(packs, p)
	}
	return packs
}
```

**Step 4: Update AllStagesWithPacks to include embedded packs**

In `trainer/internal/lesson/lesson.go`, update `AllStagesWithPacks`:

```go
func AllStagesWithPacks(dir string) ([]Stage, []error) {
	stages := AllStages()

	// Load embedded packs first (always available)
	for _, p := range EmbeddedPacks() {
		stages = append(stages, p.Stages()...)
	}

	// Then load user packs from disk (can override/extend)
	packs, errs := LoadPacks(dir)
	for _, p := range packs {
		stages = append(stages, p.Stages()...)
	}
	return stages, errs
}
```

**Step 5: Run test to verify it passes**

Run: `cd trainer && go test ./internal/lesson/ -run TestEmbeddedPacks -v`
Expected: PASS

**Step 6: Run all tests**

Run: `cd trainer && go test ./... -v`
Expected: All PASS

**Step 7: Commit**

```bash
git add trainer/internal/lesson/pack.go trainer/internal/lesson/lesson.go trainer/internal/lesson/pack_test.go
git commit -m "feat: embed repo packs so Apex stages are always available"
```

---

### Task 3: Expand Apex snippets

Add more realistic Apex code patterns to `apex.json`.

**Files:**
- Modify: `trainer/internal/lesson/packs/apex.json`

**Step 1: Add more snippets**

Expand the snippets list to 20+ patterns covering:
- Class declarations with sharing settings
- SOQL with WHERE, LIMIT, ORDER BY
- DML with error handling
- Test methods with @isTest and assertions
- Trigger patterns with context variables
- Batch and queueable patterns
- Governor limit patterns

**Step 2: Run all tests to verify JSON is valid**

Run: `cd trainer && go test ./internal/lesson/ -v`
Expected: All PASS (embedded pack parses correctly)

**Step 3: Commit**

```bash
git add trainer/internal/lesson/packs/apex.json
git commit -m "content: expand Apex snippets to 20+ real code patterns"
```

---

### Task 4: Manual smoke test

Run: `cd trainer && go run ./cmd/totem-trainer/`

Verify:
1. Picker shows "── Apex ──" group header
2. Selecting "snippets" stage shows a full Apex line to type
3. Typing it scores correctly (green/red, WPM, accuracy)
4. Selecting "keywords" stage generates words from Apex vocabulary
