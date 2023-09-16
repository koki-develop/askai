package config

import (
	"os"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey   string   `yaml:"api_key"`
	Model    string   `yaml:"model"`
	Messages Messages `yaml:"messages"`
}

type Message struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

type Messages []Message

func (m Messages) OpenAI() []openai.ChatCompletionMessage {
	msgs := make([]openai.ChatCompletionMessage, len(m))
	for i, msg := range m {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return msgs
}

func Load() (*Config, error) {
	p, _ := configPath(false)

	ok, err := isExists(p)
	if err != nil {
		return nil, err
	}
	if ok {
		cfg, err := load(p)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	p, err = configPath(true)
	if err != nil {
		return nil, err
	}

	ok, err = isExists(p)
	if err != nil {
		return nil, err
	}
	if ok {
		cfg, err := load(p)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	return &Config{}, nil
}

func Save(cfg *Config, global bool) (string, error) {
	p, err := configPath(global)
	if err != nil {
		return "", err
	}

	f, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(cfg); err != nil {
		return "", err
	}

	return p, nil
}

func configPath(global bool) (string, error) {
	p := ".askai"
	if global {
		h, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(h, p)
	}
	return p, nil
}

func load(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func isExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
