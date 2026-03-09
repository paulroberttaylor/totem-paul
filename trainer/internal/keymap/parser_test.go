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
	want := []string{"BASE", "NAV", "NUM", "SYM"}
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
		pos  int
		typ  string
		tap  string
		hold string
		char string
	}{
		{0, "kp", "Q", "", "q"},           // Q
		{4, "kp", "B", "", "b"},           // B (Colemak DHm)
		{9, "kp", "SEMI", "", ";"},        // ;
		{10, "mt", "A", "LALT", "a"},      // A/ALT
		{12, "mt", "S", "LGUI", "s"},      // S/GUI
		{13, "mt", "T", "LSHFT", "t"},     // T/SHFT
		{20, "kp", "ESC", "", ""},         // ESC (not typeable)
		{24, "kp", "D", "", "d"},          // D (Colemak DHm — bottom row)
		{27, "kp", "H", "", "h"},          // H (Colemak DHm — bottom row)
		{33, "lt", "TAB", "NAV", ""},      // TAB/NAV (layer tap)
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

func TestParseKeymap_NUMLayerNumbers(t *testing.T) {
	km, err := ParseFile(testKeymapPath(t))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	num := km.Layers[2] // NUM is now layer 2

	// Position 6 should be 7 (N7)
	if num.Bindings[6].Char != "7" {
		t.Errorf("NUM pos 6: want char '7', got %q", num.Bindings[6].Char)
	}
	// Position 19 should be 0 (N0) — right pinky home row
	if num.Bindings[19].Char != "0" {
		t.Errorf("NUM pos 19: want char '0', got %q", num.Bindings[19].Char)
	}
}

func TestParseKeymap_SYMLayerSymbols(t *testing.T) {
	km, err := ParseFile(testKeymapPath(t))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	sym := km.Layers[3] // SYM is now layer 3

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
	if len(km.Combos) != 4 {
		t.Fatalf("want 4 combos, got %d", len(km.Combos))
	}

	// Verify combo names and outputs
	wantCombos := []struct {
		name   string
		output string
	}{
		{"combo_tilde", "~"},
		{"combo_pipe", "|"},
		{"combo_backtick", "`"},
		{"combo_underscore", "_"},
	}
	for i, want := range wantCombos {
		if km.Combos[i].Name != want.name {
			t.Errorf("combo %d: want name %q, got %q", i, want.name, km.Combos[i].Name)
		}
		if km.Combos[i].Output != want.output {
			t.Errorf("combo %d: want output %q, got %q", i, want.output, km.Combos[i].Output)
		}
	}
}
