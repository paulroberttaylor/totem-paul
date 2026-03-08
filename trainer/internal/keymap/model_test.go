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
