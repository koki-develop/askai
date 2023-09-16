package ui

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
)

var (
	_ tea.Model = (*model)(nil)
)

type model struct {
	program *tea.Program

	// openai api
	aiClient *openai.Client
	aiModel  string

	// states
	err               error
	askMessage        string
	receiving         bool
	resceivingMessage string
	messages          []openai.ChatCompletionMessage

	// components
	textarea textarea.Model
}

func newModel(cfg *Config) *model {
	ta := textarea.New()
	ta.Focus()

	return &model{
		// openai api
		aiClient: openai.NewClient(cfg.APIKey),
		aiModel:  cfg.Model,

		// states
		askMessage: "",
		receiving:  false,
		messages:   []openai.ChatCompletionMessage{},

		// components
		textarea: ta,
	}
}

func (m *model) Init() tea.Cmd {
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

	{
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.error
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlD:
			askm := m.textarea.Value()
			m.messages = append(m.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: askm})
			m.askMessage = askm
			m.textarea.Reset()
			cmds = append(cmds, m.startReceiving)
		}
	case startReceivingMsg:
		m.textarea.Blur()
		m.receiving = true
		cmds = append(cmds, m.receive)
	case receivedMsg:
		m.resceivingMessage = msg.string
	case endReceivingMsg:
		m.textarea.Focus()
		m.receiving = false
		m.messages = append(m.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: m.resceivingMessage})
		m.resceivingMessage = ""
	}

	return m, tea.Batch(cmds...)
}

func (m *model) startReceiving() tea.Msg {
	return startReceivingMsg{}
}

func (m *model) receive() tea.Msg {
	ctx := context.Background()
	stream, err := m.aiClient.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: m.askMessage},
		},
		Model:  m.aiModel,
		Stream: true,
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

	for _, msg := range m.messages {
		v.WriteString(msg.Content)
		v.WriteRune('\n')
	}

	if m.receiving {
		v.WriteString(m.resceivingMessage)
	} else {
		v.WriteString(m.textarea.View())
		v.WriteRune('\n')
	}

	return v.String()
}
