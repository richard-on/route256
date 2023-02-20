// Package config provides structures for storing configurable variables and functions for parsing them.
package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Service             `yaml:"service"`
	Log                 `yaml:"log"`
	HTTP                `yaml:"http"`
	LOMS                `yaml:"loms"`
	ProductService      `yaml:"productService"`
	NotificationService `yaml:"notificationService"`
}

type Service struct {
	Name string `yaml:"name"`
}

type Log struct {
	Level string `yaml:"logLevel"`
}

type HTTP struct {
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

type LOMS struct {
	URL string `yaml:"url"`
}

type ProductService struct {
	Token string `yaml:"token"`
	URL   string `yaml:"url"`
}

type NotificationService struct{}

// New reads and returns app config.
func New() (*Config, error) {
	cfg := &Config{}

	rawYAML, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing yaml")
	}

	return cfg, nil
}
