// Package repository provides methods for database interaction.
package repository

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor"
)

type Repository struct {
	transactor.QueryEngineProvider
	transactor.ExecEngineProvider
}

func New(queryProvider transactor.QueryEngineProvider, execProvider transactor.ExecEngineProvider) *Repository {
	return &Repository{
		QueryEngineProvider: queryProvider,
		ExecEngineProvider:  execProvider,
	}
}
