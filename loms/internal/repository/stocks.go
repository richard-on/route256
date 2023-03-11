package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

// GetStocks returns a slice of model.Stock, containing availability information for a given item.
func (r *Repository) GetStocks(ctx context.Context, sku uint32) ([]model.Stock, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select("warehouse_id", "count").
		From("stocks").
		Where(sq.Eq{"sku": sku}).
		Where(sq.Gt{"count": 0}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var stocks []schema.Stock
	if err = pgxscan.Select(ctx, db, &stocks, raw, args...); err != nil {
		return nil, err
	}

	return convert.ToModelStockSlice(stocks), nil
}

// IncreaseStock increases stock count for a given item.
func (r *Repository) IncreaseStock(ctx context.Context, sku int64, stock model.Stock) error {
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

// DecreaseStock decreases stock count for a given item.
func (r *Repository) DecreaseStock(ctx context.Context, sku int64, stock model.Stock) error {
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
