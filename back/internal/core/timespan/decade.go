// Package timespan holds the time vocabulary shared by every context: decades, and
// time spans expressed in EDTF (ISO 8601-2) with a plain English label.
package timespan

import (
	"errors"
	"fmt"
	"time"
)

// FirstDecade is the oldest decade offered on the timeline.
const FirstDecade = 1840

// ErrInvalidDecade is returned for a year that is not a decade of the timeline.
var ErrInvalidDecade = errors.New("invalid decade")

// Decade is the first year of a ten-year period, e.g. 1950 for 1950–1959.
type Decade int

// NewDecade validates a decade against the timeline bounds at the given time.
func NewDecade(year int, now time.Time) (Decade, error) {
	current := now.Year() / 10 * 10
	if year%10 != 0 {
		return 0, fmt.Errorf("%w: %d is not a multiple of 10", ErrInvalidDecade, year)
	}
	if year < FirstDecade || year > current {
		return 0, fmt.Errorf("%w: %d is outside %d–%d", ErrInvalidDecade, year, FirstDecade, current)
	}
	return Decade(year), nil
}

// Contains reports whether a year belongs to the decade.
func (d Decade) Contains(year int) bool {
	return year >= int(d) && year <= int(d)+9
}

// Span returns the whole decade as a time span, e.g. "195X" / "1950s".
func (d Decade) Span() TimeSpan {
	return TimeSpan{
		EDTF:  fmt.Sprintf("%dX", int(d)/10),
		Label: fmt.Sprintf("%ds", int(d)),
	}
}
