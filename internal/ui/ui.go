package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UI struct {
	model *model
}

func New() *UI {
	return &UI{
		model: &model{},
	}
}

func (*UI) Start() error {
	p := tea.NewProgram(&model{})
	return p.Start()
}
