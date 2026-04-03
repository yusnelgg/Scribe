package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AI AIConfig `yaml:"ai"`
}

type AIConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func LoadDefault() (*Config, error) {
	paths := []string{".scribe.yaml", "scribe.yaml", ".scribe.yml", "scribe.yml"}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return Load(path)
		}
	}

	return nil, nil
}

func (c *Config) GetAIProvider() string {
	if c == nil || c.AI.Provider == "" {
		return "ollama"
	}
	return c.AI.Provider
}

func (c *Config) GetAIAPIKey() string {
	if c == nil {
		return ""
	}
	return c.AI.APIKey
}

func (c *Config) GetAIEndpoint() string {
	if c == nil || c.AI.Endpoint == "" {
		return "http://localhost:11434"
	}
	return c.AI.Endpoint
}

func (c *Config) GetAIModel() string {
	if c == nil || c.AI.Model == "" {
		return ""
	}
	return c.AI.Model
}
