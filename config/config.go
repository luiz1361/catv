package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Ollama struct {
		Model string `yaml:"model"`
		URL   string `yaml:"url"`
	} `yaml:"ollama"`
	Database struct {
		Name string `yaml:"name"`
	} `yaml:"database"`
}

// LoadConfig loads the configuration from the config.yaml file
func LoadConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
