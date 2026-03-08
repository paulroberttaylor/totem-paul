# Totem Trainer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a terminal typing trainer for the Totem 38-key split keyboard that parses the ZMK keymap, provides progressive lessons by layer, adapts to weak keys, and tracks stats over time.

**Architecture:** Go module at `trainer/` with four packages: `keymap` (ZMK parser), `lesson` (exercise generation + adaptive algorithm), `stats` (persistence + aggregation), `ui` (Bubble Tea TUI). The keymap parser reads `config/totem.keymap` as source of truth. Stats persist to `~/.config/totem-trainer/`.

**Tech Stack:** Go 1.26, Bubble Tea, lipgloss, bubbles, standard library for file I/O and JSON.

---

### Task 1: Scaffold Go module and verify build

**Files:**
- Create: `trainer/cmd/totem-trainer/main.go`
- Create: `trainer/go.mod`

**Step 1: Initialize Go module**

```bash
cd trainer && go mod init github.com/paul/totem-trainer
```

**Step 2: Write minimal main.go**

```go
package main

import "fmt"

func main() {
    fmt.Println("totem-trainer")
}
```

**Step 3: Verify it builds and runs**

Run: `cd trainer && go build ./cmd/totem-trainer && ./totem-trainer`
Expected: prints "totem-trainer"

**Step 4: Add Bubble Tea dependency**

```bash
cd trainer && go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/charmbracelet/bubbles
```

**Step 5: Commit**

```bash
git add trainer/
git commit -m "feat: scaffold totem-trainer Go module with dependencies"
```

---

### Task 2: Keymap data model

**Files:**
- Create: `trainer/internal/keymap/model.go`
- Create: `trainer/internal/keymap/model_test.go`

**Step 1: Write the test**

