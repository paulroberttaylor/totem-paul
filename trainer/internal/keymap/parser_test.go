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
		pos  int
		typ  string
		tap  string
		hold string
		char string
	}{
		{0, "kp", "Q", "", "q"},          // Q
		{9, "kp", "SEMI", "", ";"},        // ;
		{10, "mt", "A", "LGUI", "a"},      // A/GUI
		{13, "mt", "T", "LSHFT", "t"},     // T/SHFT
		{20, "kp", "ESC", "", ""},          // ESC (not typeable)
		{33, "lt", "TAB", "NAV", ""},       // TAB/NAV (layer tap)
		{34, "kp", "SPACE", "", " "},       // SPACE
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
