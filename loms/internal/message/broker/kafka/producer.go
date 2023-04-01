package kafka

import (
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
)

func NewAsyncProducer(config config.Kafka) (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(config.Brokers, NewSaramaConfig(config.ProducerConfig))
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func NewSaramaConfig(appConfig config.KafkaConfig) *sarama.Config {
	c := sarama.NewConfig()

	c.Version = sarama.MaxVersion
	c.ClientID = appConfig.ClientID

	switch appConfig.Partitioner {
	case "hash":
		c.Producer.Partitioner = sarama.NewHashPartitioner
	case "roundrobin":
		c.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "manual":
		c.Producer.Partitioner = sarama.NewManualPartitioner
	case "random":
		c.Producer.Partitioner = sarama.NewRandomPartitioner
	default:
		c.Producer.Partitioner = sarama.NewHashPartitioner
	}

	c.Producer.RequiredAcks = sarama.RequiredAcks(appConfig.RequiredAcks)
	c.Producer.Idempotent = appConfig.Idempotent
	c.Net.MaxOpenRequests = appConfig.MaxOpenRequests

	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true

	return c
}
