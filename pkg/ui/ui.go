package ui

import (
	"github.com/ricoberger/httpmonitor/pkg/target"

	tea "github.com/charmbracelet/bubbletea"
)

func Start(targets []target.Client) error {
	model := NewModel(targets)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
