package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paul/totem-trainer/internal/stats"
)

var (
	correctStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	wrongStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
	upcomingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	cursorStyle   = lipgloss.NewStyle().Reverse(true)
	wpmStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true)
	accStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	typingHelp    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(2)
)

type typingModel struct {
	exercise  string
	stage     string
	input     []rune   // what the user has typed (parallel to exercise runes)
	correct   []bool   // whether each typed rune was correct
	pos       int      // current cursor position in the exercise
	session   *stats.Session
	started   bool
	start     time.Time
	lastKey   time.Time
	done      bool
	lastTyped string   // label of last key pressed (for keyboard highlight)
}

func newTypingModel(exercise, stage string) typingModel {
	return typingModel{
		exercise: exercise,
		stage:    stage,
		input:    make([]rune, 0, len(exercise)),
		correct:  make([]bool, 0, len(exercise)),
		session:  stats.NewSession(),
	}
}

func (m typingModel) update(msg tea.Msg) (typingModel, tea.Cmd) {
	if m.done {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()

		// Handle backspace.
		if keyStr == "backspace" {
			if m.pos > 0 {
				m.pos--
				m.input = m.input[:m.pos]
				m.correct = m.correct[:m.pos]
			}
			return m, nil
		}

		// Handle escape to quit typing.
		if keyStr == "esc" {
			return m, backToMenu()
		}

		// Ignore non-character keys.
		var typed rune
		switch keyStr {
		case "enter":
			typed = '\n'
		case "tab":
			typed = '\t'
		case "space":
			typed = ' '
		default:
			runes := []rune(keyStr)
			if len(runes) != 1 {
				return m, nil
			}
			typed = runes[0]
		}

		// Start timer on first keypress.
		if !m.started {
			m.started = true
			m.start = time.Now()
			m.lastKey = m.start
		}

		exerciseRunes := []rune(m.exercise)
		if m.pos >= len(exerciseRunes) {
			return m, nil
		}

		now := time.Now()
		latency := now.Sub(m.lastKey)
		m.lastKey = now

		expected := exerciseRunes[m.pos]
		isCorrect := typed == expected
		m.input = append(m.input, typed)
		m.correct = append(m.correct, isCorrect)
		m.lastTyped = CharToLabel(typed)

		// Track stats.
		m.session.TotalChars++
		charKey := string(expected)
		if _, ok := m.session.PerKey[charKey]; !ok {
			m.session.PerKey[charKey] = &stats.KeyStats{}
		}
		if isCorrect {
			m.session.Correct++
			m.session.PerKey[charKey].RecordHit(latency)
		} else {
			m.session.PerKey[charKey].RecordMiss(latency)
		}

		m.pos++

		// Check if exercise is complete.
		if m.pos >= len(exerciseRunes) {
			m.done = true
			m.session.Duration = time.Since(m.start)
			return m, lessonDone(m.session, m.stage)
		}

		return m, nil
	}

	return m, nil
}

func (m typingModel) view(width, height int) string {
	var b strings.Builder

	exerciseRunes := []rune(m.exercise)

	// Render the exercise text with coloring.
	var textParts strings.Builder
	for i, ch := range exerciseRunes {
		s := string(ch)
		if i < m.pos {
			if m.correct[i] {
				textParts.WriteString(correctStyle.Render(s))
			} else {
				textParts.WriteString(wrongStyle.Render(s))
			}
		} else if i == m.pos {
			textParts.WriteString(cursorStyle.Render(s))
		} else {
			textParts.WriteString(upcomingStyle.Render(s))
		}
	}

	b.WriteString(textParts.String())
	b.WriteString("\n\n")

	// Live WPM and accuracy.
	var wpm float64
	var acc float64
	if m.started && m.pos > 0 {
		elapsed := time.Since(m.start)
		if m.done {
			elapsed = m.session.Duration
		}
		minutes := elapsed.Minutes()
		if minutes > 0 {
			wpm = (float64(m.pos) / 5.0) / minutes
		}
		acc = float64(m.session.Correct) / float64(m.session.TotalChars) * 100
	}

	statsLine := fmt.Sprintf("WPM: %.0f  Accuracy: %.0f%%  [%d/%d]", wpm, acc, m.pos, len(exerciseRunes))
	b.WriteString(wpmStyle.Render(statsLine))
	b.WriteString("\n")
	b.WriteString(typingHelp.Render("esc: quit"))

	b.WriteString("\n\n")
	layout := BaseLayer()
	layout.Highlight = m.lastTyped
	// Show next expected key
	if m.pos < len(exerciseRunes) {
		layout.Next = CharToLabel(exerciseRunes[m.pos])
	}
	b.WriteString(RenderKeyboard(layout))

	content := b.String()
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
