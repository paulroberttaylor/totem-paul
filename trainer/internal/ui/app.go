package ui

import (
	"fmt"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paul/totem-trainer/internal/keymap"
	"github.com/paul/totem-trainer/internal/lesson"
	"github.com/paul/totem-trainer/internal/stats"
)

// Screen identifiers.
type screen int

const (
	screenMenu screen = iota
	screenPicker
	screenTyping
	screenResults
	screenStats
)

// Messages for screen transitions.
type switchScreenMsg struct{ screen screen }
type startLessonMsg struct{ stage lesson.Stage }
type lessonDoneMsg struct{ session *stats.Session; stage string }
type backToMenuMsg struct{}

func switchScreen(s screen) tea.Cmd {
	return func() tea.Msg { return switchScreenMsg{screen: s} }
}

func startLesson(stage lesson.Stage) tea.Cmd {
	return func() tea.Msg { return startLessonMsg{stage: stage} }
}

func lessonDone(session *stats.Session, stage string) tea.Cmd {
	return func() tea.Msg { return lessonDoneMsg{session: session, stage: stage} }
}

func backToMenu() tea.Cmd {
	return func() tea.Msg { return backToMenuMsg{} }
}

// App is the top-level Bubble Tea model.
type App struct {
	keymap  *keymap.Keymap
	history *stats.History
	width   int
	height  int

	current screen
	menu    menuModel
	picker  pickerModel
	typing  typingModel
	results resultsModel
	stats   statsViewModel
}

// NewApp creates the top-level app model.
func NewApp(km *keymap.Keymap) App {
	configDir := stats.ConfigDir()
	histPath := filepath.Join(configDir, "history.jsonl")
	h := stats.NewHistory(histPath)

	return App{
		keymap:  km,
		history: h,
		current: screenMenu,
		menu:    newMenuModel(),
		picker:  newPickerModel(),
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

	case switchScreenMsg:
		a.current = msg.screen
		if msg.screen == screenStats {
			a.stats = newStatsViewModel()
		}
		return a, nil

	case startLessonMsg:
		stage := msg.stage
		var exercise string
		switch stage.Name {
		case "symbols":
			snippets := lesson.SymbolSnippets()
			exercise = lesson.GenerateExercise(snippets, 10)
		case "numbers":
			exercise = lesson.GenerateExercise(stage.Keys, 20)
		default:
			allowed := make(map[string]bool)
			for _, k := range stage.Keys {
				allowed[k] = true
			}
			allowed[" "] = true
			words := lesson.FilterWords(lesson.CommonWords(), allowed)
			exercise = lesson.GenerateExercise(words, 15)
			if exercise == "" {
				exercise = lesson.GenerateExercise(stage.Keys, 20)
			}
		}
		a.typing = newTypingModel(exercise, stage.Name)
		a.current = screenTyping
		return a, nil

	case lessonDoneMsg:
		sess := msg.session
		// Save session to history.
		record := stats.SessionRecord{
			Timestamp: time.Now(),
			Stage:     msg.stage,
			Duration:  sess.Duration.Seconds(),
			WPM:       sess.WPM(),
			Accuracy:  sess.Accuracy(),
			TotalKeys: sess.TotalChars,
			PerKey:    make(map[string]stats.KeyRecord),
		}
		for k, ks := range sess.PerKey {
			record.PerKey[k] = stats.KeyRecord{
				Hits:   ks.Hits,
				Misses: ks.Misses,
				AvgMs:  ks.AvgLatency(),
			}
		}
		_ = a.history.Append(record)

		a.results = newResultsModel(sess, msg.stage)
		a.current = screenResults
		return a, nil

	case backToMenuMsg:
		a.menu = newMenuModel()
		a.current = screenMenu
		return a, nil
	}

	// Delegate to current sub-model.
	var cmd tea.Cmd
	switch a.current {
	case screenMenu:
		a.menu, cmd = a.menu.update(msg)
	case screenPicker:
		a.picker, cmd = a.picker.update(msg)
	case screenStats:
		a.stats, cmd = a.stats.update(msg)
	case screenTyping:
		a.typing, cmd = a.typing.update(msg)
	case screenResults:
		a.results, cmd = a.results.update(msg)
	}
	return a, cmd
}

func (a App) View() string {
	switch a.current {
	case screenMenu:
		return a.menu.view(a.width, a.height)
	case screenPicker:
		return a.picker.view(a.width, a.height)
	case screenTyping:
		return a.typing.view(a.width, a.height)
	case screenResults:
		return a.results.view(a.width, a.height)
	case screenStats:
		return a.stats.view(a.width, a.height)
	}
	return fmt.Sprintf("unknown screen: %d", a.current)
}
