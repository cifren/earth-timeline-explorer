package timespan

import "testing"

func TestFromText(t *testing.T) {
	decade := Decade(1950)
	tests := []struct {
		name string
		raw  string
		want TimeSpan
	}{
		{name: "ISO date", raw: "1953-04-03", want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "ISO date with time", raw: "1953-04-03 14:24:20", want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "numeric date is read day first", raw: "03/04/1953", want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "dotted numeric date", raw: "3.4.1953", want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "day month year", raw: "3 April 1953", want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "month day year", raw: "April 3, 1953", want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "month and year", raw: "October 1958", want: TimeSpan{"1958-10", "October 1958"}},
		{name: "ISO month", raw: "1958-10", want: TimeSpan{"1958-10", "October 1958"}},
		{name: "year in a sentence", raw: "Taken in 1959", want: TimeSpan{"1959", "1959"}},
		{name: "circa", raw: "circa 1955", want: TimeSpan{"1955~", "around 1955"}},
		{name: "ca. abbreviation", raw: "ca. 1955", want: TimeSpan{"1955~", "around 1955"}},
		{name: "between two years", raw: "between 1950 and 1955", want: TimeSpan{"1950/1955", "between 1950 and 1955"}},
		{name: "year range", raw: "1951–1954", want: TimeSpan{"1951/1954", "between 1951 and 1954"}},
		{name: "HTML is stripped", raw: `<time class="dtstart" datetime="1953-04-03">3 April 1953</time>`, want: TimeSpan{"1953-04-03", "3 April 1953"}},
		{name: "empty text gives the decade", raw: "", want: TimeSpan{"195X", "1950s"}},
		{name: "unreadable text gives the decade", raw: "unknown date", want: TimeSpan{"195X", "1950s"}},
		{name: "relative date gives the decade", raw: "before 1955", want: TimeSpan{"195X", "1950s"}},
		{name: "date outside the decade gives the decade", raw: "1962-05-01", want: TimeSpan{"195X", "1950s"}},
		{name: "range crossing the decade gives the decade", raw: "between 1948 and 1952", want: TimeSpan{"195X", "1950s"}},
		{name: "two different years give the decade", raw: "1952 or 1953", want: TimeSpan{"195X", "1950s"}},
		{name: "impossible day gives the decade", raw: "31/02/1953", want: TimeSpan{"195X", "1950s"}},
		{name: "decade text gives the decade", raw: "1950s", want: TimeSpan{"195X", "1950s"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FromText(test.raw, decade); got != test.want {
				t.Fatalf("FromText(%q) = %+v, want %+v", test.raw, got, test.want)
			}
		})
	}
}
