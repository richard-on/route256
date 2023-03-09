package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/schema"
)

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
		return nil, errors.New("cart is empty")
	}

	return convert.ToModelItemSlice(items), nil
}
