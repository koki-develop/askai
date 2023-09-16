package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	_ tea.Model = (*model)(nil)
)

type model struct{}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	return "Hello, World!"
}
