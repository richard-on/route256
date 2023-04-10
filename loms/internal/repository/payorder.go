package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

// PayOrder sets order status to "paid".
func (r *Repository) PayOrder(ctx context.Context, orderID int64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("orders").
		Set("status", model.Paid).
		Where(sq.Eq{"order_id": orderID}).
		Where(sq.Eq{"status": model.AwaitingPayment}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("PayOrder", tag, err)
		return err
	}
	r.log.PGTag("PayOrder", tag)
	if tag.RowsAffected() == 0 {
		return domain.ErrNotExistsOrPaid
	}

	return nil
}
