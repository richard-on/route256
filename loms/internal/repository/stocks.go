package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

func (r *Repository) GetStocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
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

	return convert.ToDomainStock(stocks), nil
}
