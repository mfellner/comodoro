package model

// A Unit represents a fleet unit.
type Unit struct {
	Name string            `json:"name"`
	Body map[string]string `json:"body"`
}

// Units is a collection of fleet units.
type Units []Unit
