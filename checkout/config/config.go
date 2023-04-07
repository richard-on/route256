// Package config provides structures for storing configurable variables and functions for parsing them.
package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	Version      string // Version of this app.
	Build        string // Build date and time.
	ProtoVersion string // ProtoVersion is the protobuf contract version.
)

type Config struct {
	Service             `yaml:"service"`
	Log                 `yaml:"log"`
	Observability       `yaml:"observability"`
	Postgres            `yaml:"postgres"`
	Kubernetes          `yaml:"kubernetes"`
	GRPC                `yaml:"grpc"`
	LOMS                `yaml:"loms"`
	ProductService      `yaml:"productService"`
	NotificationService `yaml:"notificationService"`
}

type Service struct {
	Name           string `yaml:"name"`
	Environment    string `yaml:"environment"`
	MaxPoolWorkers int    `yaml:"maxPoolWorkers"`
}

type Log struct {
	Level string `yaml:"logLevel"`
}

type Observability struct {
	Metrics `yaml:"metrics"`
	Jaeger  `yaml:"jaeger"`
}

type Metrics struct {
	Port string `yaml:"port"`
}

type Jaeger struct {
	SamplerType  string  `yaml:"samplerType"`
	SamplerParam float64 `yaml:"samplerParam"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslMode"`
}

type Kubernetes struct {
	Namespace      string        `yaml:"namespace"`
	LabelSelector  string        `yaml:"labelSelector"`
	UpdateInterval time.Duration `yaml:"updateInterval"`
}

type GRPC struct {
	Port string `yaml:"port"`
}

type LOMS struct {
	URL string `yaml:"url"`
}

type RateLimit struct {
	Rate  int `yaml:"rate"`
	Burst int `yaml:"burst"`
}

type ProductService struct {
	Token     string    `yaml:"token"`
	URL       string    `yaml:"url"`
	RateLimit RateLimit `yaml:"rateLimit"`
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
