package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"time"
)

func (r *Repository) CreateOrder(ctx context.Context, order domain.OrderInfo) (int64, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	items := convert.ToSchemaItemArray(order.Items)

	statement := sq.Insert("\"order\"").
		Columns("user_id", "items", "status", "created_at").
		Values(sq.Expr("$1", order.User), sq.Expr("$2::item[]", items),
			sq.Expr("$3", order.Status), sq.Expr("$4", time.Now())).
		Suffix("RETURNING order_id").
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return 0, err
	}

	var orderID int64
	if err = pgxscan.Get(ctx, db, &orderID, raw, args...); err != nil {
		return 0, err
	}

	return orderID, nil
}

type AvailabilityInfo struct {
	WarehouseID int64 `db:"warehouse_id"`
	Available   int32 `db:"available"`
}

func (r *Repository) ChangeOrderStatus(ctx context.Context, orderID int64, status domain.Status) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("\"order\"").
		Set("status", status).
		Where(sq.Eq{"order_id": orderID}).
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
		return errors.New("order does not exist")
	}

	return nil
}

func (r *Repository) CheckAvailability(ctx context.Context, domainItem domain.Item) ([]AvailabilityInfo, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	item := convert.ToSchemaItem(domainItem)

	statement := sq.Select("warehouse_id", "count - reserved AS available").
		From("stocks").
		Where(sq.Eq{"sku": item.SKU}).
		Where(
			sq.Expr("? >= ?",
				sq.Select("SUM(?)").
					From("stocks").
					Where(sq.Eq{"sku": item.SKU}),
				item.Count,
			),
		)

	raw, args, err := statement.ToSql()
	if err != nil {
		return nil, err
	}

	var availabilityInfo []AvailabilityInfo
	if err = pgxscan.Get(ctx, db, &availabilityInfo, raw, args...); err != nil {
		return nil, err
	}

	return availabilityInfo, nil
}

func (r *Repository) DecreaseCount(ctx context.Context, warehouseID int64, sku uint32, count uint64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("stocks").
		Set("count", sq.Expr("count - ?", count)).
		Where(sq.Eq{"warehouse_id": warehouseID}).
		Where(sq.Eq{"sku": sku}).
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
		return errors.New("warehouse or sku does not exist")
	}

	return nil
}

func (r *Repository) ReserveItem(ctx context.Context, warehouseID int64, sku uint32, reserveAmount uint64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("stocks").
		Set("reserved", reserveAmount).
		Where(sq.Eq{"warehouse_id": warehouseID}).
		Where(sq.Eq{"sku": sku}).
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
		return errors.New("warehouse or sku does not exist")
	}

	return nil
}
