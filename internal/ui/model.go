package ui

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/sashabaranov/go-openai"
)

var (
	_ tea.Model = (*model)(nil)
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
		{k.Submit},
		{k.Quit},
	}
}

type model struct {
	program *tea.Program

	// openai api
	aiClient *openai.Client
	aiModel  string

	// states
	err               error
	receiving         bool
	resceivingMessage string
	messages          []openai.ChatCompletionMessage

	// components
	help           help.Model
	textarea       textarea.Model
	textareaKeyMap textareaKeyMap
}

func newModel(cfg *Config) *model {
	ta := textarea.New()
	ta.Placeholder = "Ask AI anything..."
	ta.ShowLineNumbers = false
	ta.Focus()

	return &model{
		// openai api
		aiClient: openai.NewClient(cfg.APIKey),
		aiModel:  cfg.Model,

		// states
		receiving: false,
		messages:  []openai.ChatCompletionMessage{},

		// components
		help:     help.New(),
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

func (m *model) Init() tea.Cmd {
	runewidth.DefaultCondition.EastAsianWidth = false
	return tea.Batch(
		textarea.Blink,
	)
}

type errMsg struct{ error }
type startReceivingMsg struct{}
type receivedMsg struct{ string }
type endReceivingMsg struct{}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.error
		return m, tea.Quit
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.textareaKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.textareaKeyMap.Submit):
			askm := m.textarea.Value()
			if strings.TrimSpace(askm) != "" {
				m.messages = append(m.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: askm})
				m.textarea.Reset()
				cmds = append(cmds, m.startReceiving)
			}
		}
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
	case startReceivingMsg:
		m.textarea.Blur()
		m.textareaKeyMap.Submit.SetEnabled(false)
		m.receiving = true
		cmds = append(cmds, m.receive)
	case receivedMsg:
		m.resceivingMessage = msg.string
		return m, nil
	case endReceivingMsg:
		m.textarea.Focus()
		m.textareaKeyMap.Submit.SetEnabled(true)
		m.receiving = false
		m.messages = append(m.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: m.resceivingMessage})
		m.resceivingMessage = ""
	}

	{
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) startReceiving() tea.Msg {
	return startReceivingMsg{}
}

func (m *model) receive() tea.Msg {
	ctx := context.Background()
	stream, err := m.aiClient.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Messages: m.messages,
		Model:    m.aiModel,
		Stream:   true,
	})
	if err != nil {
		return errMsg{err}
	}
	defer stream.Close()

	msg := new(strings.Builder)

	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errMsg{err}
		}
		msg.WriteString(resp.Choices[0].Delta.Content)
		m.program.Send(receivedMsg{msg.String()})
	}

	return endReceivingMsg{}
}

func (m *model) View() string {
	v := new(strings.Builder)
	youStyle := lipgloss.NewStyle().Background(lipgloss.Color("#00ADD8")).Foreground(lipgloss.Color("#000000")).Padding(0, 1)
	aiStyle := lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000")).Padding(0, 1)

	for _, msg := range m.messages {
		switch msg.Role {
		case openai.ChatMessageRoleUser:
			v.WriteString(youStyle.Render("You"))
		case openai.ChatMessageRoleAssistant:
			v.WriteString(aiStyle.Render("AI"))
		}

		v.WriteRune('\n')

		content, err := glamour.Render(msg.Content, "dark")
		if err != nil {
			m.program.Send(errMsg{err})
			return ""
		}
		v.WriteString(content)
	}

	if m.receiving {
		v.WriteString(aiStyle.Render("AI"))
		v.WriteRune('\n')
		content, err := glamour.Render(m.resceivingMessage, "dark")
		if err != nil {
			m.program.Send(errMsg{err})
			return ""
		}
		v.WriteString(content)
	} else {
		v.WriteString(youStyle.Render("You"))
		v.WriteString("\n\n")
		v.WriteString(m.textarea.View())
		v.WriteRune('\n')
		v.WriteString(m.help.View(m.textareaKeyMap))
		v.WriteRune('\n')
	}

	return v.String()
}
