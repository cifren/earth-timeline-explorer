// Package zoning defines the zones of the world and the table placing countries in them.
//
// Zones follow the UN M49 regions, with the Americas and Asia split along M49
// sub-regions. Countries are an internal reference: clients only see zones.
package zoning

import "slices"

// ID identifies a zone.
type ID string

// Zone identifiers, in display order.
const (
	NorthernAmerica       ID = "northern-america"
	LatinAmericaCaribbean ID = "latin-america-caribbean"
	Europe                ID = "europe"
	Africa                ID = "africa"
	MiddleEast            ID = "middle-east"
	SouthAndCentralAsia   ID = "south-central-asia"
	EastAndSouthEastAsia  ID = "east-southeast-asia"
	Oceania               ID = "oceania"
	Antarctica            ID = "antarctica"
)

// Zone is a selectable area of the world.
type Zone struct {
	ID   ID
	Name string
}

var zones = []Zone{
	{NorthernAmerica, "Northern America"},
	{LatinAmericaCaribbean, "Latin America and the Caribbean"},
	{Europe, "Europe"},
	{Africa, "Africa"},
	{MiddleEast, "Middle East"},
	{SouthAndCentralAsia, "South and Central Asia"},
	{EastAndSouthEastAsia, "East and South-East Asia"},
	{Oceania, "Oceania"},
	{Antarctica, "Antarctica"},
}

// All returns every zone in display order.
func All() []Zone {
	return slices.Clone(zones)
}

// Find returns the zone with the given identifier.
func Find(id ID) (Zone, bool) {
	index := slices.IndexFunc(zones, func(zone Zone) bool { return zone.ID == id })
	if index < 0 {
		return Zone{}, false
	}
	return zones[index], true
}
