package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/convert"
)

// InsertOrderInfo inserts order information (excluding ordered items) to a database.
func (r *Repository) InsertOrderInfo(ctx context.Context, order model.Order) (int64, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	statement := sq.Insert("orders").
		Columns("user_id", "status", "created_at").
		Values(order.User, order.Status, time.Now()).
		Suffix("RETURNING order_id").
		PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return 0, err
	}

	var orderID int64
	if err = pgxscan.Get(ctx, db, &orderID, raw, args...); err != nil {
		return 0, err
	}

	return orderID, nil
}

// InsertOrderItems inserts ordered items to a database.
func (r *Repository) InsertOrderItems(ctx context.Context, orderID int64, items []model.Item) error {
	db := r.ExecEngineProvider.GetExecEngine(ctx)

	schemaItems := convert.ToSchemaItemSlice(items)

	statement := sq.Insert("order_items").
		Columns("order_id", "sku", "count")
	for _, item := range schemaItems {
		statement = statement.Values(orderID, item.SKU, item.Count)
	}
	statement = statement.PlaceholderFormat(sq.Dollar)

	raw, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	if _, err = db.Exec(ctx, raw, args...); err != nil {
		return err
	}

	return nil
}
