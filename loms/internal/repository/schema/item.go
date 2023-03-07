package schema

import (
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type Item struct {
	SKU   int64 `db:"sku"`
	Count int32 `db:"count"`
}

func (dst *Item) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return errors.New("NULL values can't be decoded. Scan into a &*Item to handle NULLs")
	}

	if err := (pgtype.CompositeFields{&dst.SKU, &dst.Count}).DecodeBinary(ci, src); err != nil {
		return err
	}

	return nil
}

func (src Item) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	sku := pgtype.Int8{Int: src.SKU, Status: pgtype.Present}
	count := pgtype.Int4{Int: src.Count, Status: pgtype.Present}

	return (pgtype.CompositeFields{&sku, &count}).EncodeBinary(ci, buf)
}
