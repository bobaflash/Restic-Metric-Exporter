package internal

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Secret string `yaml:"secret"`
}

func NewConfig() (*Config, error) {

	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		slog.Info("No config file path defined. Using default path './config.yaml'")
		path = "./config.yaml"
	}

	slog.Info("Read config file", "path", path)

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := yaml.Unmarshal(f, &conf); err != nil {
		return nil, err
	}

	return &conf, nil

}
