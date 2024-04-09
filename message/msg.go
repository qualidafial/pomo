package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Save() tea.Msg {
	return SaveMsg{}
}

type SaveMsg struct{}

func Cancel() tea.Msg {
	return CancelMsg{}
}

type CancelMsg struct{}
