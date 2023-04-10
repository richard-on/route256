package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

// ChangeOrderStatus sets order status to a provided model.Status.
func (r *Repository) ChangeOrderStatus(ctx context.Context, orderID int64, status model.Status) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("orders").
		Set("status", status).
		Where(sq.Eq{"order_id": orderID}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("ChangeOrderStatus", tag, err)
		return err
	}
	r.log.PGTag("ChangeOrderStatus", tag)
	if tag.RowsAffected() == 0 {
		return domain.ErrEmptyOrder
	}

	return nil
}
