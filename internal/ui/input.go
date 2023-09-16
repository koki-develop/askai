package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	_ tea.Model = (*inputModel)(nil)
)

type inputModel struct {
	abort    bool
	quitting bool
	value    string
	textarea textarea.Model
}

func newInputModel() *inputModel {
	ta := textarea.New()
	ta.Placeholder = "Ask AI anything..."
	ta.ShowLineNumbers = false
	ta.Focus()

	return &inputModel{
		textarea: ta,
	}
}

func (m *inputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			m.abort = true
			return m, tea.Quit
		case tea.KeyCtrlD:
			m.quitting = true
			m.value = m.textarea.Value()
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *inputModel) View() string {
	if m.quitting {
		return ""
	}

	you := lipgloss.NewStyle().Background(lipgloss.Color("#00ADD8")).Foreground(lipgloss.Color("#000000")).Padding(0, 1).Render("You")
	return you + "\n" + m.textarea.View()
}
