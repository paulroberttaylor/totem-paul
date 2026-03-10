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

	thumbKeyStyle = lipgloss.NewStyle().
			Width(5).
			Height(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99"))

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
	Title  string
	Top    [10]string // 5 left + 5 right
	Home   [12]string // 1 pinky + 5 left + 5 right + 1 pinky
	Bottom [10]string // 5 left + 5 right
	Thumbs [6]string  // 3 left + 3 right
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

func renderKey(label string, style lipgloss.Style) string {
	if label == "" {
		label = " "
	}
	return style.Render(label)
}

// RenderKeyboard draws the Totem split keyboard layout.
func RenderKeyboard(layout TotemLayout) string {
	var rows []string

	if layout.Title != "" {
		rows = append(rows, keyboardTitle.Render(layout.Title))
	}

	// Top row: 5 left + gap + 5 right
	leftTop := make([]string, 5)
	rightTop := make([]string, 5)
	for i := 0; i < 5; i++ {
		leftTop[i] = renderKey(layout.Top[i], blankKeyStyle)
		rightTop[i] = renderKey(layout.Top[5+i], blankKeyStyle)
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
		leftHome[i] = renderKey(layout.Home[1+i], blankKeyStyle)
		rightHome[i] = renderKey(layout.Home[6+i], blankKeyStyle)
	}
	homeRow := lipgloss.JoinHorizontal(lipgloss.Center,
		renderKey(layout.Home[0], blankKeyStyle),
		lipgloss.JoinHorizontal(lipgloss.Center, leftHome...),
		"   ",
		lipgloss.JoinHorizontal(lipgloss.Center, rightHome...),
		renderKey(layout.Home[11], blankKeyStyle),
	)
	rows = append(rows, homeRow)

	// Bottom row: 5 left + gap + 5 right
	leftBot := make([]string, 5)
	rightBot := make([]string, 5)
	for i := 0; i < 5; i++ {
		leftBot[i] = renderKey(layout.Bottom[i], blankKeyStyle)
		rightBot[i] = renderKey(layout.Bottom[5+i], blankKeyStyle)
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
		leftThumbs[i] = renderKey(layout.Thumbs[i], thumbKeyStyle)
		rightThumbs[i] = renderKey(layout.Thumbs[3+i], thumbKeyStyle)
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
