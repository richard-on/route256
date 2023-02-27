package loms

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// LOMS is a wrapper for Logistics and Order Management System gRPC server API.
type LOMS struct {
	loms.UnimplementedLOMSServer
	domain *domain.Domain
}

// New creates a new LOMS.
func New(domain *domain.Domain) *LOMS {
	return &LOMS{
		loms.UnimplementedLOMSServer{},
		domain,
	}
}
