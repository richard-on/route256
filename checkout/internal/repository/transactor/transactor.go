package transactor

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

const (
	txKey = "tx"
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

type QueryRowEngine interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type QueryRowEngineProvider interface {
	GetQueryRowEngine(ctx context.Context) QueryRowEngine
}

type TransactionManager struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (t *TransactionManager) RunRepeatableRead(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, txKey, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(txKey).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return t.pool
}

func (t *TransactionManager) GetQueryRowEngine(ctx context.Context) QueryRowEngine {
	tx, ok := ctx.Value(txKey).(QueryRowEngine)
	if ok && tx != nil {
		return tx
	}

	return t.pool
}

func (t *TransactionManager) GetExecEngine(ctx context.Context) ExecEngine {
	tx, ok := ctx.Value(txKey).(ExecEngine)
	if ok && tx != nil {
		return tx
	}

	return t.pool
}