```go
package keymap

import "testing"

func TestKeyPosition_Hand(t *testing.T) {
    layout := NewTotemLayout()

    tests := []struct {
        pos  int
        hand string
        row  int
    }{
        {0, "left", 0},   // Q
        {5, "right", 0},  // J
        {10, "left", 1},  // A
        {20, "left", 2},  // ESC (pinky)
        {31, "right", 2}, // backslash (pinky)
        {32, "left", 3},  // DEL (thumb)
        {37, "right", 3}, // BSPC (thumb)
    }

    for _, tt := range tests {
        info := layout.Position(tt.pos)
        if info.Hand != tt.hand {
            t.Errorf("pos %d: want hand %s, got %s", tt.pos, tt.hand, info.Hand)
        }
        if info.Row != tt.row {
            t.Errorf("pos %d: want row %d, got %d", tt.pos, tt.row, info.Row)
        }
    }
}

func TestZMKCodeToChar(t *testing.T) {
    tests := []struct {
        code string
        char string
    }{
        {"A", "a"},
        {"SEMI", ";"},
        {"COMMA", ","},
        {"DOT", "."},
        {"FSLH", "/"},
        {"BSLH", "\\"},
        {"N7", "7"},
        {"EXCL", "!"},
        {"AT", "@"},
        {"HASH", "#"},
        {"DLLR", "$"},
        {"PRCNT", "%"},
        {"CARET", "^"},
        {"AMPS", "&"},
        {"ASTRK", "*"},
        {"SQT", "'"},
        {"DQT", "\""},
        {"TILDE", "~"},
        {"GRAVE", "`"},
        {"UNDER", "_"},
        {"PIPE", "|"},
        {"LBRC", "{"},
        {"RBRC", "}"},
        {"LBKT", "["},
        {"RBKT", "]"},
        {"LPAR", "("},
        {"RPAR", ")"},
        {"PLUS", "+"},
        {"MINUS", "-"},
        {"EQUAL", "="},
    }

    for _, tt := range tests {
        got := ZMKCodeToChar(tt.code)
        if got != tt.char {
            t.Errorf("ZMKCodeToChar(%q) = %q, want %q", tt.code, got, tt.char)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd trainer && go test ./internal/keymap/ -v`
Expected: FAIL — types don't exist yet

**Step 3: Write model.go**

```go
package keymap

// Key represents a single key binding on a layer.
type Key struct {
    Position int    // 0-37 physical position
    Type     string // "kp", "mt", "lt", "trans", "none", "sys", "bt", "out"
    Tap      string // character produced on tap (ZMK code)
    Hold     string // modifier or layer on hold
    Char     string // resolved typeable character (empty if not typeable)
}

// Layer represents one layer of the keymap.
type Layer struct {
    Name     string
    Index    int
    Bindings []Key // exactly 38 entries
}

// Combo represents a ZMK combo.
type Combo struct {
    Name      string
    Positions []int
    Output    string // resolved character
}

// Keymap is the full parsed keymap.
type Keymap struct {
    Layers []Layer
    Combos []Combo
}

// PositionInfo describes the physical location of a key.
type PositionInfo struct {
    Hand   string // "left" or "right"
    Row    int    // 0=top, 1=home, 2=bottom, 3=thumb
    Col    int    // column within the hand
    Finger string // "pinky", "ring", "middle", "index", "thumb"
}

// TotemLayout maps physical positions to layout metadata.
type TotemLayout struct {
    positions [38]PositionInfo
}

func NewTotemLayout() *TotemLayout {
    l := &TotemLayout{}

    // Row 0: positions 0-4 (left), 5-9 (right)
    fingers := []string{"pinky", "ring", "middle", "index", "index"}
    for i := 0; i < 5; i++ {
        l.positions[i] = PositionInfo{Hand: "left", Row: 0, Col: i, Finger: fingers[i]}
        l.positions[i+5] = PositionInfo{Hand: "right", Row: 0, Col: i, Finger: fingers[4-i]}
    }

    // Row 1: positions 10-14 (left), 15-19 (right)
    for i := 0; i < 5; i++ {
        l.positions[10+i] = PositionInfo{Hand: "left", Row: 1, Col: i, Finger: fingers[i]}
        l.positions[15+i] = PositionInfo{Hand: "right", Row: 1, Col: i, Finger: fingers[4-i]}
    }

    // Row 2: position 20 (left pinky extra), 21-25 (left), 26-30 (right), 31 (right pinky extra)
    l.positions[20] = PositionInfo{Hand: "left", Row: 2, Col: -1, Finger: "pinky"}
    for i := 0; i < 5; i++ {
        l.positions[21+i] = PositionInfo{Hand: "left", Row: 2, Col: i, Finger: fingers[i]}
        l.positions[26+i] = PositionInfo{Hand: "right", Row: 2, Col: i, Finger: fingers[4-i]}
    }
    l.positions[31] = PositionInfo{Hand: "right", Row: 2, Col: 5, Finger: "pinky"}

    // Row 3: thumbs 32-34 (left), 35-37 (right)
    for i := 0; i < 3; i++ {
        l.positions[32+i] = PositionInfo{Hand: "left", Row: 3, Col: i, Finger: "thumb"}
        l.positions[35+i] = PositionInfo{Hand: "right", Row: 3, Col: i, Finger: "thumb"}
    }

    return l
}

func (l *TotemLayout) Position(pos int) PositionInfo {
    if pos < 0 || pos >= 38 {
        return PositionInfo{}
    }
    return l.positions[pos]
}

// ZMK code to typeable character mapping.
var zmkCharMap = map[string]string{
    // Letters (single uppercase letter → lowercase char)
    "A": "a", "B": "b", "C": "c", "D": "d", "E": "e",
    "F": "f", "G": "g", "H": "h", "I": "i", "J": "j",
    "K": "k", "L": "l", "M": "m", "N": "n", "O": "o",
    "P": "p", "Q": "q", "R": "r", "S": "s", "T": "t",
    "U": "u", "V": "v", "W": "w", "X": "x", "Y": "y",
    "Z": "z",

    // Numbers
    "N0": "0", "N1": "1", "N2": "2", "N3": "3", "N4": "4",
    "N5": "5", "N6": "6", "N7": "7", "N8": "8", "N9": "9",

    // Punctuation
    "SEMI": ";", "COMMA": ",", "DOT": ".", "FSLH": "/", "BSLH": "\\",
    "SQT": "'", "DQT": "\"", "GRAVE": "`", "TILDE": "~",
    "MINUS": "-", "UNDER": "_", "EQUAL": "=", "PLUS": "+",

    // Brackets
    "LBKT": "[", "RBKT": "]", "LBRC": "{", "RBRC": "}",
    "LPAR": "(", "RPAR": ")",

    // Symbols
    "EXCL": "!", "AT": "@", "HASH": "#", "DLLR": "$",
    "PRCNT": "%", "CARET": "^", "AMPS": "&", "ASTRK": "*",
    "PIPE": "|",

    // Whitespace / control
    "SPACE": " ", "RET": "\n", "TAB": "\t",
}

// ZMKCodeToChar converts a ZMK keycode to its typeable character.
// Returns empty string if the key doesn't produce a typeable character.
func ZMKCodeToChar(code string) string {
    if ch, ok := zmkCharMap[code]; ok {
        return ch
    }
    return ""
}
```

**Step 4: Run tests**

Run: `cd trainer && go test ./internal/keymap/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add trainer/internal/keymap/
git commit -m "feat: add keymap data model with Totem layout and ZMK code mapping"
```

---

### Task 3: ZMK keymap parser

**Files:**
- Create: `trainer/internal/keymap/parser.go`
- Create: `trainer/internal/keymap/parser_test.go`

**Step 1: Write the test**

The test reads the actual `config/totem.keymap` file from the repo.

```go
package keymap

import (
    "os"
    "path/filepath"
    "testing"
)

func testKeymapPath(t *testing.T) string {
    t.Helper()
    // Walk up from trainer/internal/keymap/ to repo root
    path := filepath.Join("..", "..", "..", "config", "totem.keymap")
    abs, err := filepath.Abs(path)
    if err != nil {
        t.Fatalf("resolve keymap path: %v", err)
    }
    if _, err := os.Stat(abs); err != nil {
        t.Skipf("keymap not found at %s", abs)
    }
    return abs
}

func TestParseKeymap_LayerCount(t *testing.T) {
    km, err := ParseFile(testKeymapPath(t))
    if err != nil {
        t.Fatalf("ParseFile: %v", err)
    }
    if len(km.Layers) != 4 {
        t.Errorf("want 4 layers, got %d", len(km.Layers))
    }
}

func TestParseKeymap_LayerNames(t *testing.T) {
    km, err := ParseFile(testKeymapPath(t))
    if err != nil {
        t.Fatalf("ParseFile: %v", err)
    }
    want := []string{"BASE", "NAV", "SYM", "ADJ"}
    for i, name := range want {
        if km.Layers[i].Name != name {
            t.Errorf("layer %d: want name %q, got %q", i, name, km.Layers[i].Name)
        }
    }
}

func TestParseKeymap_BaseLayerKeys(t *testing.T) {
    km, err := ParseFile(testKeymapPath(t))
    if err != nil {
        t.Fatalf("ParseFile: %v", err)
    }
    base := km.Layers[0]
    if len(base.Bindings) != 38 {
        t.Fatalf("BASE layer: want 38 bindings, got %d", len(base.Bindings))
    }

    tests := []struct {
        pos      int
        typ      string
        tap      string
        hold     string
        char     string
    }{
        {0, "kp", "Q", "", "q"},          // Q
        {9, "kp", "SEMI", "", ";"},        // ;
        {10, "mt", "A", "LGUI", "a"},      // A/GUI
        {13, "mt", "T", "LSHFT", "t"},     // T/SHFT
        {20, "kp", "ESC", "", ""},         // ESC (not typeable)
        {33, "lt", "TAB", "NAV", ""},      // TAB/NAV (layer tap, not simple char)
        {34, "kp", "SPACE", "", " "},      // SPACE
    }

    for _, tt := range tests {
        k := base.Bindings[tt.pos]
        if k.Type != tt.typ {
            t.Errorf("pos %d: want type %q, got %q", tt.pos, tt.typ, k.Type)
        }
        if k.Tap != tt.tap {
            t.Errorf("pos %d: want tap %q, got %q", tt.pos, tt.tap, k.Tap)
        }
        if k.Hold != tt.hold {
            t.Errorf("pos %d: want hold %q, got %q", tt.pos, tt.hold, k.Hold)
        }
        if k.Char != tt.char {
            t.Errorf("pos %d: want char %q, got %q", tt.pos, tt.char, k.Char)
        }
    }
}

func TestParseKeymap_NAVLayerNumbers(t *testing.T) {
    km, err := ParseFile(testKeymapPath(t))
    if err != nil {
        t.Fatalf("ParseFile: %v", err)
    }
    nav := km.Layers[1]

    // Position 6 should be 7 (N7)
    if nav.Bindings[6].Char != "7" {
        t.Errorf("NAV pos 6: want char '7', got %q", nav.Bindings[6].Char)
    }
    // Position 36 should be 0 (N0)
    if nav.Bindings[36].Char != "0" {
        t.Errorf("NAV pos 36: want char '0', got %q", nav.Bindings[36].Char)
    }
}

func TestParseKeymap_SYMLayerSymbols(t *testing.T) {
    km, err := ParseFile(testKeymapPath(t))
    if err != nil {
        t.Fatalf("ParseFile: %v", err)
    }
    sym := km.Layers[2]

    // Position 0 should be !
    if sym.Bindings[0].Char != "!" {
        t.Errorf("SYM pos 0: want char '!', got %q", sym.Bindings[0].Char)
    }
    // Position 3 should be $
    if sym.Bindings[3].Char != "$" {
        t.Errorf("SYM pos 3: want char '$', got %q", sym.Bindings[3].Char)
    }
}

func TestParseKeymap_Combos(t *testing.T) {
    km, err := ParseFile(testKeymapPath(t))
    if err != nil {
        t.Fatalf("ParseFile: %v", err)
    }
    if len(km.Combos) != 1 {
        t.Fatalf("want 1 combo, got %d", len(km.Combos))
    }
    c := km.Combos[0]
    if c.Name != "combo_esc" {
        t.Errorf("combo name: want 'combo_esc', got %q", c.Name)
    }
    if len(c.Positions) != 2 || c.Positions[0] != 0 || c.Positions[1] != 1 {
        t.Errorf("combo positions: want [0,1], got %v", c.Positions)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd trainer && go test ./internal/keymap/ -v -run TestParse`
Expected: FAIL — ParseFile doesn't exist

**Step 3: Write parser.go**

```go
package keymap

import (
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

// ParseFile reads a ZMK .keymap file and extracts layers and combos.
func ParseFile(path string) (*Keymap, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read keymap: %w", err)
    }
    return Parse(string(data))
}

// Parse parses ZMK keymap content.
func Parse(content string) (*Keymap, error) {
    km := &Keymap{}

    // Extract combos
    km.Combos = parseCombos(content)

    // Extract layers from the keymap node
    km.Layers = parseLayers(content)

    if len(km.Layers) == 0 {
        return nil, fmt.Errorf("no layers found in keymap")
    }

    return km, nil
}

var bindingsRe = regexp.MustCompile(`bindings\s*=\s*<([^>]+)>`)
var labelRe = regexp.MustCompile(`label\s*=\s*"([^"]+)"`)
var comboNameRe = regexp.MustCompile(`(\w+)\s*\{`)
var keyPosRe = regexp.MustCompile(`key-positions\s*=\s*<([^>]+)>`)

func parseLayers(content string) []Layer {
    // Find the keymap node
    kmStart := strings.Index(content, "keymap {")
    if kmStart == -1 {
        kmStart = strings.Index(content, "keymap{")
    }
    if kmStart == -1 {
        return nil
    }

    // Find the keymap block by brace matching
    kmBlock := extractBlock(content[kmStart:])

    // Split into layer blocks — find each "*_layer {" or block with label
    var layers []Layer
    idx := 0
    for idx < len(kmBlock) {
        // Find next "bindings = <"
        bindLoc := bindingsRe.FindStringIndex(kmBlock[idx:])
        if bindLoc == nil {
            break
        }

        // Look backwards from bindings for the label
        prefix := kmBlock[idx : idx+bindLoc[0]]
        labelMatch := labelRe.FindStringSubmatch(prefix)
        name := ""
        if labelMatch != nil {
            name = labelMatch[1]
        }

        // Extract bindings
        bindMatch := bindingsRe.FindStringSubmatch(kmBlock[idx:])
        if bindMatch == nil {
            break
        }

        keys := parseBindings(bindMatch[1])
        layers = append(layers, Layer{
            Name:     name,
            Index:    len(layers),
            Bindings: keys,
        })

        idx += bindLoc[1]
    }

    return layers
}

func parseBindings(raw string) []Key {
    // Split on & to get individual bindings
    // Raw looks like: "&kp Q  &kp W  &mt LGUI A ..."
    raw = strings.TrimSpace(raw)

    var keys []Key
    parts := strings.Split(raw, "&")

    for _, part := range parts {
        part = strings.TrimSpace(part)
        if part == "" {
            continue
        }

        fields := strings.Fields(part)
        if len(fields) == 0 {
            continue
        }

        k := Key{Position: len(keys)}
        behavior := fields[0]

        switch behavior {
        case "kp":
            k.Type = "kp"
            if len(fields) >= 2 {
                k.Tap = fields[1]
                k.Char = ZMKCodeToChar(fields[1])
            }
        case "mt":
            k.Type = "mt"
            if len(fields) >= 3 {
                k.Hold = fields[1]
                k.Tap = fields[2]
                k.Char = ZMKCodeToChar(fields[2])
            }
        case "lt":
            k.Type = "lt"
            if len(fields) >= 3 {
                k.Hold = fields[1]
                k.Tap = fields[2]
                // layer-tap: TAB and ESC aren't simple typeable chars in this context
                // but we still resolve for completeness
                k.Char = ZMKCodeToChar(fields[2])
            }
        case "trans":
            k.Type = "trans"
        case "none":
            k.Type = "none"
        case "mo":
            k.Type = "mo"
            if len(fields) >= 2 {
                k.Hold = fields[1]
            }
        case "sys_reset":
            k.Type = "sys"
            k.Tap = "RESET"
        case "bootloader":
            k.Type = "sys"
            k.Tap = "BOOT"
        case "bt":
            k.Type = "bt"
            if len(fields) >= 2 {
                k.Tap = strings.Join(fields[1:], " ")
            }
        case "out":
            k.Type = "out"
            if len(fields) >= 2 {
                k.Tap = fields[1]
            }
        default:
            k.Type = behavior
            if len(fields) >= 2 {
                k.Tap = strings.Join(fields[1:], " ")
            }
        }

        keys = append(keys, k)
    }

    return keys
}

func parseCombos(content string) []Combo {
    // Find the combos node
    cStart := strings.Index(content, "combos {")
    if cStart == -1 {
        cStart = strings.Index(content, "combos{")
    }
    if cStart == -1 {
        return nil
    }

    cBlock := extractBlock(content[cStart:])

    var combos []Combo

    // Find each combo sub-block (skip "compatible" property)
    lines := strings.Split(cBlock, "\n")
    var currentName string
    var currentPositions []int
    var currentBinding string
    inCombo := false
    depth := 0

    for _, line := range lines {
        trimmed := strings.TrimSpace(line)

        if strings.Contains(trimmed, "compatible") {
            continue
        }

        // Detect combo block start: "name {"
        if match := comboNameRe.FindStringSubmatch(trimmed); match != nil && !strings.Contains(trimmed, "combos") {
            if depth == 1 { // inside combos node
                currentName = match[1]
                inCombo = true
            }
        }

        for _, ch := range trimmed {
            if ch == '{' {
                depth++
            } else if ch == '}' {
                depth--
                if depth == 1 && inCombo {
                    // End of combo block
                    combos = append(combos, Combo{
                        Name:      currentName,
                        Positions: currentPositions,
                        Output:    currentBinding,
                    })
                    currentName = ""
                    currentPositions = nil
                    currentBinding = ""
                    inCombo = false
                }
            }
        }

        if inCombo {
            // Parse key-positions
            if posMatch := keyPosRe.FindStringSubmatch(trimmed); posMatch != nil {
                for _, s := range strings.Fields(posMatch[1]) {
                    if n, err := strconv.Atoi(s); err == nil {
                        currentPositions = append(currentPositions, n)
                    }
                }
            }
            // Parse bindings
            if bindMatch := bindingsRe.FindStringSubmatch(trimmed); bindMatch != nil {
                // Extract the character from the binding
                parts := strings.Fields(strings.TrimPrefix(strings.TrimSpace(bindMatch[1]), "&"))
                if len(parts) >= 2 && parts[0] == "kp" {
                    currentBinding = ZMKCodeToChar(parts[1])
                }
            }
        }
    }

    return combos
}

// extractBlock extracts content between the first { and its matching }.
func extractBlock(s string) string {
    start := strings.Index(s, "{")
    if start == -1 {
        return ""
    }

    depth := 0
    for i := start; i < len(s); i++ {
        switch s[i] {
        case '{':
            depth++
        case '}':
            depth--
            if depth == 0 {
                return s[start : i+1]
            }
        }
    }
    return ""
}
```

**Step 4: Run tests**

Run: `cd trainer && go test ./internal/keymap/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add trainer/internal/keymap/parser.go trainer/internal/keymap/parser_test.go
git commit -m "feat: add ZMK keymap parser with layer and combo extraction"
```

---

### Task 4: Lesson engine — stages and exercise generation

**Files:**
- Create: `trainer/internal/lesson/lesson.go`
- Create: `trainer/internal/lesson/lesson_test.go`
- Create: `trainer/data/wordlists/common.txt`

**Step 1: Create a basic wordlist**

Create `trainer/data/wordlists/common.txt` with the 200 most common English words (one per line, lowercase). These are short, common words good for typing practice.

**Step 2: Write the test**

```go
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

func containsBigram(s, bg string) bool {
    for i := 0; i < len(s)-1; i++ {
        if s[i:i+2] == bg {
            return true
        }
    }
    return false
}
```

**Step 3: Run test to verify it fails**

Run: `cd trainer && go test ./internal/lesson/ -v`
Expected: FAIL

**Step 4: Write lesson.go**

```go
package lesson

import (
    "math/rand"
    "strings"
)

// Stage defines a set of keys to practice.
type Stage struct {
    Name string
    Keys []string // typeable characters in this stage
}

func HomeRow() Stage {
    return Stage{Name: "home_row", Keys: []string{"a", "r", "s", "t", "d", "h", "n", "e", "i", "o"}}
}

func TopRow() Stage {
    return Stage{Name: "top_row", Keys: []string{"q", "w", "f", "p", "g", "j", "l", "u", "y", ";"}}
}

func BottomRow() Stage {
    return Stage{Name: "bottom_row", Keys: []string{"z", "x", "c", "v", "b", "k", "m", ",", ".", "/"}}
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
```

**Step 5: Run tests**

Run: `cd trainer && go test ./internal/lesson/ -v`
Expected: PASS

**Step 6: Commit**

```bash
git add trainer/internal/lesson/ trainer/data/
git commit -m "feat: add lesson engine with stages, word filtering, and exercise generation"
```

---

### Task 5: Stats tracking and persistence

**Files:**
- Create: `trainer/internal/stats/stats.go`
- Create: `trainer/internal/stats/stats_test.go`

**Step 1: Write the test**

```go
package stats

import (
    "os"
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
        TotalChars: 250,      // 250 characters
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
    path := filepath.Join(dir, "history.json")
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
    // Just verify ConfigDir returns something reasonable
    dir := ConfigDir()
    if dir == "" {
        t.Error("ConfigDir should not be empty")
    }
    if !filepath.IsAbs(dir) {
        t.Error("ConfigDir should return absolute path")
    }

    // Clean up — don't actually create the dir in home during tests
    _ = os.RemoveAll(dir)
}
```

**Step 2: Run test to verify it fails**

Run: `cd trainer && go test ./internal/stats/ -v`
Expected: FAIL

**Step 3: Write stats.go**

```go
package stats

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
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
    records, _ := h.Load() // ignore error if file doesn't exist yet
    records = append(records, record)
    return h.write(records)
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
    if err := json.Unmarshal(data, &records); err != nil {
        return nil, fmt.Errorf("parse history: %w", err)
    }
    return records, nil
}

func (h *History) write(records []SessionRecord) error {
    if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
        return err
    }
    data, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(h.path, data, 0o644)
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
```

**Step 4: Run tests**

Run: `cd trainer && go test ./internal/stats/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add trainer/internal/stats/
git commit -m "feat: add stats tracking with session recording and persistence"
```

---

### Task 6: TUI — Bubble Tea app skeleton with screen routing

**Files:**
- Create: `trainer/internal/ui/app.go`
- Create: `trainer/internal/ui/menu.go`
- Update: `trainer/cmd/totem-trainer/main.go`

**Step 1: Write app.go — top-level Bubble Tea model with screen routing**

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
    screenMenu screen = iota
    screenPicker
    screenTyping
    screenResults
    screenStats
)

// App is the top-level Bubble Tea model.
type App struct {
    current    screen
    menu       MenuModel
    picker     PickerModel
    typing     TypingModel
    results    ResultsModel
    statsView  StatsViewModel
    keymapPath string
    width      int
    height     int
}

func NewApp(keymapPath string) App {
    return App{
        current:    screenMenu,
        menu:       NewMenuModel(),
        keymapPath: keymapPath,
    }
}

func (a App) Init() tea.Cmd {
    return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        a.width = msg.Width
        a.height = msg.Height

    case tea.KeyMsg:
        if msg.String() == "ctrl+c" {
            return a, tea.Quit
        }
    }

    switch a.current {
    case screenMenu:
        return a.updateMenu(msg)
    case screenPicker:
        return a.updatePicker(msg)
    case screenTyping:
        return a.updateTyping(msg)
    case screenResults:
        return a.updateResults(msg)
    case screenStats:
        return a.updateStats(msg)
    }

    return a, nil
}

func (a App) View() string {
    switch a.current {
    case screenMenu:
        return a.menu.View()
    case screenPicker:
        return a.picker.View()
    case screenTyping:
        return a.typing.View(a.width)
    case screenResults:
        return a.results.View()
    case screenStats:
        return a.statsView.View()
    default:
        return ""
    }
}

// Screen transition messages
type switchScreen struct{ to screen }
type startLesson struct{ stageIdx int }
type lessonDone struct{}
type backToMenu struct{}
```

**Step 2: Write menu.go — main menu**

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        MarginBottom(1)

    menuItemStyle = lipgloss.NewStyle().
        PaddingLeft(2)

    selectedStyle = lipgloss.NewStyle().
        PaddingLeft(2).
        Foreground(lipgloss.Color("170")).
        Bold(true)

    dimStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))
)

type menuChoice int

const (
    menuPractice menuChoice = iota
    menuStats
    menuQuit
)

type MenuModel struct {
    cursor  int
    choices []string
}

func NewMenuModel() MenuModel {
    return MenuModel{
        choices: []string{"Practice", "Stats", "Quit"},
    }
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter":
            switch menuChoice(m.cursor) {
            case menuPractice:
                return m, func() tea.Msg { return switchScreen{to: screenPicker} }
            case menuStats:
                return m, func() tea.Msg { return switchScreen{to: screenStats} }
            case menuQuit:
                return m, tea.Quit
            }
        case "q":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m MenuModel) View() string {
    s := titleStyle.Render("totem-trainer") + "\n\n"

    for i, choice := range m.choices {
        if i == m.cursor {
            s += selectedStyle.Render("> "+choice) + "\n"
        } else {
            s += menuItemStyle.Render("  "+choice) + "\n"
        }
    }

    s += "\n" + dimStyle.Render("j/k to move, enter to select, q to quit")

    return s
}
```

**Step 3: Write screen transition handlers in app.go**

Add these methods to `app.go`:

```go
func (a App) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case switchScreen:
        a.current = msg.to
        if msg.to == screenPicker {
            a.picker = NewPickerModel()
        }
        if msg.to == screenStats {
            a.statsView = NewStatsViewModel()
        }
        return a, nil
    default:
        var cmd tea.Cmd
        a.menu, cmd = a.menu.Update(msg)
        return a, cmd
    }
}

