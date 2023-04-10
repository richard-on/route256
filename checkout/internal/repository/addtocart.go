package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/convert"
)

// AddToCart adds provided item to user's cart.
func (r *Repository) AddToCart(ctx context.Context, userID int64, modelItem model.Item) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	item := convert.ToSchemaItem(modelItem)

	statement := sq.Insert("cart_items").
		Columns("user_id", "sku", "count").
		Values(userID, item.SKU, item.Count).
		Suffix("ON CONFLICT (user_id, sku) DO UPDATE SET count = cart_items.count + ?", item.Count).
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, raw, args...)
	if err != nil {
		r.log.PGTag("AddToCart", tag, err)
		return err
	}
	r.log.PGTag("AddToCart", tag)

	return nil
}
