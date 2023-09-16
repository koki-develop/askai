package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/koki-develop/askai/internal/config"
	"github.com/koki-develop/askai/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

/*
 * TODO: refactor
 */

var (
	flagGlobal      bool   // -g, --global
	flagConfigure   bool   // --configure
	flagAPIKey      string // -k, --api-key
	flagModel       string // -m, --model
	flagInteractive bool   // -i, --interactive
)

func configure(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}

	if !cmd.Flag("api-key").Changed && !cmd.Flag("model").Changed {
		fmt.Print("OpenAI API Key: ")
		key, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		if len(key) != 0 {
			cfg.APIKey = string(key)
		}

		fmt.Print("\nChat Completion Model: ")
		var m string
		fmt.Scanln(&m)
		if m != "" {
			cfg.Model = m
		}
	}

	if cmd.Flag("api-key").Changed {
		cfg.APIKey = flagAPIKey
	}
	if cmd.Flag("model").Changed {
		cfg.Model = flagModel
	}

	p, err := config.Save(cfg, flagGlobal)
	if err != nil {
		return err
	}

	fmt.Printf("Configuration saved to %s\n", p)
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "askai [flags] [question]",
	Short: "AI is with you",
	Long:  "AI is with you.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagConfigure {
			return configure(cmd, args)
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		uicfg := &ui.Config{
			APIKey:      cfg.APIKey,
			Model:       cfg.Model,
			Interactive: flagInteractive,
		}

		if cmd.Flag("api-key").Changed {
			uicfg.APIKey = flagAPIKey
		}
		if strings.TrimSpace(uicfg.APIKey) == "" {
			return fmt.Errorf("OpenAI API Key is required")
		}

		if cmd.Flag("model").Changed {
			uicfg.Model = flagModel
		}
		if uicfg.Model == "" {
			uicfg.Model = openai.GPT3Dot5Turbo
		}

		q := strings.Join(args, " ")
		if q != "" {
			uicfg.Question = &q
		}

		ui := ui.New(uicfg)
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
	rootCmd.Flags().BoolVarP(&flagGlobal, "global", "g", false, "configure askai globally (only for --configure)")
	rootCmd.Flags().BoolVar(&flagConfigure, "configure", false, "configure askai")
	rootCmd.Flags().StringVarP(&flagAPIKey, "api-key", "k", "", "the OpenAI API key")
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", openai.GPT3Dot5Turbo, "the chat completion model to use")
	rootCmd.Flags().BoolVarP(&flagInteractive, "interactive", "i", false, "interactive mode")
}
