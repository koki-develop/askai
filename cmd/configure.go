package cmd

import (
	"fmt"
	"syscall"

	"github.com/koki-develop/askai/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
