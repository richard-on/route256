package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/convert"
)

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
		return 0, errors.New("item is not in the cart")
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) DecreaseItemCartCount(ctx context.Context, userID int64, modelItem model.Item) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	item := convert.ToSchemaItem(modelItem)

	statement := sq.Insert("cart_items").
		Columns("user_id", "sku", "count").
		Values(userID, item.SKU, item.Count).
		Suffix("ON CONFLICT (user_id, sku) DO UPDATE SET count = cart_items.count + ?").
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
