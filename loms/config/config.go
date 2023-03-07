// Package config provides structures for storing configurable variables and functions for parsing them.
package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Service             `yaml:"service"`
	Log                 `yaml:"log"`
	Postgres            `yaml:"postgres"`
	GRPC                `yaml:"grpc"`
	Checkout            `yaml:"checkout"`
	ProductService      `yaml:"productService"`
	NotificationService `yaml:"notificationService"`
}

type Service struct {
	Name string `yaml:"name"`
}

type Log struct {
	Level string `yaml:"logLevel"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslMode"`
}

type GRPC struct {
	Port string `yaml:"port"`
}

type Checkout struct {
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
	//b, err := exec.Command("cmd", "/C", "dir").Output()
	//fmt.Println(string(b))
	rawYAML, err := os.ReadFile("./loms/config/config.yaml")
	if err != nil {
		return nil, errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing yaml")
	}

	return cfg, nil
}
