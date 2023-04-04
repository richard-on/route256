package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	Version string // Version of this app.
	Build   string // Build date and time.
)

type Config struct {
	Service `yaml:"service"`
	Log     `yaml:"log"`
	Kafka   `yaml:"kafka"`
}

type Service struct {
	Name string `yaml:"name"`
}

type Log struct {
	Level string `yaml:"logLevel"`
}

type Kafka struct {
	Topic          string      `yaml:"topic"`
	Brokers        []string    `yaml:"brokers"`
	ConsumerConfig KafkaConfig `yaml:"consumerConfig"`
}

type KafkaConfig struct {
	ClientID         string   `yaml:"clientId"`
	GroupID          string   `yaml:"groupId"`
	GroupStrategies  []string `yaml:"groupStrategies"`
	InitialOffset    int      `yaml:"initialOffset"`
	EnableAutoCommit bool     `yaml:"enableAutoCommit"`
}

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
