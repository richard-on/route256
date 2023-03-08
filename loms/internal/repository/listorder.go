package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

func (r *Repository) ListOrderInfo(ctx context.Context, orderID int64) (domain.OrderInfo, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select("status", "user_id").
		From("orders").
		Where(sq.Eq{"order_id": orderID}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := query.ToSql()
	if err != nil {
		return domain.OrderInfo{}, err
	}

	var orderInfo schema.OrderInfo
	if err = pgxscan.Get(ctx, db, &orderInfo, raw, args...); err != nil {
		return domain.OrderInfo{}, err
	}

	return convert.ToDomainOrderInfo(orderInfo), nil
}

func (r *Repository) ListOrderItems(ctx context.Context, orderID int64) ([]domain.Item, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select("sku", "count").
		From("order_items").
		Where(sq.Eq{"order_id": orderID}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.Item
	if err = pgxscan.Select(ctx, db, &items, raw, args...); err != nil {
		return nil, err
	}

	return convert.ToDomainItems(items), nil
}
