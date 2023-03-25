package transactor

//go:generate minimock -i DB -o ./mocks/ -s "_minimock.go"

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"
)

type txKey string

const Key = txKey("tx")

type Transactor struct {
	db DB
}

type DB interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func New(pool DB) *Transactor {
	return &Transactor{
		db: pool,
	}
}

func (t *Transactor) RunReadCommitted(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, Key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *Transactor) RunRepeatableRead(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, Key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *Transactor) RunSerializable(ctx context.Context, f func(txCtx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, Key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}
