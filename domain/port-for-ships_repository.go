package domain

// Package domain contains the core business entities and logic.
// In Domain-Driven Design (DDD), this package is the heart of the business logic,
// encapsulating the domain model and rules.

import (
	"context"
	"fmt"

	"ports-service/internal/ports"
)

// PortRepository is an interface defining the contract for persistence operations
// related to the Port entity, which is the aggregate root in this domain.
// In DDD, repositories abstract away the details of data persistence, allowing the
// domain model (entities/aggregates) to be designed and evolved without being
// coupled to specific database technologies or schemas.
//
// The repository interface resides within the domain package, close to the
// aggregate root (Port), emphasizing the principle of Ubiquitous Language and
// keeping domain logic cohesive. This closeness to the aggregate root reinforces
// the idea that saving and retrieving aggregates should be domain-centric,
// preserving the integrity and invariants of the aggregate.
//
// Additionally, having the repository interface within the domain package aligns
// with the DDD concept of a Bounded Context, ensuring that the domain layer
// remains focused on business logic and rules, while the implementation details
// of persistence (how data is stored, queried, etc.) are handled by the infrastructure layer.
type PortRepository interface {
	// Store persists a Port aggregate, handling its state storage in a way
	// that is consistent with the domain model's invariants and business rules.
	// The use of context.Context allows for operation cancellation, deadlines,
	// and passing request-scoped values, making the method more robust and flexible.
	Store(context.Context, Port) error
}

type StorePortRepository struct {
	Data ports.Store[Port]
}

func (s StorePortRepository) Store(ctx context.Context, port Port) error {
	if err := s.Data.Set(ctx, port.Key, port); err != nil {
		return fmt.Errorf("method of PortRepository Store can not Set data: %w", err)
	}

	return nil
}
