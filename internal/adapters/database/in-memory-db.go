package database

import (
	"context"
)

// MemDB represents a simplistic in-memory database abstraction in Go.
// It's a generic type, allowing it to store any type of value.
// The database is represented as a map, with string keys and generic type values.
type MemDB[T any] struct {
	// TODO: Add mutex
	DB map[string]T // Map acting as the in-memory storage.
}

// Set adds or updates a value in the in-memory database.
func (db *MemDB[T]) Set(ctx context.Context, key string, value T) error {
	db.DB[key] = value // Store or update the value in the map.
	return nil         // In this simple implementation, no error handling is performed.
}
