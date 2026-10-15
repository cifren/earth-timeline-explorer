package zoning

import (
	"maps"
	"slices"

	"frise-chrono/back/var/geozone"
)

// countryZones places each country in its zone. Keys are English short names as used
// by Wikimedia Commons categories ("the Netherlands", "Georgia (country)"), including
// historical entities, which are placed in the zone of their present-day territory.
var countryZones = map[string]ID{
	// Northern America
	"Canada": NorthernAmerica, "the United States": NorthernAmerica, "Greenland": NorthernAmerica,
	"Bermuda": NorthernAmerica, "Saint Pierre and Miquelon": NorthernAmerica,

	// Latin America and the Caribbean
	"Mexico": LatinAmericaCaribbean, "Guatemala": LatinAmericaCaribbean, "Belize": LatinAmericaCaribbean,
	"El Salvador": LatinAmericaCaribbean, "Honduras": LatinAmericaCaribbean, "Nicaragua": LatinAmericaCaribbean,
	"Costa Rica": LatinAmericaCaribbean, "Panama": LatinAmericaCaribbean, "Cuba": LatinAmericaCaribbean,
	"Jamaica": LatinAmericaCaribbean, "Haiti": LatinAmericaCaribbean, "the Dominican Republic": LatinAmericaCaribbean,
	"Puerto Rico": LatinAmericaCaribbean, "the Bahamas": LatinAmericaCaribbean, "Barbados": LatinAmericaCaribbean,
	"Trinidad and Tobago": LatinAmericaCaribbean, "Grenada": LatinAmericaCaribbean, "Saint Lucia": LatinAmericaCaribbean,
	"Saint Vincent and the Grenadines": LatinAmericaCaribbean, "Antigua and Barbuda": LatinAmericaCaribbean,
	"Dominica": LatinAmericaCaribbean, "Saint Kitts and Nevis": LatinAmericaCaribbean, "Aruba": LatinAmericaCaribbean,
	"Curaçao": LatinAmericaCaribbean, "Sint Maarten": LatinAmericaCaribbean, "Guadeloupe": LatinAmericaCaribbean,
	"Martinique": LatinAmericaCaribbean, "the Cayman Islands": LatinAmericaCaribbean,
	"the British Virgin Islands": LatinAmericaCaribbean, "the United States Virgin Islands": LatinAmericaCaribbean,
	"Anguilla": LatinAmericaCaribbean, "Montserrat": LatinAmericaCaribbean, "the Turks and Caicos Islands": LatinAmericaCaribbean,
	"Bonaire": LatinAmericaCaribbean, "Colombia": LatinAmericaCaribbean, "Venezuela": LatinAmericaCaribbean,
	"Guyana": LatinAmericaCaribbean, "Suriname": LatinAmericaCaribbean, "French Guiana": LatinAmericaCaribbean,
	"Ecuador": LatinAmericaCaribbean, "Peru": LatinAmericaCaribbean, "Brazil": LatinAmericaCaribbean,
	"Bolivia": LatinAmericaCaribbean, "Paraguay": LatinAmericaCaribbean, "Chile": LatinAmericaCaribbean,
	"Argentina": LatinAmericaCaribbean, "Uruguay": LatinAmericaCaribbean, "the Falkland Islands": LatinAmericaCaribbean,
	"the Netherlands Antilles": LatinAmericaCaribbean,

	// Europe (M49 places the whole of Russia in Europe)
	"Albania": Europe, "Andorra": Europe, "Austria": Europe, "Belarus": Europe, "Belgium": Europe,
	"Bosnia and Herzegovina": Europe, "Bulgaria": Europe, "Croatia": Europe, "the Czech Republic": Europe,
	"Czechia": Europe, "Denmark": Europe, "Estonia": Europe, "the Faroe Islands": Europe, "Finland": Europe,
	"France": Europe, "Germany": Europe, "Gibraltar": Europe, "Greece": Europe, "Guernsey": Europe,
	"Hungary": Europe, "Iceland": Europe, "Ireland": Europe, "the Isle of Man": Europe, "Italy": Europe,
	"Jersey": Europe, "Kosovo": Europe, "Latvia": Europe, "Liechtenstein": Europe, "Lithuania": Europe,
	"Luxembourg": Europe, "Malta": Europe, "Moldova": Europe, "Monaco": Europe, "Montenegro": Europe,
	"the Netherlands": Europe, "North Macedonia": Europe, "Norway": Europe, "Poland": Europe, "Portugal": Europe,
	"Romania": Europe, "Russia": Europe, "San Marino": Europe, "Serbia": Europe, "Slovakia": Europe,
	"Slovenia": Europe, "Spain": Europe, "Sweden": Europe, "Switzerland": Europe, "Ukraine": Europe,
	"the United Kingdom": Europe, "Vatican City": Europe,
	"Austria-Hungary": Europe, "the Soviet Union": Europe, "Yugoslavia": Europe, "Czechoslovakia": Europe,
	"East Germany": Europe, "West Germany": Europe, "Prussia": Europe, "the Russian Empire": Europe,
	"the German Empire": Europe,

	// Africa
	"Algeria": Africa, "Egypt": Africa, "Libya": Africa, "Morocco": Africa, "Sudan": Africa, "Tunisia": Africa,
	"Western Sahara": Africa, "Benin": Africa, "Burkina Faso": Africa, "Cape Verde": Africa, "Ivory Coast": Africa,
	"Côte d'Ivoire": Africa, "the Gambia": Africa, "Ghana": Africa, "Guinea": Africa, "Guinea-Bissau": Africa,
	"Liberia": Africa, "Mali": Africa, "Mauritania": Africa, "Niger": Africa, "Nigeria": Africa, "Senegal": Africa,
	"Sierra Leone": Africa, "Togo": Africa, "Burundi": Africa, "the Comoros": Africa, "Djibouti": Africa,
	"Eritrea": Africa, "Ethiopia": Africa, "Kenya": Africa, "Madagascar": Africa, "Malawi": Africa,
	"Mauritius": Africa, "Mozambique": Africa, "Rwanda": Africa, "Seychelles": Africa, "Somalia": Africa,
	"South Sudan": Africa, "Tanzania": Africa, "Uganda": Africa, "Zambia": Africa, "Zimbabwe": Africa,
	"Angola": Africa, "Cameroon": Africa, "the Central African Republic": Africa, "Chad": Africa,
	"the Republic of the Congo": Africa, "the Democratic Republic of the Congo": Africa,
	"Equatorial Guinea": Africa, "Gabon": Africa, "São Tomé and Príncipe": Africa, "Botswana": Africa,
	"Eswatini": Africa, "Lesotho": Africa, "Namibia": Africa, "South Africa": Africa, "Réunion": Africa,
	"Mayotte": Africa, "Saint Helena": Africa,
	"Belgian Congo": Africa, "Somaliland": Africa, "Rhodesia": Africa, "Gold Coast": Africa, "Zaire": Africa,
	"French West Africa": Africa, "French Equatorial Africa": Africa, "Tanganyika": Africa, "Zanzibar": Africa,

	// Middle East (M49 Western Asia)
	"Armenia": MiddleEast, "Azerbaijan": MiddleEast, "Bahrain": MiddleEast, "Cyprus": MiddleEast,
	"Georgia (country)": MiddleEast, "Iraq": MiddleEast, "Israel": MiddleEast, "Jordan": MiddleEast,
	"Kuwait": MiddleEast, "Lebanon": MiddleEast, "Oman": MiddleEast, "Qatar": MiddleEast,
	"Saudi Arabia": MiddleEast, "Palestine": MiddleEast, "the State of Palestine": MiddleEast,
	"Syria": MiddleEast, "Turkey": MiddleEast, "the United Arab Emirates": MiddleEast, "Yemen": MiddleEast,
	"the Ottoman Empire": MiddleEast, "Mandatory Palestine": MiddleEast,

	// South and Central Asia (M49 Southern Asia and Central Asia)
	"Kazakhstan": SouthAndCentralAsia, "Kyrgyzstan": SouthAndCentralAsia, "Tajikistan": SouthAndCentralAsia,
	"Turkmenistan": SouthAndCentralAsia, "Uzbekistan": SouthAndCentralAsia, "Afghanistan": SouthAndCentralAsia,
	"Bangladesh": SouthAndCentralAsia, "Bhutan": SouthAndCentralAsia, "India": SouthAndCentralAsia,
	"Iran": SouthAndCentralAsia, "the Maldives": SouthAndCentralAsia, "Nepal": SouthAndCentralAsia,
	"Pakistan": SouthAndCentralAsia, "Sri Lanka": SouthAndCentralAsia,
	"British India": SouthAndCentralAsia, "Ceylon": SouthAndCentralAsia, "Persia": SouthAndCentralAsia,
	"East Pakistan": SouthAndCentralAsia,

	// East and South-East Asia (M49 Eastern Asia and South-eastern Asia)
	"China": EastAndSouthEastAsia, "Hong Kong": EastAndSouthEastAsia, "Macau": EastAndSouthEastAsia,
	"Japan": EastAndSouthEastAsia, "Mongolia": EastAndSouthEastAsia, "North Korea": EastAndSouthEastAsia,
	"South Korea": EastAndSouthEastAsia, "Korea": EastAndSouthEastAsia, "Taiwan": EastAndSouthEastAsia,
	"Brunei": EastAndSouthEastAsia, "Cambodia": EastAndSouthEastAsia, "Indonesia": EastAndSouthEastAsia,
	"Laos": EastAndSouthEastAsia, "Malaysia": EastAndSouthEastAsia, "Myanmar": EastAndSouthEastAsia,
	"Burma": EastAndSouthEastAsia, "the Philippines": EastAndSouthEastAsia, "Singapore": EastAndSouthEastAsia,
	"Thailand": EastAndSouthEastAsia, "Timor-Leste": EastAndSouthEastAsia, "East Timor": EastAndSouthEastAsia,
	"Vietnam": EastAndSouthEastAsia, "Tibet": EastAndSouthEastAsia,
	"the Dutch East Indies": EastAndSouthEastAsia, "French Indochina": EastAndSouthEastAsia,
	"Manchukuo": EastAndSouthEastAsia, "Siam": EastAndSouthEastAsia, "Malaya": EastAndSouthEastAsia,

	// Oceania
	"Australia": Oceania, "New Zealand": Oceania, "Fiji": Oceania, "Papua New Guinea": Oceania,
	"the Solomon Islands": Oceania, "Vanuatu": Oceania, "New Caledonia": Oceania, "Guam": Oceania,
	"Kiribati": Oceania, "the Marshall Islands": Oceania, "the Federated States of Micronesia": Oceania,
	"Nauru": Oceania, "the Northern Mariana Islands": Oceania, "Palau": Oceania, "American Samoa": Oceania,
	"the Cook Islands": Oceania, "French Polynesia": Oceania, "Niue": Oceania, "Samoa": Oceania,
	"Tonga": Oceania, "Tuvalu": Oceania, "Wallis and Futuna": Oceania, "Norfolk Island": Oceania,
	"the Pitcairn Islands": Oceania, "Tokelau": Oceania, "German New Guinea": Oceania,

	// Antarctica
	"Antarctica": Antarctica,
}

