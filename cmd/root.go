package cmd

import (
	"os"
	"strings"

	"github.com/koki-develop/askai/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	flagModel       string // -m, --model
	flagInteractive bool   // -i, --interactive
)

var rootCmd = &cobra.Command{
	Use:   "askai [flags] [question]",
	Short: "AI is with you.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: from config file
		key := os.Getenv("OPENAI_API_KEY")

		cfg := &ui.Config{
			APIKey:      key,
			Model:       flagModel,
			Interactive: flagInteractive,
		}

		q := strings.Join(args, " ")
		if q != "" {
			cfg.Question = &q
		}

		ui := ui.New(cfg)
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

func init() {
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", openai.GPT3Dot5Turbo, "the chat completion model to use")
	rootCmd.Flags().BoolVarP(&flagInteractive, "interactive", "i", false, "interactive mode")
}
