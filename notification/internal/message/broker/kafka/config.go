package kafka

import (
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
)

func NewSaramaConsumerConfig(appConfig config.KafkaConfig) *sarama.Config {
	c := sarama.NewConfig()

	c.Version = sarama.MaxVersion
	c.ClientID = appConfig.ClientID

	for _, strategy := range appConfig.GroupStrategies {
		var saramaStrategy sarama.BalanceStrategy
		switch strategy {
		case "range":
			saramaStrategy = sarama.BalanceStrategyRange
		case "roundRobin":
			saramaStrategy = sarama.BalanceStrategyRoundRobin
		case "sticky":
			saramaStrategy = sarama.BalanceStrategySticky
		default:
			saramaStrategy = sarama.BalanceStrategyRange
		}

		c.Consumer.Group.Rebalance.GroupStrategies = append(c.Consumer.Group.Rebalance.GroupStrategies, saramaStrategy)
	}
	c.Consumer.Offsets.Initial = int64(appConfig.InitialOffset)
	c.Consumer.Offsets.AutoCommit.Enable = appConfig.EnableAutoCommit

	c.Consumer.Return.Errors = true

	return c
}
