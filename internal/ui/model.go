package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	openai "github.com/sashabaranov/go-openai"
)

var (
	_ tea.Model = (*model)(nil)
)

type model struct {
	// states
	receiving bool
	messages  []openai.ChatCompletionMessage

	// components
	textarea textarea.Model
}

func newModel() *model {
	ta := textarea.New()
	ta.Focus()

	return &model{
		receiving: false,
		textarea:  ta,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	{
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlD:
			m.messages = append(m.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: m.textarea.Value()})
			m.textarea.Reset()
			m.receiving = true
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	v := new(strings.Builder)

	if m.receiving {
		for _, msg := range m.messages {
			v.WriteString(msg.Content)
			v.WriteRune('\n')
		}
		v.WriteString("...")
	} else {
		v.WriteString(m.textarea.View())
		v.WriteRune('\n')
		for _, msg := range m.messages {
			v.WriteString(msg.Content)
			v.WriteRune('\n')
		}
	}

	return v.String()
}
