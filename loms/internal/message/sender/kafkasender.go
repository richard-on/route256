package sender

import (
	"runtime"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
)

type Sender struct {
	asyncProducer sarama.AsyncProducer
	config        config.Kafka
	onSuccess     HandlerFunc
	onFail        HandlerFunc

	wg sync.WaitGroup
}

type HandlerFunc func(id int64)

type StatusSenderOption func(*Sender)

func WithSuccessFunc(fn HandlerFunc) StatusSenderOption {
	return func(s *Sender) {
		s.onSuccess = fn
	}
}

func WithFailFunc(fn HandlerFunc) StatusSenderOption {
	return func(s *Sender) {
		s.onFail = fn
	}
}

func NewStatusSender(producer sarama.AsyncProducer, cfg config.Kafka, options ...StatusSenderOption) *Sender {
	s := &Sender{
		asyncProducer: producer,
		config:        cfg,
		wg:            sync.WaitGroup{},
	}

	for _, opt := range options {
		opt(s)
	}

	s.wg.Add(1)
	go func() {
		for m := range producer.Successes() {
			id := m.Metadata.(int64)
			if s.onSuccess != nil {
				s.onSuccess(id)
			}
		}
		s.wg.Done()
	}()

	s.wg.Add(1)
	go func() {
		for e := range producer.Errors() {
			id := e.Msg.Metadata.(int64)
			if s.onFail != nil {
				s.onFail(id)
			}
		}
		s.wg.Done()
	}()

	return s
}

// Close calls AsyncClose() on producer and then waits for Success and Errors channels to close and drain.
func (s *Sender) Close() {
	s.asyncProducer.AsyncClose()
	s.wg.Wait()
}

func (s *Sender) SendWithKey(id int64, key string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: s.config.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("go_ver"),
				Value: []byte(runtime.Version()),
			},
			{
				Key:   []byte("app_ver"),
				Value: []byte(s.config.Headers.AppVersion),
			},
			{
				Key:   []byte("app_build"),
				Value: []byte(s.config.Headers.AppBuild),
			},
			{
				Key:   []byte("proto_ver"),
				Value: []byte(s.config.Headers.ProtoVersion),
			},
		},
		Metadata:  id,
		Partition: -1,
		Timestamp: time.Now(),
	}

	s.asyncProducer.Input() <- msg
}
