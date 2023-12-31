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
	writer      io.Writer
	client      *openai.Client
	model       string
	interactive bool
	question    *string
	messages    []openai.ChatCompletionMessage
}

type Config struct {
	APIKey      string
	Model       string
	Interactive bool
	Question    *string
	Messages    []openai.ChatCompletionMessage
}

func New(cfg *Config) *UI {
	client := openai.NewClient(cfg.APIKey)

	return &UI{
		writer:      os.Stdout,
		client:      client,
		model:       cfg.Model,
		interactive: cfg.Interactive,
		question:    cfg.Question,
		messages:    cfg.Messages,
	}
}

func (ui *UI) Start() error {
	ctx := context.Background()

	if !ui.interactive && ui.question == nil {
		return errors.New("question is required when interactive mode is disabled")
	}

	for {
		var msg string
		if ui.question == nil {
			ipt, ok, err := ui.readInput()
			if err != nil {
				return err
			}
			if !ok {
				break
			}
			msg = ipt
		} else {
			msg = *ui.question
			ui.question = nil
		}

		ui.messages = append(ui.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg})
		if ui.interactive {
			_, _ = ui.writer.Write([]byte(youHeader))
			_, _ = ui.writer.Write([]byte{'\n'})
			_, _ = ui.writer.Write([]byte(strings.TrimSpace(msg)))
			_, _ = ui.writer.Write([]byte{'\n', '\n'})
		}

		ans, err := ui.printAnswer(ctx)
		if err != nil {
			return err
		}

		if !ui.interactive {
			break
		}

		ui.messages = append(ui.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: ans})
	}

	return nil
}

func (ui *UI) readInput() (string, bool, error) {
	m := newInputModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return "", false, err
	}
	if m.abort {
		return "", false, nil
	}
	return m.value, true, nil
}

func (ui *UI) printAnswer(ctx context.Context) (string, error) {
	stream, err := ui.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Messages: ui.messages,
		Model:    ui.model,
		Stream:   true,
	})
	if err != nil {
		return "", err
	}
	defer stream.Close()

	b := new(strings.Builder)
	if ui.interactive {
		_, _ = ui.writer.Write([]byte(aiHeader))
		_, _ = ui.writer.Write([]byte{'\n'})
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		s := resp.Choices[0].Delta.Content
		b.WriteString(s)
		_, _ = ui.writer.Write([]byte(s))
	}
	if ui.interactive {
		_, _ = ui.writer.Write([]byte{'\n'})
	}
	_, _ = ui.writer.Write([]byte{'\n'})

	return b.String(), nil
}
