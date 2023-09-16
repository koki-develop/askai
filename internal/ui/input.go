package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	_ tea.Model = (*inputModel)(nil)
)

type textareaKeyMap struct {
	Submit key.Binding
	Quit   key.Binding
}

func (k textareaKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Submit, k.Quit}
}

func (k textareaKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Submit, k.Quit},
	}
}

type inputModel struct {
	abort    bool
	quitting bool
	value    string

	help           help.Model
	textarea       textarea.Model
	textareaKeyMap textareaKeyMap
}

func newInputModel() *inputModel {
	ta := textarea.New()
	ta.Placeholder = "Ask AI anything..."
	ta.ShowLineNumbers = false
	ta.Focus()

	return &inputModel{
		help:     help.NewModel(),
		textarea: ta,
		textareaKeyMap: textareaKeyMap{
			Submit: key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("Ctrl+d", "Submit message"),
			),
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("Ctrl+c/esc", "Quit the program"),
			),
		},
	}
}

func (m *inputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.textareaKeyMap.Quit):
			m.quitting = true
			m.abort = true
			return m, tea.Quit
		case key.Matches(msg, m.textareaKeyMap.Submit):
			if strings.TrimSpace(m.textarea.Value()) != "" {
				m.quitting = true
				m.value = m.textarea.Value()
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
	}

	return m, cmd
}

func (m *inputModel) View() string {
	if m.quitting {
		return ""
	}

	return youHeader + "\n" + m.textarea.View() + "\n" + m.help.View(m.textareaKeyMap)
}
