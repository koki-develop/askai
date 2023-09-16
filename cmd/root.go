package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "CLI_NAME", // TODO
	RunE: func(cmd *cobra.Command, args []string) error {
		q := strings.Join(args, " ")

		ctx := context.Background()

		key := os.Getenv("OPENAI_API_KEY")
		client := openai.NewClient(key)

		// TODO: from config
		model := openai.GPT3Dot5Turbo

		stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: q},
			},
			Model:  model,
			Stream: true,
		})
		if err != nil {
			return err
		}
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}

			fmt.Print(resp.Choices[0].Delta.Content)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
