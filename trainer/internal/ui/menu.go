package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	itemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

type menuChoice int

const (
	menuPractice menuChoice = iota
	menuStats
	menuQuit
)

type menuModel struct {
	cursor  int
	choices []string
}

func newMenuModel() menuModel {
	return menuModel{
		choices: []string{"Practice", "Stats", "Quit"},
	}
}

func (m menuModel) update(msg tea.Msg) (menuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = len(m.choices) - 1
			}
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "enter":
			switch menuChoice(m.cursor) {
			case menuPractice:
				return m, switchScreen(screenPicker)
			case menuStats:
				return m, func() tea.Msg {
					return switchScreenMsg{screen: screenStats}
				}
			case menuQuit:
				return m, tea.Quit
			}
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) view(width, height int) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Totem Trainer"))
	b.WriteString("\n\n")

	for i, choice := range m.choices {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("> %s", choice)))
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("  %s", choice)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("j/k: navigate • enter: select • q: quit"))

	content := b.String()
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
