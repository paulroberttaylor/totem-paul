package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paul/totem-trainer/internal/keymap"
	"github.com/paul/totem-trainer/internal/ui"
)

func main() {
	keymapPath := "config/totem.keymap"
	if len(os.Args) > 1 {
		keymapPath = os.Args[1]
	}

	km, err := keymap.ParseFile(keymapPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading keymap %s: %v\n", keymapPath, err)
		os.Exit(1)
	}

	app := ui.NewApp(km)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
