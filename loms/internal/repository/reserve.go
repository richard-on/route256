package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

// ReserveItem inserts item to reserves table.
func (r *Repository) ReserveItem(ctx context.Context, orderID int64, sku int64, stock model.Stock) error {
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

// RemoveItemsFromReserved removes items from a given order from reserves table.
func (r *Repository) RemoveItemsFromReserved(ctx context.Context, orderID int64) ([]int64, []model.Stock, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	statement := sq.Delete("reserves").
		Where(sq.Eq{"order_id": orderID}).
		Suffix("RETURNING sku, warehouse_id, count").
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return nil, nil, err
	}

	var reserves []schema.Reserve
	err = pgxscan.Select(ctx, db, &reserves, raw, args...)
	if err != nil {
		return nil, nil, err
	}

	skus := make([]int64, 0, len(reserves))
	stocks := make([]model.Stock, 0, len(reserves))
	for _, res := range reserves {
		skus = append(skus, res.SKU)
		stocks = append(stocks, model.Stock{
			WarehouseID: res.WarehouseID,
			Count:       uint64(res.Count),
		})
	}

	return skus, stocks, nil
}
