package config

import (
	"errors"
	"fmt"
	"os"
)

type Constructor struct {
	loader *Loader
}

func NewConstructor() *Constructor {
	return &Constructor{
		loader: NewLoader(),
	}
}

func (c *Constructor) Init(defaultPath string) (*Config, error) {
	path := os.Getenv("APP_CONFIG_PATH")
	if path == "" {
		if defaultPath == "" {
			return nil, errors.New("empty path to config file")
		}
		path = defaultPath
	}

	cfg, err := c.loader.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
