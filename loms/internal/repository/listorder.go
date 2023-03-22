package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

// ListOrderInfo gets order information (excluding ordered items) from a database.
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
	err = pgxscan.Get(ctx, db, &orderInfo, raw, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Order{}, domain.ErrEmptyOrder
	}
	if err != nil {
		return model.Order{}, err
	}

	return convert.ToModelOrder(orderInfo), nil
}

// ListOrderItems gets ordered items from a database.
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

// ListUnpaidOrders gets all order, that are awaiting payment for more than provided duration.
func (r *Repository) ListUnpaidOrders(ctx context.Context, paymentWait time.Duration) ([]int64, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select("order_id").
		From("orders").
		Where(sq.Lt{"created_at": time.Now().Add(-paymentWait)}).
		Where(sq.Eq{"status": model.AwaitingPayment}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var ids []int64
	if err = pgxscan.Select(ctx, db, &ids, raw, args...); err != nil {
		return nil, err
	}

	return ids, nil
}
