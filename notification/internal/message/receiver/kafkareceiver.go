package receiver

import (
	"context"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/pkg/logger"
	"sync"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/internal/message/broker/kafka"
)

type StatusReceiver struct {
	consumerGroup  sarama.ConsumerGroup
	statusConsumer kafka.Consumer
	log            logger.Logger
}

func NewStatusReceiver(consumerGroup sarama.ConsumerGroup, consumer kafka.Consumer, log logger.Logger) *StatusReceiver {
	return &StatusReceiver{
		consumerGroup:  consumerGroup,
		statusConsumer: consumer,
		log:            log,
	}
}

func (r *StatusReceiver) Subscribe(ctx context.Context, topic string) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := r.consumerGroup.Consume(ctx, []string{topic}, &r.statusConsumer); err != nil {
				r.log.Errorf(err, "error consuming kafka topic: %v", topic)
			}
			// check if context was cancelled, signaling that the consumer should stop.
			if ctx.Err() != nil {
				r.log.Debug("got signal to stop consumer")
				return
			}
		}
	}()

	<-r.statusConsumer.Ready() // Wait till the consumer has been set up.
	r.log.Info("sarama consumer up and running")

	wg.Wait()
	if err := r.consumerGroup.Close(); err != nil {
		r.log.Fatal(err, "closing client")
	}

	return nil
}
