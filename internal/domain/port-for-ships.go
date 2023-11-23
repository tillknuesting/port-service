package domain

// Package domain contains the core business entities and logic.
// In Domain-Driven Design (DDD), this package is the heart of the business logic,
// encapsulating the domain model and rules.

// Port is the aggregate root and an entity in the domain model of a logistics or
// maritime system. In DDD, an aggregate root is a main entity within an aggregate,
// a cluster of domain objects that can be treated as a single unit for data changes.
// The Port aggregate would include all the domain logic and rules applicable to a port.

type Port struct {
	Key         string    `json:"key"`         // Unique identifier for the Port.
	Name        string    `json:"name"`        // Human-readable name of the Port.
	City        string    `json:"city"`        // City where the Port is located.
	Country     string    `json:"country"`     // Country where the Port is located.
	Alias       []string  `json:"alias"`       // Alternative names or identifiers for the Port.
	Regions     []string  `json:"regions"`     // Geographical or administrative regions associated with the Port.
	Coordinates []float64 `json:"coordinates"` // Geographical coordinates of the Port, typically latitude and longitude.
	Province    string    `json:"province"`    // Province or state where the Port is located.
	Timezone    string    `json:"timezone"`    // Time zone of the Port.
	Unlocs      []string  `json:"unlocs"`      // United Nations Location Codes for the Port.
	Code        string    `json:"code"`        // Additional coding system
}
