package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/schema"
)

// GetCartItems returns a slice of model.Item, which contains all items currently in user's cart.
func (r *Repository) GetCartItems(ctx context.Context, userID int64) ([]model.Item, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	statement := sq.Select("sku", "count").
		From("cart_items").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return nil, err
	}

	var items []schema.Item
	if err = pgxscan.Select(ctx, db, &items, raw, args...); err != nil {
		return nil, err
	}

	if items == nil {
		return nil, domain.ErrEmptyCart
	}

	return convert.ToModelItemSlice(items), nil
}

// GetItemCartCount returns a number of item with given sku in user's cart.
func (r *Repository) GetItemCartCount(ctx context.Context, userID int64, modelItem model.Item) (int32, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	item := convert.ToSchemaItem(modelItem)

	statement := sq.Select("count").
		From("cart_items").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sku": item.SKU}).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return 0, err
	}

	var count int32
	err = pgxscan.Get(ctx, db, &count, raw, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, domain.ErrNotInCart
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}
