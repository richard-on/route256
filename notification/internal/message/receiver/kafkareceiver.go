package receiver

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/message/broker/kafka"
)

type StatusReceiver struct {
	consumerGroup  sarama.ConsumerGroup
	statusConsumer kafka.Consumer
}

func NewStatusReceiver(consumerGroup sarama.ConsumerGroup, consumer kafka.Consumer) *StatusReceiver {
	return &StatusReceiver{
		consumerGroup:  consumerGroup,
		statusConsumer: consumer,
	}
}

func (r *StatusReceiver) Subscribe(ctx context.Context, topic string) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := r.consumerGroup.Consume(ctx, []string{topic}, &r.statusConsumer); err != nil {
				log.Error().Err(err).Msgf("error consuming kafka topic: %v", topic)
			}
			// check if context was cancelled, signaling that the consumer should stop.
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-r.statusConsumer.Ready() // Wait till the consumer has been set up.
	log.Info().Msg("sarama consumer up and running")

	wg.Wait()
	if err := r.consumerGroup.Close(); err != nil {
		log.Fatal().Err(err).Msg("error closing client")
	}

	return nil

}
