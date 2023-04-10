package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/logger/zerolog"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/message/broker/kafka"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/message/receiver"
)

func Run(cfg *config.Config) {
	log := zerolog.New(
		os.Stdout,
		cfg.Log.Level,
		cfg.Service.Name,
	)
	log.Info("config and logger init success")

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers,
		cfg.Kafka.ConsumerConfig.GroupID, kafka.NewSaramaConsumerConfig(cfg.Kafka.ConsumerConfig))
	if err != nil {
		log.Fatal(err, "creating kafka consumer group")
	}

	consumer := kafka.NewConsumer(log)

	statusReceiver := receiver.NewStatusReceiver(consumerGroup, consumer, log)

	model := domain.New(cfg.Service, statusReceiver)

	err = model.StatusReceiver.Subscribe(ctx, cfg.Kafka.Topic)
	if err != nil {
		log.Fatalf(err, "subscribing to kafka topic: %v", cfg.Kafka.Topic)
	}
	log.Info("successfully subscribed to kafka topic")

	<-ctx.Done()
	cancel()

	log.Info("shutting down: notification service")
}
