package transactor

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine
}

type ExecEngine interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type ExecEngineProvider interface {
	GetExecEngine(ctx context.Context) ExecEngine
}

func (t *Transactor) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return t.pool
}

func (t *Transactor) GetExecEngine(ctx context.Context) ExecEngine {
	tx, ok := ctx.Value(key).(ExecEngine)
	if ok && tx != nil {
		return tx
	}

	return t.pool
}
