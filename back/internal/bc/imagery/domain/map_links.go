package domain

import "strconv"

// MapLinks points to the exact spot on map services; there is none unless the point is
// where the image was taken.
func (l Location) MapLinks() []Link {
	if l.Granularity != GranularityPoint || l.Point == nil {
		return []Link{}
	}

	latitude, longitude := degrees(l.Point.Latitude), degrees(l.Point.Longitude)
	return []Link{
		{Label: "GeoHack", URL: "https://geohack.toolforge.org/geohack.php?params=" + geoHackParams(*l.Point)},
		{
			Label: "OpenStreetMap",
			URL:   "https://www.openstreetmap.org/?mlat=" + latitude + "&mlon=" + longitude + "#map=16/" + latitude + "/" + longitude,
		},
		{Label: "Google Maps", URL: "https://www.google.com/maps/search/?api=1&query=" + latitude + "%2C" + longitude},
	}
}

// geoHackParams writes a point as GeoHack expects it, e.g. "48.85_N_2.35_E".
func geoHackParams(point Point) string {
	return hemisphere(point.Latitude, "N", "S") + "_" + hemisphere(point.Longitude, "E", "W")
}

func hemisphere(value float64, positive, negative string) string {
	if value < 0 {
		return degrees(-value) + "_" + negative
	}
	return degrees(value) + "_" + positive
}

func degrees(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
