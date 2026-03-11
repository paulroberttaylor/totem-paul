package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paul/totem-trainer/internal/stats"
)

var (
	statsViewTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	statsHeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75"))
	statsRowStyle       = lipgloss.NewStyle()
	statsViewHelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)

	// Heatmap intensity colors (dark to bright green, like GitHub)
	heatNone = lipgloss.NewStyle().Foreground(lipgloss.Color("236")) // empty
	heat1    = lipgloss.NewStyle().Foreground(lipgloss.Color("22"))  // dark green
	heat2    = lipgloss.NewStyle().Foreground(lipgloss.Color("28"))  // medium green
	heat3    = lipgloss.NewStyle().Foreground(lipgloss.Color("34"))  // bright green
	heat4    = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))  // vivid green
	heatDay  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	heatMon  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type statsViewModel struct {
	records []stats.SessionRecord
	err     error
}

func newStatsViewModel() statsViewModel {
	configDir := stats.ConfigDir()
	histPath := filepath.Join(configDir, "history.jsonl")
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
		// Activity heatmap
		b.WriteString(m.renderHeatmap())
		b.WriteString("\n\n")

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

func (m statsViewModel) renderHeatmap() string {
	// Build a map of date -> session count
	dayCounts := make(map[string]int)
	for _, r := range m.records {
		day := r.Timestamp.Format("2006-01-02")
		dayCounts[day]++
	}

	// Find max for scaling
	maxCount := 0
	for _, c := range dayCounts {
		if c > maxCount {
			maxCount = c
		}
	}

	// Show last 20 weeks (140 days) — fits nicely in a terminal
	now := time.Now()
	weeks := 20

	// Find the start: go back to the most recent Sunday, then back `weeks` weeks
	// GitHub-style: columns are weeks, rows are days (Mon-Sun)
	end := now
	// Walk back to Sunday
	for end.Weekday() != time.Saturday {
		end = end.AddDate(0, 0, 1)
	}
	start := end.AddDate(0, 0, -(weeks*7 - 1))

	// Month labels
	var monthRow strings.Builder
	monthRow.WriteString("       ") // indent for day labels
	lastMonth := -1
	for w := 0; w < weeks; w++ {
		day := start.AddDate(0, 0, w*7)
		month := int(day.Month())
		if month != lastMonth {
			monthRow.WriteString(heatMon.Render(day.Format("Jan")))
			lastMonth = month
			// Pad to align — "Jan" is 3 chars, each week column is 2 chars
			// We used 3 chars, so skip ahead
		} else {
			monthRow.WriteString("  ")
		}
	}

	// Day labels and grid
	dayNames := []string{"Mon", "   ", "Wed", "   ", "Fri", "   ", "Sun"}
	var rows [7]strings.Builder
	for i := range rows {
		rows[i].WriteString(heatDay.Render(fmt.Sprintf("  %-4s", dayNames[i])))
	}

	block := "█"

	for w := 0; w < weeks; w++ {
		for d := 0; d < 7; d++ {
			day := start.AddDate(0, 0, w*7+d)
			if day.After(now) {
				rows[d].WriteString("  ")
				continue
			}

			key := day.Format("2006-01-02")
			count := dayCounts[key]

			var style lipgloss.Style
			if count == 0 {
				style = heatNone
			} else if maxCount <= 1 || count == 1 {
				style = heat1
			} else {
				ratio := float64(count) / float64(maxCount)
				switch {
				case ratio < 0.33:
					style = heat2
				case ratio < 0.66:
					style = heat3
				default:
					style = heat4
				}
			}

			rows[d].WriteString(style.Render(block) + " ")
		}
	}

	var b strings.Builder
	b.WriteString("  " + statsHeaderStyle.Render("Practice Activity") + "\n\n")
	b.WriteString(monthRow.String() + "\n")
	for _, row := range rows {
		b.WriteString(row.String() + "\n")
	}

	// Total sessions and streak
	totalSessions := len(m.records)
	streak := m.currentStreak(now)
	b.WriteString(fmt.Sprintf("\n  %d sessions total", totalSessions))
	if streak > 0 {
		b.WriteString(fmt.Sprintf("  •  %d day streak", streak))
	}

	return b.String()
}

func (m statsViewModel) currentStreak(now time.Time) int {
	// Count consecutive days with sessions, ending today or yesterday
	dayCounts := make(map[string]bool)
	for _, r := range m.records {
		dayCounts[r.Timestamp.Format("2006-01-02")] = true
	}

	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	// Start from today or yesterday
	var check time.Time
	if dayCounts[today] {
		check = now
	} else if dayCounts[yesterday] {
		check = now.AddDate(0, 0, -1)
	} else {
		return 0
	}

	streak := 0
	for {
		key := check.Format("2006-01-02")
		if !dayCounts[key] {
			break
		}
		streak++
		check = check.AddDate(0, 0, -1)
	}

	return streak
}
