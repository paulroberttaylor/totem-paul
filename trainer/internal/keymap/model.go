package keymap

// Key represents a single key binding on a layer.
type Key struct {
	Position int    // 0-37 physical position
	Type     string // "kp", "mt", "lt", "trans", "none", "sys", "bt", "out", "mo"
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

// NewTotemLayout creates a TotemLayout with the Totem 38-key physical layout.
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

// Position returns layout metadata for a physical key position (0-37).
func (l *TotemLayout) Position(pos int) PositionInfo {
	if pos < 0 || pos >= 38 {
		return PositionInfo{}
	}
	return l.positions[pos]
}

// ZMK code to typeable character mapping.
var zmkCharMap = map[string]string{
	// Letters (single uppercase letter -> lowercase char)
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
