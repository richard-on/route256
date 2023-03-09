// Package schema defines main models for repository layer.
package schema

// Item represents a product as it is stored in database.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU int64 `db:"sku"`
	// Count is the number of product's with this SKU.
	Count int32 `db:"count"`
}

/*func (dst *Item) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
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
}*/
