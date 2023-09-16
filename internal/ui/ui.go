package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UI struct {
	model *model
}

func New() *UI {
	return &UI{
		model: newModel(),
	}
}

func (ui *UI) Start() error {
	p := tea.NewProgram(ui.model)
	return p.Start()
}
