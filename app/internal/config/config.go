package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

var (
	ErrConfigFileDoesNotExist = errors.New("config file does not exist")
)

type Config struct {
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`

	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`

	JWT struct {
		TokenSecondsTTL int    `yaml:"token_seconds_ttl"`
		SecretKey       string `yaml:"secret_key"`
	} `yaml:"jwt"`
}

func New(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileDoesNotExist, configPath)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
