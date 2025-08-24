package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WorkerURL       string `yaml:"worker_url"`
	ConfirmURL      string `yaml:"confirm_url"`
	PollingInterval int    `yaml:"polling_interval"`
	Port            string `yaml:"port"`
	LogLevel        string `yaml:"log_level"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
