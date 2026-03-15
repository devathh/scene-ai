package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type Loader struct {
	validator *validator.Validate
}

func NewLoader() *Loader {
	return &Loader{
		validator: validator.New(),
	}
}

func (l *Loader) Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path to config file")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	data = []byte(os.ExpandEnv(string(data)))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := l.validator.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &cfg, nil
}
