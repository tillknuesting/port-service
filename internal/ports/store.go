package ports

// package ports defines interfaces for external system
// integrations that are implemented by adapter layer.

import "context"

// Store is a generic persistence interface for saving,
// data elements. The T type parameter allows loose coupling for the value objects
// to be stored without assuming specific implementation.
type Store[T any] interface {
	// Set stores a value for a given key
	Set(ctx context.Context, key string, value T) error

	// Similarly Get and Delete would define
	// access capabilities.

	// Concrete implementing adapters would
	// provide specifics like serialization,
	// database/storage connectivity etc.
}
