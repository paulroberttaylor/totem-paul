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
	resultsTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	resultGoodStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green >95%
	resultOkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow >85%
	resultBadStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red <85%
	resultsHelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

type resultsModel struct {
	session *stats.Session
	stage   string
}

func newResultsModel(session *stats.Session, stage string) resultsModel {
	return resultsModel{session: session, stage: stage}
}

func (m resultsModel) update(msg tea.Msg) (resultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "p":
			// Practice again: re-pick the same stage.
			// We go back to the picker screen.
			return m, switchScreen(screenPicker)
		case "esc", "q":
			return m, backToMenu()
		}
	}
	return m, nil
}

func (m resultsModel) view(width, height int) string {
	var b strings.Builder

	b.WriteString(resultsTitleStyle.Render("Results"))
	b.WriteString("\n\n")

	sess := m.session

	b.WriteString(fmt.Sprintf("  Stage:     %s\n", m.stage))
	b.WriteString(fmt.Sprintf("  WPM:       %.0f\n", sess.WPM()))
	b.WriteString(fmt.Sprintf("  Accuracy:  %.1f%%\n", sess.Accuracy()*100))
	b.WriteString(fmt.Sprintf("  Duration:  %.1fs\n", sess.Duration.Seconds()))
	b.WriteString(fmt.Sprintf("  Total keys: %d\n", sess.TotalChars))

	// Top 5 weakest keys.
	type keyAcc struct {
		key string
		acc float64
	}
	var keys []keyAcc
	for k, ks := range sess.PerKey {
		keys = append(keys, keyAcc{key: k, acc: ks.Accuracy()})
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].acc < keys[j].acc
	})

	b.WriteString("\n  Weakest keys:\n")
	limit := 5
	if len(keys) < limit {
		limit = len(keys)
	}
	for i := 0; i < limit; i++ {
		ka := keys[i]
		pct := ka.acc * 100
		display := ka.key
		if display == " " {
			display = "space"
		}
		line := fmt.Sprintf("    %s: %.0f%%", display, pct)
		switch {
		case pct >= 95:
			b.WriteString(resultGoodStyle.Render(line))
		case pct >= 85:
			b.WriteString(resultOkStyle.Render(line))
		default:
			b.WriteString(resultBadStyle.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(resultsHelpStyle.Render("enter/p: practice again • esc: menu"))

	content := b.String()
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
