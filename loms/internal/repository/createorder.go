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

func (r *Repository) InsertOrderInfo(ctx context.Context, order domain.OrderInfo) (int64, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	statement := sq.Insert("orders").
		Columns("user_id", "status", "created_at").
		Values(order.User, order.Status, time.Now()).
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

func (r *Repository) InsertOrderItems(ctx context.Context, orderID int64, domainItems []domain.Item) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	items := convert.ToSchemaItemArray(domainItems)

	statement := sq.Insert("order_items").
		Columns("order_id", "sku", "count")
	for _, item := range items {
		statement = statement.Values(orderID, item.SKU, item.Count)
	}
	statement = statement.PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if _, err = db.Exec(ctx, raw, args...); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ChangeOrderStatus(ctx context.Context, orderID int64, status domain.Status) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("orders").
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

func (r *Repository) DecreaseStock(ctx context.Context, sku int64, stock domain.Stock) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("stocks").
		Set("count", sq.Expr("count - ?", stock.Count)).
		Where(sq.Eq{"sku": sku}).
		Where(sq.Eq{"warehouse_id": stock.WarehouseID}).
		Where(sq.Gt{"count": 0}).
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

func (r *Repository) IncreaseStock(ctx context.Context, sku int64, stock domain.Stock) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Update("stocks").
		Set("count", sq.Expr("count + ?", stock.Count)).
		Where(sq.Eq{"sku": sku}).
		Where(sq.Eq{"warehouse_id": stock.WarehouseID}).
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

func (r *Repository) ReserveItem(ctx context.Context, orderID int64, sku int64, stock domain.Stock) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Insert("reserves").
		Values(orderID, sku, stock.WarehouseID, stock.Count).
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

type Reserve struct {
	SKU         int64
	WarehouseID int64
	Count       int64
}

func (r *Repository) RemoveItemsFromReserved(ctx context.Context, orderID int64) ([]int64, []domain.Stock, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	statement := sq.Delete("reserves").
		Where(sq.Eq{"order_id": orderID}).
		Suffix("RETURNING sku, warehouse_id, count").
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return nil, nil, err
	}

	var reserves []Reserve
	err = pgxscan.Select(ctx, db, &reserves, raw, args...)
	if err != nil {
		return nil, nil, err
	}

	skus := make([]int64, 0, len(reserves))
	stocks := make([]domain.Stock, 0, len(reserves))
	for _, res := range reserves {
		skus = append(skus, res.SKU)
		stocks = append(stocks, domain.Stock{
			WarehouseID: res.WarehouseID,
			Count:       uint64(res.Count),
		})
	}

	return skus, stocks, nil
}
