package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	keyStyle = lipgloss.NewStyle().
			Width(5).
			Height(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	blankKeyStyle = keyStyle.Foreground(lipgloss.Color("245"))

	highlightKeyStyle = lipgloss.NewStyle().
				Width(5).
				Height(1).
				Align(lipgloss.Center).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("10")).
				Foreground(lipgloss.Color("10")).
				Bold(true)

	nextKeyStyle = lipgloss.NewStyle().
			Width(5).
			Height(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("75")).
			Foreground(lipgloss.Color("75")).
			Bold(true)

	thumbKeyStyle = lipgloss.NewStyle().
			Width(5).
			Height(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99"))

	thumbHighlightStyle = lipgloss.NewStyle().
				Width(5).
				Height(1).
				Align(lipgloss.Center).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("10")).
				Foreground(lipgloss.Color("10")).
				Bold(true)

	thumbNextStyle = lipgloss.NewStyle().
			Width(5).
			Height(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("75")).
			Foreground(lipgloss.Color("75")).
			Bold(true)

	// Invisible placeholder to fill the empty corners
	emptyCell = lipgloss.NewStyle().
			Width(7). // key width (5) + border (2)
			Height(3).
			Render("")

	keyboardTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			MarginBottom(0)
)

// TotemLayout holds the key labels for rendering the Totem keyboard.
// Empty strings render as blank keys.
type TotemLayout struct {
	Title     string
	Top       [10]string // 5 left + 5 right
	Home      [12]string // 1 pinky + 5 left + 5 right + 1 pinky
	Bottom    [10]string // 5 left + 5 right
	Thumbs    [6]string  // 3 left + 3 right
	Highlight string     // key label to highlight (last pressed)
	Next      string     // key label to highlight as next target
}

// charToKeyLabel maps a typed character to its label on the keyboard.
var charToKeyLabel = map[rune]string{
	' ': "SPC", '\t': "TAB", '\n': "ENT", '\b': "BSP",
	'\'': "'", '"': "\"", '\\': "\\",
}

// BlankTotem returns a layout with all blank keys.
func BlankTotem() TotemLayout {
	return TotemLayout{
		Title: "TOTEM",
	}
}

// BaseLayer returns the Colemak DHm base layer layout.
func BaseLayer() TotemLayout {
	return TotemLayout{
		Title:  "BASE",
		Top:    [10]string{"q", "w", "f", "p", "b", "j", "l", "u", "y", ";"},
		Home:   [12]string{"ESC", "a", "r", "s", "t", "g", "m", "n", "e", "i", "o", "TMX"},
		Bottom: [10]string{"z", "x", "c", "d", "v", "k", "h", ",", ".", "/"},
		Thumbs: [6]string{"DEL", "TAB", "SPC", "BSP", "ENT", "DEL"},
	}
}

// SymLayer returns the SYM layer layout.
func SymLayer() TotemLayout {
	return TotemLayout{
		Title:  "SYM (hold BSPC)",
		Top:    [10]string{"!", "@", "#", "$", "%", "^", "&", "*", "'", "\""},
		Home:   [12]string{"", "~", "`", "_", "|", "{", "}", "SFT", "GUI", "CTL", "ALT", ""},
		Bottom: [10]string{"\\", "(", ")", "[", "]", "-", "+", "=", "<", ">"},
		Thumbs: [6]string{"", "", "", "", "███", ""},
	}
}

// NumLayer returns the NUM layer layout.
func NumLayer() TotemLayout {
	return TotemLayout{
		Title:  "NUM (hold DEL)",
		Top:    [10]string{"", "", "", "", "", "`", "7", "8", "9", "~"},
		Home:   [12]string{"", "ALT", "CTL", "GUI", "SFT", "", "-", "4", "5", "6", "0", ""},
		Bottom: [10]string{"", "", "", "", "", "/", "1", "2", "3", "."},
		Thumbs: [6]string{"███", "", "", "=", "+", ""},
	}
}

// LayerForStage returns the appropriate keyboard layout for a given stage name.
func LayerForStage(stage string) TotemLayout {
	switch stage {
	case "symbols":
		return SymLayer()
	case "numbers":
		return NumLayer()
	default:
		return BaseLayer()
	}
}