func (a App) updatePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case backToMenu:
        a.current = screenMenu
        return a, nil
    case startLesson:
        a.current = screenTyping
        a.typing = NewTypingModel(msg.stageIdx, a.keymapPath)
        return a, nil
    default:
        var cmd tea.Cmd
        a.picker, cmd = a.picker.Update(msg)
        return a, cmd
    }
}

func (a App) updateTyping(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case lessonDone:
        a.current = screenResults
        a.results = NewResultsModel(a.typing.session)
        return a, nil
    default:
        var cmd tea.Cmd
        a.typing, cmd = a.typing.Update(msg)
        return a, cmd
    }
}

func (a App) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case backToMenu:
        a.current = screenMenu
        return a, nil
    case switchScreen:
        a.current = msg.(switchScreen).to
        if msg.(switchScreen).to == screenPicker {
            a.picker = NewPickerModel()
        }
        return a, nil
    default:
        var cmd tea.Cmd
        a.results, cmd = a.results.Update(msg)
        return a, cmd
    }
}

func (a App) updateStats(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case backToMenu:
        a.current = screenMenu
        return a, nil
    default:
        var cmd tea.Cmd
        a.statsView, cmd = a.statsView.Update(msg)
        return a, cmd
    }
}
```

**Step 4: Update main.go**

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/paul/totem-trainer/internal/ui"
)

func main() {
    // Find keymap relative to binary or use flag
    keymapPath := "config/totem.keymap"
    if len(os.Args) > 1 {
        keymapPath = os.Args[1]
    }

    // Resolve to absolute path
    if !filepath.IsAbs(keymapPath) {
        if abs, err := filepath.Abs(keymapPath); err == nil {
            keymapPath = abs
        }
    }

    p := tea.NewProgram(ui.NewApp(keymapPath), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**Step 5: Create stub files for remaining screen models** (so it compiles)

Create `trainer/internal/ui/picker.go`, `trainer/internal/ui/typing.go`, `trainer/internal/ui/results.go`, `trainer/internal/ui/statsview.go` with minimal stubs that satisfy the interface.

**Step 6: Verify it builds**

Run: `cd trainer && go build ./cmd/totem-trainer`
Expected: compiles without error

**Step 7: Commit**

```bash
git add trainer/
git commit -m "feat: add Bubble Tea TUI skeleton with menu and screen routing"
```

---

### Task 7: TUI — Lesson picker screen

**Files:**
- Create: `trainer/internal/ui/picker.go`

**Step 1: Write picker.go**

```go
package ui

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/paul/totem-trainer/internal/lesson"
)

