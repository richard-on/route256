package repository

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/transactor"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/logger"
)

type Repository struct {
	transactor.QueryEngineProvider
	transactor.ExecEngineProvider
	log logger.Logger
}

func New(queryProvider transactor.QueryEngineProvider, execProvider transactor.ExecEngineProvider,
	log logger.Logger) *Repository {
	return &Repository{
		QueryEngineProvider: queryProvider,
		ExecEngineProvider:  execProvider,
		log:                 log,
	}
}
