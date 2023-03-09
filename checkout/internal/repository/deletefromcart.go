package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/convert"
)

// DecreaseItemCartCount removes a given number of item from user's cart.
func (r *Repository) DecreaseItemCartCount(ctx context.Context, userID int64, modelItem model.Item) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	item := convert.ToSchemaItem(modelItem)

	statement := sq.Update("cart_items").
		Set("count", sq.Expr("cart_items.count - ?", item.Count)).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sku": item.SKU}).
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

// DeleteItemCart removes item record from user's cart.
func (r *Repository) DeleteItemCart(ctx context.Context, userID int64, sku uint32) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Delete("cart_items").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sku": int64(sku)}).
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

// ClearCart clears cart for a given user.
func (r *Repository) ClearCart(ctx context.Context, userID int64) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	statement := sq.Delete("cart_items").
		Where(sq.Eq{"user_id": userID}).
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
