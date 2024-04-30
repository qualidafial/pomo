package main

import (
	"fmt"
	"github.com/qualidafial/pomo/store"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qualidafial/pomo/app"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("accessing user home dir: %w", err))
	}
	dataDir := filepath.Join(homeDir, ".pomo")

	s, err := store.New(dataDir)
	if err != nil {
		log.Fatal(fmt.Errorf("creating pomo data store: %w", err))
	}

	p := tea.NewProgram(app.New(s))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
