package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model/outbox"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

func (r *Repository) AddMessageWithKey(ctx context.Context, key string, payload []byte) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Insert("outbox").
		Columns("key", "payload", "status").
		Values(key, payload, outbox.NotStarted).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	r.log.RawSQL("AddMessageWithKey", raw, args)

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("AddMessageWithKey", tag, err)
		return err
	}
	r.log.PGTag("AddMessageWithKey", tag)

	return nil
}

func (r *Repository) AddMessageWithoutKey(ctx context.Context, payload []byte) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Insert("outbox").
		Columns("payload", "status").
		Values(payload, outbox.NotStarted).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	r.log.RawSQL("AddMessageWithoutKey", raw, args)

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("AddMessageWithoutKey", tag, err)
		return err
	}
	r.log.PGTag("AddMessageWithoutKey", tag)

	return nil
}

func (r *Repository) UpdateMessageStatus(ctx context.Context, id int64, status outbox.Status) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("outbox").
		Set("status", status).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	r.log.RawSQL("UpdateMessageStatus", raw, args)

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("UpdateMessageStatus", tag, err)
		return err
	}
	r.log.PGTag("UpdateMessageStatus", tag)
	if tag.RowsAffected() == 0 {
		return domain.ErrRecordNotExists
	}

	return nil
}

func (r *Repository) DeleteMessage(ctx context.Context, id int64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Delete("outbox").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	r.log.RawSQL("DeleteMessage", raw, args)

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("DeleteMessage", tag, err)
		return err
	}
	r.log.PGTag("DeleteMessage", tag)
	if tag.RowsAffected() == 0 {
		return domain.ErrRecordNotExists
	}

	return nil
}

func (r *Repository) ListUnsent(ctx context.Context) ([]outbox.Message, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select("id", "key", "payload", "status").
		From("outbox").
		Where(sq.Or{
			sq.Eq{"status": outbox.NotStarted},
			sq.Eq{"status": outbox.Failed},
		}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var messages []schema.Message
	err = pgxscan.Select(ctx, db, &messages, raw, args...)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return convert.ToOutboxMessageSlice(messages), nil
}
