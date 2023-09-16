package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UI struct {
	program *tea.Program
	model   *model
}

type Config struct {
	APIKey string
	Model  string
}

func New(cfg *Config) *UI {
	m := newModel(cfg)
	p := tea.NewProgram(m)
	m.program = p

	return &UI{
		program: p,
		model:   m,
	}
}

func (ui *UI) Start() error {
	if _, err := ui.program.Run(); err != nil {
		return err
	}
	if ui.model.err != nil {
		return ui.model.err
	}

	return nil
}
