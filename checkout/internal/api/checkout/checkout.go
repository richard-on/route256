package checkout

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

// Checkout is a wrapper for checkout service gRPC server API.
type Checkout struct {
	checkout.UnimplementedCheckoutServer
	// domain represents business logic of checkout service.
	domain *domain.Domain
}

// New creates a new Checkout.
func New(domain *domain.Domain) *Checkout {
	return &Checkout{
		checkout.UnimplementedCheckoutServer{},
		domain,
	}
}