// CharToLabel converts a typed rune to its keyboard label.
func CharToLabel(r rune) string {
	if label, ok := charToKeyLabel[r]; ok {
		return label
	}
	return strings.ToLower(string(r))
}

func pickStyle(label, highlight, next string, normal, hlStyle, nxStyle lipgloss.Style) lipgloss.Style {
	lower := strings.ToLower(label)
	if highlight != "" && lower == strings.ToLower(highlight) {
		return hlStyle
	}
	if next != "" && lower == strings.ToLower(next) {
		return nxStyle
	}
	return normal
}

func renderKey(label, highlight, next string, normal, hlStyle, nxStyle lipgloss.Style) string {
	if label == "" {
		label = " "
	}
	style := pickStyle(label, highlight, next, normal, hlStyle, nxStyle)
	return style.Render(label)
}

// RenderKeyboard draws the Totem split keyboard layout.
func RenderKeyboard(layout TotemLayout) string {
	var rows []string
	hl := layout.Highlight
	nx := layout.Next

	if layout.Title != "" {
		rows = append(rows, keyboardTitle.Render(layout.Title))
	}

	// Top row: 5 left + gap + 5 right
	leftTop := make([]string, 5)
	rightTop := make([]string, 5)
	for i := 0; i < 5; i++ {
		leftTop[i] = renderKey(layout.Top[i], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle)
		rightTop[i] = renderKey(layout.Top[5+i], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle)
	}
	topRow := lipgloss.JoinHorizontal(lipgloss.Center,
		emptyCell,
		lipgloss.JoinHorizontal(lipgloss.Center, leftTop...),
		"   ",
		lipgloss.JoinHorizontal(lipgloss.Center, rightTop...),
		emptyCell,
	)
	rows = append(rows, topRow)

	// Home row: pinky + 5 left + gap + 5 right + pinky
	leftHome := make([]string, 5)
	rightHome := make([]string, 5)
	for i := 0; i < 5; i++ {
		leftHome[i] = renderKey(layout.Home[1+i], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle)
		rightHome[i] = renderKey(layout.Home[6+i], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle)
	}
	homeRow := lipgloss.JoinHorizontal(lipgloss.Center,
		renderKey(layout.Home[0], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle),
		lipgloss.JoinHorizontal(lipgloss.Center, leftHome...),
		"   ",
		lipgloss.JoinHorizontal(lipgloss.Center, rightHome...),
		renderKey(layout.Home[11], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle),
	)
	rows = append(rows, homeRow)

	// Bottom row: 5 left + gap + 5 right
	leftBot := make([]string, 5)
	rightBot := make([]string, 5)
	for i := 0; i < 5; i++ {
		leftBot[i] = renderKey(layout.Bottom[i], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle)
		rightBot[i] = renderKey(layout.Bottom[5+i], hl, nx, blankKeyStyle, highlightKeyStyle, nextKeyStyle)
	}
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Center,
		emptyCell,
		lipgloss.JoinHorizontal(lipgloss.Center, leftBot...),
		"   ",
		lipgloss.JoinHorizontal(lipgloss.Center, rightBot...),
		emptyCell,
	)
	rows = append(rows, bottomRow)

	// Thumb row: 3 left + gap + 3 right (centered)
	leftThumbs := make([]string, 3)
	rightThumbs := make([]string, 3)
	for i := 0; i < 3; i++ {
		leftThumbs[i] = renderKey(layout.Thumbs[i], hl, nx, thumbKeyStyle, thumbHighlightStyle, thumbNextStyle)
		rightThumbs[i] = renderKey(layout.Thumbs[3+i], hl, nx, thumbKeyStyle, thumbHighlightStyle, thumbNextStyle)
	}
	thumbRow := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, leftThumbs...),
		"   ",
		lipgloss.JoinHorizontal(lipgloss.Center, rightThumbs...),
	)
	// Center thumbs under the main rows
	thumbRowCentered := lipgloss.NewStyle().
		Width(lipgloss.Width(topRow)).
		Align(lipgloss.Center).
		Render(thumbRow)
	rows = append(rows, thumbRowCentered)

	return strings.Join(rows, "\n")
}
