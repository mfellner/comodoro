package app

// A FleetUnit represents a fleet unit.
type FleetUnit struct {
	Name string                 `json:"name"`
	Body map[string]interface{} `json:"body"`
}

// Units is a collection of fleet units.
type FleetUnits []FleetUnit
