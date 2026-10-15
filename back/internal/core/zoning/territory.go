package zoning

// territory is an area that UN M49 lists apart from its sovereign country, located by a
// bounding box in WGS 84 degrees. Every land inside a box belongs to the box's zone.
type territory struct {
	west, south, east, north float64
	zone                     ID
}

// overseasTerritories are the territories whose zone differs from their country's.
var overseasTerritories = []territory{
	{west: -55, south: 1, east: -51, north: 6.5, zone: LatinAmericaCaribbean},       // French Guiana
	{west: -62, south: 14, east: -60.5, north: 16.6, zone: LatinAmericaCaribbean},   // Guadeloupe, Martinique
	{west: -69, south: 11.9, east: -62.9, north: 18.1, zone: LatinAmericaCaribbean}, // Caribbean Netherlands
	{west: 44.9, south: -21.5, east: 55.9, north: -12.5, zone: Africa},              // Réunion, Mayotte
}

// ZoneOfTerritory returns the zone of an overseas territory containing the point.
func ZoneOfTerritory(longitude, latitude float64) (ID, bool) {
	for _, area := range overseasTerritories {
		if longitude >= area.west && longitude <= area.east && latitude >= area.south && latitude <= area.north {
			return area.zone, true
		}
	}
	return "", false
}
