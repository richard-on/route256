package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/notification/config"
)

type StatusReceiver interface {
	Subscribe(ctx context.Context, topic string) error
}

// Domain represents business logic of Notification service.
// It should wrap interfaces used in a service.
type Domain struct {
	config config.Service
	StatusReceiver
}

// New creates a new Domain.
func New(config config.Service, receiver StatusReceiver) *Domain {
	return &Domain{
		config,
		receiver,
	}
}

// NewMockDomain creates a new mock Domain used for testing.
func NewMockDomain(opts ...any) *Domain {
	d := Domain{}

	for _, v := range opts {
		switch s := v.(type) {
		case config.Service:
			d.config = s
		case StatusReceiver:
			d.StatusReceiver = s
		}
	}

	return &d
}