type PickerModel struct {
    cursor int
    stages []lesson.Stage
}

func NewPickerModel() PickerModel {
    return PickerModel{stages: lesson.AllStages()}
}

func (m PickerModel) Update(msg tea.Msg) (PickerModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.stages)-1 {
                m.cursor++
            }
        case "enter":
            idx := m.cursor
            return m, func() tea.Msg { return startLesson{stageIdx: idx} }
        case "esc", "q":
            return m, func() tea.Msg { return backToMenu{} }
        }
    }
    return m, nil
}

func (m PickerModel) View() string {
    s := titleStyle.Render("Choose a lesson") + "\n\n"

    for i, stage := range m.stages {
        label := fmt.Sprintf("%-12s  %s", stage.Name, formatKeys(stage.Keys))
        if i == m.cursor {
            s += selectedStyle.Render("> "+label) + "\n"
        } else {
            s += menuItemStyle.Render("  "+label) + "\n"
        }
    }

    s += "\n" + dimStyle.Render("j/k to move, enter to start, esc to go back")
    return s
}

func formatKeys(keys []string) string {
    if len(keys) > 15 {
        return fmt.Sprintf("[%d keys]", len(keys))
    }
    s := ""
    for i, k := range keys {
        if i > 0 {
            s += " "
        }
        s += k
    }
    return s
}
```

**Step 2: Verify it builds**

Run: `cd trainer && go build ./cmd/totem-trainer`
Expected: compiles

**Step 3: Commit**

```bash
git add trainer/internal/ui/picker.go
git commit -m "feat: add lesson picker screen"
```

---

### Task 8: TUI — Typing session screen (core experience)

**Files:**
- Create: `trainer/internal/ui/typing.go`
- Modify: `trainer/internal/ui/app.go` (if needed for message types)

**Step 1: Write typing.go**

```go
package ui

