package repository

import (
	"context"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
)

func LoadCustomTypes(ctx context.Context, conn *pgx.Conn) error {
	dataType, err := pgxtype.LoadDataType(ctx, conn, conn.ConnInfo(), "item")
	if err != nil {
		return err
	}
	conn.ConnInfo().RegisterDataType(dataType)

	dataType, err = pgxtype.LoadDataType(ctx, conn, conn.ConnInfo(), "_item")
	if err != nil {
		return err
	}
	conn.ConnInfo().RegisterDataType(dataType)

	return nil
}
