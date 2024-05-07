package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/qualidafial/pomo/app"
	"github.com/qualidafial/pomo/store"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("accessing user home dir: %w", err))
	}

	dataDir := filepath.Join(homeDir, ".pomo")
	err = os.MkdirAll(dataDir, 0700)
	if err != nil {
		log.Fatal(fmt.Errorf("creating data dir: %w", err))
	}

	logFile := filepath.Join(dataDir, "log.txt")
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(fmt.Errorf("creating log file: %w", err))
	}
	defer func() {
		_ = f.Close()
	}()
	log.SetOutput(f)
	log.SetFormatter(log.TextFormatter)

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
