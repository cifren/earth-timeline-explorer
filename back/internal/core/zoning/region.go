package zoning

// subregionZones places UN M49 sub-regions whose zone differs from their region's.
var subregionZones = map[string]ID{
	"Northern America":   NorthernAmerica,
	"Central America":    LatinAmericaCaribbean,
	"South America":      LatinAmericaCaribbean,
	"Caribbean":          LatinAmericaCaribbean,
	"Western Asia":       MiddleEast,
	"Southern Asia":      SouthAndCentralAsia,
	"Central Asia":       SouthAndCentralAsia,
	"Eastern Asia":       EastAndSouthEastAsia,
	"South-Eastern Asia": EastAndSouthEastAsia,
}

// regionZones places the UN M49 regions that form a zone as a whole.
var regionZones = map[string]ID{
	"Europe":     Europe,
	"Africa":     Africa,
	"Oceania":    Oceania,
	"Antarctica": Antarctica,
}

// ZoneOfRegion returns the zone of a UN M49 region and sub-region, named in English.
func ZoneOfRegion(region, subregion string) (ID, bool) {
	if zone, found := subregionZones[subregion]; found {
		return zone, true
	}
	zone, found := regionZones[region]
	return zone, found
}
