package transactor

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

type txKey string

const key = txKey("tx")

type Transactor struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Transactor {
	return &Transactor{
		pool: pool,
	}
}

func (t *Transactor) RunReadCommitted(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *Transactor) RunRepeatableRead(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *Transactor) RunSerializable(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}
