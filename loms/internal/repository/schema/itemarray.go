package schema

import (
	"encoding/binary"
	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type ItemArray struct {
	Elements   []Item `db:"_item"`
	Dimensions []pgtype.ArrayDimension
	Status     pgtype.Status
}

func (dst *ItemArray) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	var arrayHeader pgtype.ArrayHeader

	rp, err := arrayHeader.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	if len(arrayHeader.Dimensions) == 0 {
		*dst = ItemArray{Dimensions: arrayHeader.Dimensions, Status: pgtype.Present}
		return nil
	}

	elementCount := arrayHeader.Dimensions[0].Length
	for _, d := range arrayHeader.Dimensions[1:] {
		elementCount *= d.Length
	}

	elements := make([]Item, elementCount)

	for i := range elements {
		elemLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4
		var elemSrc []byte
		if elemLen >= 0 {
			elemSrc = src[rp : rp+elemLen]
			rp += elemLen
		}
		err = elements[i].DecodeBinary(ci, elemSrc)
		if err != nil {
			return err
		}
	}

	*dst = ItemArray{Elements: elements, Dimensions: arrayHeader.Dimensions, Status: pgtype.Present}
	return nil
}

func (src ItemArray) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	arrayHeader := pgtype.ArrayHeader{
		Dimensions: src.Dimensions,
	}

	if dt, ok := ci.DataTypeForName("item"); ok {
		arrayHeader.ElementOID = int32(dt.OID)
	} else {
		return nil, errors.New("unable to find oid for type name Item")
	}

	buf = arrayHeader.EncodeBinary(ci, buf)

	for i := range src.Elements {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		elemBuf, err := src.Elements[i].EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if elemBuf != nil {
			buf = elemBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
	}

	return buf, nil
}
