package outbox

const (
	NotStarted Status = iota
	Processing
	Failed
)

// Status is an enumeration that represents a status of message in outbox.
type Status uint8

type Message struct {
	ID      int64
	Key     string
	Payload []byte
	Status  Status
}
