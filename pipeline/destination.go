// Destination requirements:
// - storage in multiple destinations (i.e. in a file, written to a network, etc.)
// - elegantly handle failures
// - gracefully handle termination
// - provide reusable destination component for different destinations
// - extend the system as requirements evolve.
// - indepdently testable

// MVP:
// - write to a file

package pipeline

type Destination interface {
	// TODO: make sure it can only send.
	Messages(chan Message)
}
