package ports

// package ports defines interfaces for external system
// integrations that are implemented by adapter layer.

import "context"

// Streamer is a generic interface that exposes a stream
// processing/consumption capability for objects of
// arbitrary types.
// Rather than directly depending on heavy streaming
// implementations, the domain layer interacts with a
// stream via this minimal interface to keep domain
// model isolated from those outer layer details.
// T represents the application domain object type
// we want streamed, allowing flexibility.
type Streamer[T any] interface {
	// StreamObjects returns for example a Go channel of T
	// that gets populated with a stream of objects.

	// StreamObjects bufferSize hints at desired channel buffer length.
	StreamObjects(ctx context.Context, bufferSize int) (<-chan T, error)
}
