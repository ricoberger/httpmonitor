package ui

import (
	"github.com/ricoberger/httpmonitor/pkg/target"

	tea "charm.land/bubbletea/v2"
)

func Start(targets []target.Client) error {
	model := NewModel(targets)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