import (
    "fmt"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/paul/totem-trainer/internal/lesson"
    "github.com/paul/totem-trainer/internal/stats"
)

var (
    correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
    wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
    cursorStyle  = lipgloss.NewStyle().Background(lipgloss.Color("240"))
    upcomingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
    statsBarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).MarginTop(1)
)

type tickMsg time.Time

type TypingModel struct {
    exercise   string
    input      []rune
    cursor     int
    startTime  time.Time
    started    bool
    session    *stats.Session
    stage      lesson.Stage
    lastKeyAt  time.Time
}

func NewTypingModel(stageIdx int, keymapPath string) TypingModel {
    stages := lesson.AllStages()
    stage := stages[stageIdx]

    // Build allowed keys set
    allowed := make(map[string]bool)
    for _, k := range stage.Keys {
        allowed[k] = true
    }

    // Load wordlist — use built-in common words for now
    words := lesson.CommonWords()
    filtered := lesson.FilterWords(words, allowed)

    var exercise string
    if stage.Name == "mixed" {
        // Use terminal snippets for mixed mode
        snippets := lesson.TerminalSnippets()
        exercise = lesson.GenerateExercise(snippets, 5)
    } else if len(filtered) > 0 {
        exercise = lesson.GenerateExercise(filtered, 15)
    } else {
        // Fallback: just repeat the keys
        exercise = lesson.GenerateExercise(stage.Keys, 20)
    }

    return TypingModel{
        exercise: exercise,
        input:    make([]rune, 0),
        session:  stats.NewSession(),
        stage:    stage,
    }
}

