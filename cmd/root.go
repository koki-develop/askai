package cmd

import (
	"os"
	"strings"

	"github.com/koki-develop/askai/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagAPIKey      string // -k, --api-key
	flagModel       string // -m, --model
	flagInteractive bool   // -i, --interactive
)

type config struct {
	APIKey string `json:"api-key"`
	Model  string `json:"model"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure askai",
	Long:  "Configure askai.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config
		if err := viper.ReadInConfig(); err == nil {
			if err := viper.Unmarshal(&cfg); err != nil {
				return err
			}
		}

		if cmd.Flag("api-key").Changed {
			viper.Set("api-key", flagAPIKey)
		} else if cfg.APIKey != "" {
			flagAPIKey = cfg.APIKey
		}
		if cmd.Flag("model").Changed {
			viper.Set("model", flagModel)
		} else if cfg.Model != "" {
			flagModel = cfg.Model
		}
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		return nil
	},
}

var rootCmd = &cobra.Command{
	Use:   "askai [flags] [question]",
	Short: "AI is with you",
	Long:  "AI is with you.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &ui.Config{
			APIKey:      flagAPIKey,
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
	viper.SetConfigName(".askai")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	rootCmd.Flags().StringVarP(&flagAPIKey, "api-key", "k", "", "the OpenAI API key")
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", openai.GPT3Dot5Turbo, "the chat completion model to use")
	rootCmd.Flags().BoolVarP(&flagInteractive, "interactive", "i", false, "interactive mode")

	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&flagAPIKey, "api-key", "k", "", "the OpenAI API key")
	initCmd.Flags().StringVarP(&flagModel, "model", "m", "", "the chat completion model to use")
}