// ZoneOfCountry returns the zone of a country, identified by its English short name.
func ZoneOfCountry(country string) (ID, bool) {
	zone, found := countryZones[country]
	return zone, found
}

// presentDayCountries places historical entities at the present-day country covering
// most of their territory, and other names at their usual one, to locate them on a map.
var presentDayCountries = map[string]string{
	"Austria-Hungary": "Austria", "the Soviet Union": "Russia", "the Russian Empire": "Russia",
	"Yugoslavia": "Serbia", "Czechoslovakia": "Czechia", "East Germany": "Germany", "West Germany": "Germany",
	"Prussia": "Germany", "the German Empire": "Germany", "the Ottoman Empire": "Turkey",
	"Mandatory Palestine": "Israel", "British India": "India", "Ceylon": "Sri Lanka", "Persia": "Iran",
	"East Pakistan": "Bangladesh", "Belgian Congo": "the Democratic Republic of the Congo",
	"Zaire": "the Democratic Republic of the Congo", "Rhodesia": "Zimbabwe", "Gold Coast": "Ghana",
	"Tanganyika": "Tanzania", "Burma": "Myanmar", "Siam": "Thailand", "the Dutch East Indies": "Indonesia",
	"German New Guinea": "Papua New Guinea", "the State of Palestine": "Palestine",
	"the Netherlands Antilles": "Curaçao",
}

// Countries returns the names of every country of the table, sorted.
func Countries() []string {
	return slices.Sorted(maps.Keys(countryZones))
}

// PresentDayCountry returns the usual present-day name of a country, which is the name
// itself for most countries.
func PresentDayCountry(country string) string {
	if presentDay, historical := presentDayCountries[country]; historical {
		return presentDay
	}
	return country
}

// PointOfCountry returns a point inside the country, [longitude, latitude], where a map
// would place its label. Territories without a border of their own have none.
func PointOfCountry(country string) ([2]float64, bool) {
	point, found := geozone.CountryPoints[country]
	return point, found
}
