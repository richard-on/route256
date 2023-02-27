package loms

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

type Implementation struct {
	loms.UnimplementedLOMSServer
	domain *domain.Domain
}

func New(domain *domain.Domain) *Implementation {
	return &Implementation{
		loms.UnimplementedLOMSServer{},
		domain,
	}
}