func (m TypingModel) Init() tea.Cmd {
    return nil
}

func (m TypingModel) Update(msg tea.Msg) (TypingModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            return m, func() tea.Msg { return backToMenu{} }
        case "ctrl+c":
            return m, tea.Quit
        case "backspace":
            if m.cursor > 0 {
                m.cursor--
                if len(m.input) > 0 {
                    m.input = m.input[:len(m.input)-1]
                }
            }
            return m, nil
        }

        // Handle regular character input
        if len(msg.Runes) == 0 {
            return m, nil
        }

        ch := msg.Runes[0]

        if !m.started {
            m.started = true
            m.startTime = time.Now()
            m.lastKeyAt = time.Now()
        }

        now := time.Now()
        latency := now.Sub(m.lastKeyAt)
        m.lastKeyAt = now

        // Record the keystroke
        expected := rune(0)
        if m.cursor < len(m.exercise) {
            expected = rune(m.exercise[m.cursor])
        }

        m.session.TotalChars++
        key := string(expected)
        if _, ok := m.session.PerKey[key]; !ok {
            m.session.PerKey[key] = &stats.KeyStats{}
        }

        if ch == expected {
            m.session.Correct++
            m.session.PerKey[key].RecordHit(latency)
        } else {
            m.session.PerKey[key].RecordMiss(latency)
        }

        m.input = append(m.input, ch)
        m.cursor++

        // Check if exercise is complete
        if m.cursor >= len(m.exercise) {
            m.session.Duration = time.Since(m.startTime)
            return m, func() tea.Msg { return lessonDone{} }
        }
    }

    return m, nil
}

