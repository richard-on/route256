package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

// CancelOrder sets order status to "cancelled".
func (r *Repository) CancelOrder(ctx context.Context, orderID int64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("orders").
		Set("status", model.Cancelled).
		Where(sq.Eq{"order_id": orderID}).
		Where(sq.Or{
			sq.Eq{"status": model.AwaitingPayment},
			sq.Eq{"status": model.NewOrder},
		}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	r.log.RawSQL("CancelOrder", raw, args)

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("CancelOrder", tag, err)
		return err
	}
	r.log.PGTag("CancelOrder", tag)
	if tag.RowsAffected() == 0 {
		return domain.ErrOrderCancelled
	}

	return nil
}
