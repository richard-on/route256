package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
)

func (r *Repository) CancelOrder(ctx context.Context, orderID int64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("\"order\"").
		Set("status", domain.Cancelled).
		Where(sq.Eq{"order_id": orderID}).
		Where(sq.NotEq{"status": domain.Cancelled}).
		Where(sq.NotEq{"status": domain.Failed}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	exec, err := db.Exec(ctx, raw, args...)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return errors.New("order does not exist or has already been cancelled")
	}

	return nil
}