func (m TypingModel) View(width int) string {
    if width == 0 {
        width = 80
    }

    var b strings.Builder

    b.WriteString(titleStyle.Render(m.stage.Name) + "\n\n")

    // Render the exercise with colored characters
    for i, ch := range m.exercise {
        s := string(ch)
        if i < m.cursor {
            // Already typed
            if i < len(m.input) && m.input[i] == ch {
                b.WriteString(correctStyle.Render(s))
            } else {
                b.WriteString(wrongStyle.Render(s))
            }
        } else if i == m.cursor {
            // Current position
            b.WriteString(cursorStyle.Render(s))
        } else {
            // Upcoming
            b.WriteString(upcomingStyle.Render(s))
        }
    }

    b.WriteString("\n")

    // Stats bar
    if m.started {
        elapsed := time.Since(m.startTime)
        minutes := elapsed.Minutes()
        wpm := 0.0
        if minutes > 0 {
            wpm = (float64(m.session.TotalChars) / 5.0) / minutes
        }
        acc := m.session.Accuracy() * 100

        bar := fmt.Sprintf("WPM: %.0f  |  Accuracy: %.0f%%  |  %d/%d",
            wpm, acc, m.cursor, len(m.exercise))
        b.WriteString(statsBarStyle.Render(bar))
    }

    b.WriteString("\n\n" + dimStyle.Render("esc to quit"))

    return b.String()
}
```

**Step 2: Add CommonWords function to lesson package**

Add to `trainer/internal/lesson/lesson.go`:

```go
// CommonWords returns a built-in list of common English words.
func CommonWords() []string {
    return []string{
        "the", "and", "that", "have", "for", "not", "with", "you",
        "this", "but", "his", "from", "they", "she", "what", "their",
        "will", "each", "make", "like", "long", "thing", "see", "him",
        "two", "has", "her", "there", "one", "our", "out", "are",
        "other", "were", "all", "time", "when", "use", "your", "how",
        "said", "some", "than", "them", "would", "into", "then", "its",
        "over", "also", "did", "down", "only", "way", "find", "get",
        "come", "made", "after", "most", "just", "know", "take", "some",
        "could", "good", "state", "year", "much", "new", "been", "now",
        "old", "great", "where", "still", "hand", "high", "here", "end",
        "does", "done", "last", "long", "too", "own", "same", "tell",
        "set", "need", "home", "head", "stand", "start", "might", "show",
        "hard", "far", "run", "help", "turn", "move", "live", "real",
        "left", "late", "line", "read", "seem", "side", "went", "world",
        "data", "test", "code", "file", "list", "name", "sort", "send",
        "note", "tree", "node", "edit", "done", "mode", "kind", "next",
    }
}
```

**Step 3: Verify it builds**

Run: `cd trainer && go build ./cmd/totem-trainer`
Expected: compiles

**Step 4: Run it manually**

Run: `cd trainer && ./totem-trainer ../config/totem.keymap`
Expected: shows menu, can navigate to picker, start a lesson, type characters

**Step 5: Commit**

```bash
git add trainer/
git commit -m "feat: add typing session screen with real-time feedback"
```

---

### Task 9: TUI — Results screen

**Files:**
- Create: `trainer/internal/ui/results.go`

**Step 1: Write results.go**

```go
package ui

import (
    "fmt"
    "sort"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/paul/totem-trainer/internal/stats"
)

var (
    goodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
    okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
    badStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

type ResultsModel struct {
    session *stats.Session
}

func NewResultsModel(session *stats.Session) ResultsModel {
    return ResultsModel{session: session}
}

func (m ResultsModel) Update(msg tea.Msg) (ResultsModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter", "p":
            return m, func() tea.Msg { return switchScreen{to: screenPicker} }
        case "q", "esc":
            return m, func() tea.Msg { return backToMenu{} }
        }
    }
    return m, nil
}

func (m ResultsModel) View() string {
    s := m.session
    var b strings.Builder

    b.WriteString(titleStyle.Render("Results") + "\n\n")

    wpm := s.WPM()
    acc := s.Accuracy() * 100

    // WPM with color
    wpmStr := fmt.Sprintf("%.0f WPM", wpm)
    if wpm >= 40 {
        b.WriteString("  Speed:    " + goodStyle.Render(wpmStr) + "\n")
    } else if wpm >= 20 {
        b.WriteString("  Speed:    " + okStyle.Render(wpmStr) + "\n")
    } else {
        b.WriteString("  Speed:    " + badStyle.Render(wpmStr) + "\n")
    }

    // Accuracy with color
    accStr := fmt.Sprintf("%.0f%%", acc)
    if acc >= 95 {
        b.WriteString("  Accuracy: " + goodStyle.Render(accStr) + "\n")
    } else if acc >= 85 {
        b.WriteString("  Accuracy: " + okStyle.Render(accStr) + "\n")
    } else {
        b.WriteString("  Accuracy: " + badStyle.Render(accStr) + "\n")
    }

    b.WriteString(fmt.Sprintf("  Duration: %.0fs\n", s.Duration.Seconds()))
    b.WriteString(fmt.Sprintf("  Keys:     %d\n", s.TotalChars))

    // Per-key breakdown — show worst keys
    b.WriteString("\n" + titleStyle.Render("Weakest Keys") + "\n\n")

    type keyAcc struct {
        key string
        acc float64
    }
    var sorted []keyAcc
    for k, ks := range s.PerKey {
        if k == " " {
            k = "space"
        }
        sorted = append(sorted, keyAcc{key: k, acc: ks.Accuracy()})
    }
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].acc < sorted[j].acc
    })

    shown := 0
    for _, ka := range sorted {
        if shown >= 5 {
            break
        }
        accPct := ka.acc * 100
        style := goodStyle
        if accPct < 85 {
            style = badStyle
        } else if accPct < 95 {
            style = okStyle
        }
        b.WriteString(fmt.Sprintf("  %-8s %s\n", ka.key, style.Render(fmt.Sprintf("%.0f%%", accPct))))
        shown++
    }

    b.WriteString("\n" + dimStyle.Render("p to practice again, esc to menu"))

    return b.String()
}
```

**Step 2: Verify it builds**

Run: `cd trainer && go build ./cmd/totem-trainer`
Expected: compiles

**Step 3: Commit**

```bash
git add trainer/internal/ui/results.go
git commit -m "feat: add results screen with per-key accuracy breakdown"
```

---

### Task 10: TUI — Stats view screen

**Files:**
- Create: `trainer/internal/ui/statsview.go`

**Step 1: Write statsview.go**

```go
package ui

import (
    "fmt"
    "path/filepath"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/paul/totem-trainer/internal/stats"
)

type StatsViewModel struct {
    records []stats.SessionRecord
    err     error
}

func NewStatsViewModel() StatsViewModel {
    dir := stats.ConfigDir()
    h := stats.NewHistory(filepath.Join(dir, "history.json"))
    records, err := h.Load()
    return StatsViewModel{records: records, err: err}
}

func (m StatsViewModel) Update(msg tea.Msg) (StatsViewModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "esc":
            return m, func() tea.Msg { return backToMenu{} }
        }
    }
    return m, nil
}

