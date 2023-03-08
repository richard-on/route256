package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

func (r *Repository) ListOrderInfo(ctx context.Context, orderID int64) (model.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select("status", "user_id").
		From("orders").
		Where(sq.Eq{"order_id": orderID}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := query.ToSql()
	if err != nil {
		return model.Order{}, err
	}

	var orderInfo schema.Order
	if err = pgxscan.Get(ctx, db, &orderInfo, raw, args...); err != nil {
		return model.Order{}, err
	}

	return convert.ToModelOrder(orderInfo), nil
}

func (r *Repository) ListOrderItems(ctx context.Context, orderID int64) ([]model.Item, error) {
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

	return convert.ToModelItemSlice(items), nil
}
