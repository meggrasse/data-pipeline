package pipeline

type Message struct {
	// TODO: bytes.
	Payload string
	// TODO: use uuid
	ID int
	SchemaType string
	SchemaVersion string
}

