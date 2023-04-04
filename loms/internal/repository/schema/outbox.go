package schema

type Message struct {
	ID      int64   `db:"id"`
	Key     *string `db:"key"`
	Payload []byte  `db:"payload"`
	Status  int8    `db:"status"`
}
