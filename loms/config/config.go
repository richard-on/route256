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
	Postgres            `yaml:"postgres"`
	GRPC                `yaml:"grpc"`
	Kafka               `yaml:"kafka"`
	Checkout            `yaml:"checkout"`
	ProductService      `yaml:"productService"`
	NotificationService `yaml:"notificationService"`
}

type Service struct {
	Name           string        `yaml:"name"`
	Environment    string        `yaml:"environment"`
	PaymentTimeout time.Duration `yaml:"paymentTimeout"`
	CancelInterval time.Duration `yaml:"cancelInterval"`
	SendInterval   time.Duration `yaml:"sendInterval"`
	MaxPoolWorkers int           `yaml:"maxPoolWorkers"`
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

type Kafka struct {
	Topic          string       `yaml:"topic"`
	Brokers        []string     `yaml:"brokers"`
	ProducerConfig KafkaConfig  `yaml:"producerConfig"`
	Headers        KafkaHeaders `yaml:"headers"`
}

type KafkaConfig struct {
	ClientID        string `yaml:"clientId"`
	Partitioner     string `yaml:"partitioner"`
	RequiredAcks    int    `yaml:"requiredAcks"`
	Idempotent      bool   `yaml:"idempotent"`
	MaxOpenRequests int    `yaml:"maxOpenRequests"`
}

type KafkaHeaders struct {
	AppVersion   string `yaml:"appVersion"`
	AppBuild     string `yaml:"appBuild"`
	ProtoVersion string `yaml:"protoVersion"`
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

	rawYAML, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing yaml")
	}

	cfg.Headers.AppVersion = Version
	cfg.Headers.AppBuild = Build
	cfg.Headers.ProtoVersion = ProtoVersion

	return cfg, nil
}