func (m StatsViewModel) View() string {
    var b strings.Builder

    b.WriteString(titleStyle.Render("Stats") + "\n\n")

    if m.err != nil {
        b.WriteString("  Error loading stats: " + m.err.Error() + "\n")
        b.WriteString("\n" + dimStyle.Render("esc to go back"))
        return b.String()
    }

    if len(m.records) == 0 {
        b.WriteString("  No sessions recorded yet. Go practice!\n")
        b.WriteString("\n" + dimStyle.Render("esc to go back"))
        return b.String()
    }

    // Recent sessions
    b.WriteString("  Recent sessions:\n\n")
    start := len(m.records) - 10
    if start < 0 {
        start = 0
    }
    for i := len(m.records) - 1; i >= start; i-- {
        r := m.records[i]
        wpmStyle := goodStyle
        if r.WPM < 20 {
            wpmStyle = badStyle
        } else if r.WPM < 40 {
            wpmStyle = okStyle
        }
        accStyle := goodStyle
        if r.Accuracy < 0.85 {
            accStyle = badStyle
        } else if r.Accuracy < 0.95 {
            accStyle = okStyle
        }

        ts := r.Timestamp.Format("Jan 02 15:04")
        b.WriteString(fmt.Sprintf("  %s  %-12s  %s  %s\n",
            dimStyle.Render(ts),
            r.Stage,
            wpmStyle.Render(fmt.Sprintf("%4.0f WPM", r.WPM)),
            accStyle.Render(fmt.Sprintf("%3.0f%%", r.Accuracy*100)),
        ))
    }

    // Averages
    if len(m.records) >= 2 {
        b.WriteString("\n  Averages (last 10):\n")
        count := 10
        if count > len(m.records) {
            count = len(m.records)
        }
        var totalWPM, totalAcc float64
        for i := len(m.records) - count; i < len(m.records); i++ {
            totalWPM += m.records[i].WPM
            totalAcc += m.records[i].Accuracy
        }
        avgWPM := totalWPM / float64(count)
        avgAcc := totalAcc / float64(count) * 100
        b.WriteString(fmt.Sprintf("  WPM: %.0f  Accuracy: %.0f%%\n", avgWPM, avgAcc))
    }

    b.WriteString("\n" + dimStyle.Render("esc to go back"))
    return b.String()
}
```

**Step 2: Wire up session saving — after results, persist to disk**

Add to `app.go` updateTyping:

```go
// In the lessonDone case, after creating results, save the session
func (a App) saveSession() {
    dir := stats.ConfigDir()
    h := stats.NewHistory(filepath.Join(dir, "history.json"))

    record := stats.SessionRecord{
        Timestamp: time.Now(),
        Stage:     a.typing.stage.Name,
        Duration:  a.typing.session.Duration.Seconds(),
        WPM:       a.typing.session.WPM(),
        Accuracy:  a.typing.session.Accuracy(),
        TotalKeys: a.typing.session.TotalChars,
        PerKey:    make(map[string]stats.KeyRecord),
    }
    for k, ks := range a.typing.session.PerKey {
        record.PerKey[k] = stats.KeyRecord{
            Hits:   ks.Hits,
            Misses: ks.Misses,
            AvgMs:  ks.AvgLatency(),
        }
    }

    _ = h.Append(record) // best-effort save
}
```

**Step 3: Verify it builds and run end-to-end**

Run: `cd trainer && go build ./cmd/totem-trainer && ./totem-trainer ../config/totem.keymap`
Expected: full flow — menu → pick lesson → type → see results → stats shows history

**Step 4: Commit**

```bash
git add trainer/
git commit -m "feat: add stats view and session persistence"
```

---

### Task 11: Create common wordlist file and embed it

**Files:**
- Create: `trainer/data/wordlists/common.txt`
- Modify: `trainer/internal/lesson/lesson.go` — use `embed` to load wordlist

**Step 1: Create common.txt**

Curated list of ~200 common short English words, one per line. Focus on words typeable with common keys and useful for practice.

**Step 2: Update lesson.go to use go:embed**

```go
import "embed"

//go:embed ../../data/wordlists/common.txt
var commonWordsFile string

func CommonWords() []string {
    var words []string
    for _, line := range strings.Split(commonWordsFile, "\n") {
        w := strings.TrimSpace(line)
        if w != "" {
            words = append(words, w)
        }
    }
    return words
}
```

**Step 3: Verify it builds**

Run: `cd trainer && go build ./cmd/totem-trainer`
Expected: compiles with embedded wordlist

**Step 4: Commit**

```bash
git add trainer/data/ trainer/internal/lesson/
git commit -m "feat: embed common wordlist using go:embed"
```

---

### Task 12: Integration test — full parse-to-lesson pipeline

**Files:**
- Create: `trainer/internal/lesson/integration_test.go`

**Step 1: Write integration test**

```go
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

    // Verify NAV layer has numbers
    nav := km.Layers[1]
    var nums []string
    for _, k := range nav.Bindings {
        if k.Char >= "0" && k.Char <= "9" {
            nums = append(nums, k.Char)
        }
    }
    if len(nums) != 10 {
        t.Errorf("NAV layer should have 10 numbers, got %d: %v", len(nums), nums)
    }
}
```

**Step 2: Run the test**

Run: `cd trainer && go test ./internal/lesson/ -v -run Integration`
Expected: PASS

**Step 3: Commit**

```bash
git add trainer/internal/lesson/integration_test.go
git commit -m "test: add integration test for keymap-to-lesson pipeline"
```

---

### Task 13: Polish and final verification

**Step 1: Run all tests**

Run: `cd trainer && go test ./... -v`
Expected: all PASS

**Step 2: Run the app end-to-end**

Run: `cd trainer && go build ./cmd/totem-trainer && ./totem-trainer ../config/totem.keymap`
Manual test checklist:
- [ ] Menu shows, j/k navigation works
- [ ] Lesson picker shows all 7 stages
- [ ] Home row lesson generates words with only a,r,s,t,d,h,n,e,i,o
- [ ] Typing shows green/red feedback
- [ ] WPM and accuracy update live
- [ ] Results screen shows per-key breakdown
- [ ] Stats screen shows session history after completing a lesson
- [ ] ESC and q navigation works throughout

**Step 3: Final commit**

```bash
git add -A trainer/
git commit -m "feat: totem-trainer v0.1 — typing trainer for Totem keyboard"
```
