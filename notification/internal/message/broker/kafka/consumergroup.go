package kafka

import (
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"gitlab.ozon.dev/rragusskiy/homework-1/notification/pkg/logger"
	"google.golang.org/protobuf/proto"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan struct{}
	log   logger.Logger
}

// NewConsumer - constructor
func NewConsumer(log logger.Logger) Consumer {
	return Consumer{
		ready: make(chan struct{}),
		log:   log,
	}
}

func (c *Consumer) Ready() <-chan struct{} {
	return c.ready
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msg := &loms.OrderStatus{}
	for {
		select {
		case message := <-claim.Messages():
			err := proto.Unmarshal(message.Value, msg)
			if err != nil {
				return err
			}
			c.log.Infof("message claimed: value = %s, timestamp = %v, topic = %s",
				msg, message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
