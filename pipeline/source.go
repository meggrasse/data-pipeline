// Reader requirements:
// - process streaming data from multiple sources
// - different message types
// - elegantly handle failures
// - gracefully handle termination
// - provide reusable reader components for different sources
// - extend the system as requirements evolve.
// - indepdently testable

// MVP:
// - single message source
// - single message type

package pipeline

// import "errors"

type Result[T any] struct {
	Value T
	Error error
}

type Source interface {
	// This means that implmenenters of Source send values in a channel of Results.
	// A Result type contains either a Message or an error.
	// What is a Message? It is a type that represents the data being sent.
	// Also a key design feature is that it kicks off a goroutine (i.e. async reading)
    // Messages() <-chan Result[Message]
	// TODO, make sure it can only recieve.
	// TODO: handle errors
	Messages(chan Message)
}

// type Message interface {
	// For now, clients are required to pass data as a string.
	// We may later move to bytes. If they use a reader, they're responsible for siphoning into strings.
	// Data() string
	// Ideally, clients would register their schema with the system, and then use the schema type and version to validate the data.
	// SchemaType() string
	// SchemaVersion() string
// }

// Message validation will only occur at the processing stage.
