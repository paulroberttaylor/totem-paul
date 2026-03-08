package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paul/totem-trainer/internal/stats"
)

var (
	statsViewTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	statsHeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75"))
	statsRowStyle       = lipgloss.NewStyle()
	statsViewHelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

type statsViewModel struct {
	records []stats.SessionRecord
	err     error
}

func newStatsViewModel() statsViewModel {
	configDir := stats.ConfigDir()
	histPath := filepath.Join(configDir, "history.json")
	h := stats.NewHistory(histPath)
	records, err := h.Load()
	return statsViewModel{records: records, err: err}
}

func (m statsViewModel) update(msg tea.Msg) (statsViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, backToMenu()
		}
	}
	return m, nil
}

func (m statsViewModel) view(width, height int) string {
	var b strings.Builder

	b.WriteString(statsViewTitleStyle.Render("Session History"))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(fmt.Sprintf("  Error loading history: %v\n", m.err))
	} else if len(m.records) == 0 {
		b.WriteString("  No sessions recorded yet.\n")
	} else {
		// Show header.
		header := fmt.Sprintf("  %-20s %-14s %6s %8s %5s", "Timestamp", "Stage", "WPM", "Accuracy", "Keys")
		b.WriteString(statsHeaderStyle.Render(header))
		b.WriteString("\n")

		// Show last 10 sessions.
		start := 0
		if len(m.records) > 10 {
			start = len(m.records) - 10
		}
		recent := m.records[start:]

		for _, r := range recent {
			ts := r.Timestamp.Format("2006-01-02 15:04")
			line := fmt.Sprintf("  %-20s %-14s %6.0f %7.1f%% %5d", ts, r.Stage, r.WPM, r.Accuracy*100, r.TotalKeys)
			b.WriteString(statsRowStyle.Render(line))
			b.WriteString("\n")
		}

		// Rolling averages.
		if len(recent) > 0 {
			var totalWPM, totalAcc float64
			for _, r := range recent {
				totalWPM += r.WPM
				totalAcc += r.Accuracy
			}
			avgWPM := totalWPM / float64(len(recent))
			avgAcc := totalAcc / float64(len(recent)) * 100

			b.WriteString("\n")
			avgLine := fmt.Sprintf("  Rolling avg (last %d): WPM %.0f, Accuracy %.1f%%", len(recent), avgWPM, avgAcc)
			b.WriteString(statsHeaderStyle.Render(avgLine))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(statsViewHelpStyle.Render("esc: back"))

	content := b.String()
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
