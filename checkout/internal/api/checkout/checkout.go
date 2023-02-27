package checkout

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

type Implementation struct {
	checkout.UnimplementedCheckoutServer
	domain *domain.Domain
}

func New(domain *domain.Domain) *Implementation {
	return &Implementation{
		checkout.UnimplementedCheckoutServer{},
		domain,
	}
}
