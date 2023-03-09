package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/convert"
)

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

	if _, err = db.Exec(ctx, raw, args...); err != nil {
		return err
	}

	return nil
}
