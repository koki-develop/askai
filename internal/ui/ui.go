package ui

import (
	"errors"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"
)

var (
	youHeader = lipgloss.NewStyle().Background(lipgloss.Color("#00ADD8")).Foreground(lipgloss.Color("#000000")).Padding(0, 1).Render("You")
	aiHeader  = lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000")).Padding(0, 1).Render("AI")
)

type UI struct {
	writer io.Writer
	client *openai.Client
	model  string
}

type Config struct {
	APIKey string
	Model  string
}

func New(cfg *Config) *UI {
	client := openai.NewClient(cfg.APIKey)

	return &UI{
		writer: os.Stdout,
		client: client,
		model:  cfg.Model,
	}
}

func (ui *UI) Start() error {
	ctx := context.Background()

	messages := []openai.ChatCompletionMessage{}

	for {
		m := newInputModel()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			return err
		}
		if m.abort {
			break
		}
		ui.writer.Write([]byte(youHeader))
		ui.writer.Write([]byte{'\n'})
		ui.writer.Write([]byte(strings.TrimSpace(m.value)))
		ui.writer.Write([]byte{'\n', '\n'})
		messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: m.value})

		stream, err := ui.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Messages: messages,
			Model:    ui.model,
			Stream:   true,
		})
		if err != nil {
			return err
		}
		defer stream.Close()

		b := new(strings.Builder)
		ui.writer.Write([]byte(aiHeader))
		ui.writer.Write([]byte{'\n'})
		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			s := resp.Choices[0].Delta.Content
			b.WriteString(s)
			ui.writer.Write([]byte(s))
		}
		ui.writer.Write([]byte{'\n', '\n'})
		messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: b.String()})
	}

	return nil
}
