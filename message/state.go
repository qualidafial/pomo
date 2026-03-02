package message

import (
	tea "charm.land/bubbletea/v2"
	"github.com/qualidafial/pomo"
)

func LoadState(current pomo.Pomo, previous []pomo.Pomo) tea.Cmd {
	return func() tea.Msg {
		return LoadStateMsg{
			Current:  current,
			Previous: previous,
		}
	}
}

type LoadStateMsg struct {
	Current  pomo.Pomo
	Previous []pomo.Pomo
}
