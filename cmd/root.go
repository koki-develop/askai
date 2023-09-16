package cmd

import (
	"os"

	"github.com/koki-develop/askai/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "askai",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: from config file
		key := os.Getenv("OPENAI_API_KEY")
		model := openai.GPT3Dot5Turbo

		ui := ui.New(&ui.Config{
			APIKey: key,
			Model:  model,
		})
		if err := ui.Start(); err != nil {
			return err
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
