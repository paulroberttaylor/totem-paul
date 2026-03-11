package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paul/totem-trainer/internal/lesson"
)

var (
	pickerTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	stageNameStyle    = lipgloss.NewStyle().Bold(true)
	stageKeysStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	pickerHelpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	packHeaderStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214")).MarginTop(1)
)

type pickerModel struct {
	cursor int
	stages []lesson.Stage
}

func newPickerModel(stages []lesson.Stage) pickerModel {
	return pickerModel{
		stages: stages,
	}
}

func (m pickerModel) update(msg tea.Msg) (pickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "n", "down":
			m.cursor++
			if m.cursor >= len(m.stages) {
				m.cursor = len(m.stages) - 1
			}
		case "k", "e", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "enter":
			stage := m.stages[m.cursor]
			return m, startLesson(stage)
		case "esc", "q":
			return m, backToMenu()
		}
	}
	return m, nil
}

func (m pickerModel) view(width, height int) string {
	var b strings.Builder
	b.WriteString(pickerTitleStyle.Render("Select a Lesson"))
	b.WriteString("\n\n")

	lastPack := ""
	for i, stage := range m.stages {
		// Show pack group header when pack changes
		if stage.Pack != lastPack && stage.Pack != "" {
			b.WriteString(packHeaderStyle.Render(fmt.Sprintf("── %s ──", stage.Pack)))
			b.WriteString("\n")
		}
		lastPack = stage.Pack

		keyPreview := strings.Join(stage.Keys, " ")
		if len(keyPreview) > 40 {
			keyPreview = keyPreview[:37] + "..."
		}

		name := stageNameStyle.Render(stage.Name)
		keys := stageKeysStyle.Render(fmt.Sprintf("[%s]", keyPreview))

		if i == m.cursor {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("> %s %s", name, keys)))
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("  %s %s", name, keys)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(pickerHelpStyle.Render("n/e: navigate • enter: start • esc: back"))

	content := b.String()
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
