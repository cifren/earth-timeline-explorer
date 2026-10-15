package timespan

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimeSpan is a date or a period: EDTF is the source of truth, Label is for display.
type TimeSpan struct {
	EDTF  string
	Label string
}

var (
	tags       = regexp.MustCompile(`<[^>]*>`)
	spaces     = regexp.MustCompile(`\s+`)
	relative   = regexp.MustCompile(`(?i)\b(before|after|until|since|or)\b`)
	between    = regexp.MustCompile(`(?i)\bbetween\s+(\d{4})\s+and\s+(\d{4})\b`)
	yearRange  = regexp.MustCompile(`\b(\d{4})\s*[-–/]\s*(\d{4})\b`)
	circa      = regexp.MustCompile(`(?i)\b(?:circa|ca\.?|c\.|around|about|approx\.?)\s*(\d{4})\b`)
	isoDate    = regexp.MustCompile(`\b(\d{4})-(\d{2})-(\d{2})\b`)
	numeric    = regexp.MustCompile(`\b(\d{1,2})[./-](\d{1,2})[./-](\d{4})\b`)
	isoMonth   = regexp.MustCompile(`\b(\d{4})-(\d{2})\b`)
	dayMonth   = regexp.MustCompile(`(?i)\b(\d{1,2})\s+(` + monthNames + `)\s+(\d{4})\b`)
	monthDay   = regexp.MustCompile(`(?i)\b(` + monthNames + `)\s+(\d{1,2}),?\s+(\d{4})\b`)
	monthYear  = regexp.MustCompile(`(?i)\b(` + monthNames + `)\s+(\d{4})\b`)
	singleYear = regexp.MustCompile(`\b(1[89]\d{2}|20\d{2})\b`)
)

const monthNames = `january|february|march|april|may|june|july|august|september|october|november|december`

// FromText converts a free-text date of a source into a time span inside the decade.
//
// Ambiguous numeric dates are read day first (European order). A text that cannot be
// read, is relative ("before 1962"), or falls outside the decade gives the decade itself:
// the decade of the source category is the most reliable date.
func FromText(raw string, decade Decade) TimeSpan {
	text := strings.TrimSpace(spaces.ReplaceAllString(html.UnescapeString(tags.ReplaceAllString(raw, " ")), " "))
	if text == "" || relative.MatchString(text) {
		return decade.Span()
	}
	span, ok := parse(text, decade)
	if !ok {
		return decade.Span()
	}
	return span
}

func parse(text string, decade Decade) (TimeSpan, bool) {
	if m := between.FindStringSubmatch(text); m != nil {
		return yearsBetween(atoi(m[1]), atoi(m[2]), decade)
	}
	if m := isoDate.FindStringSubmatch(text); m != nil {
		return day(atoi(m[1]), atoi(m[2]), atoi(m[3]), decade)
	}
	if m := numeric.FindStringSubmatch(text); m != nil {
		return day(atoi(m[3]), atoi(m[2]), atoi(m[1]), decade)
	}
	if m := yearRange.FindStringSubmatch(text); m != nil {
		return yearsBetween(atoi(m[1]), atoi(m[2]), decade)
	}
	if m := circa.FindStringSubmatch(text); m != nil {
		return approximateYear(atoi(m[1]), decade)
	}
	if m := isoMonth.FindStringSubmatch(text); m != nil {
		return month(atoi(m[1]), atoi(m[2]), decade)
	}
	if m := dayMonth.FindStringSubmatch(text); m != nil {
		return day(atoi(m[3]), monthNumber(m[2]), atoi(m[1]), decade)
	}
	if m := monthDay.FindStringSubmatch(text); m != nil {
		return day(atoi(m[3]), monthNumber(m[1]), atoi(m[2]), decade)
	}
	if m := monthYear.FindStringSubmatch(text); m != nil {
		return month(atoi(m[2]), monthNumber(m[1]), decade)
	}
	return singleYearOf(text, decade)
}

func day(year, monthOfYear, dayOfMonth int, decade Decade) (TimeSpan, bool) {
	if !decade.Contains(year) || monthOfYear < 1 || monthOfYear > 12 {
		return TimeSpan{}, false
	}
	date := time.Date(year, time.Month(monthOfYear), dayOfMonth, 0, 0, 0, 0, time.UTC)
	if date.Day() != dayOfMonth {
		return TimeSpan{}, false
	}
	return TimeSpan{EDTF: date.Format(time.DateOnly), Label: date.Format("2 January 2006")}, true
}

func month(year, monthOfYear int, decade Decade) (TimeSpan, bool) {
	if !decade.Contains(year) || monthOfYear < 1 || monthOfYear > 12 {
		return TimeSpan{}, false
	}
	return TimeSpan{
		EDTF:  fmt.Sprintf("%04d-%02d", year, monthOfYear),
		Label: fmt.Sprintf("%s %d", time.Month(monthOfYear), year),
	}, true
}

func approximateYear(year int, decade Decade) (TimeSpan, bool) {
	if !decade.Contains(year) {
		return TimeSpan{}, false
	}
	return TimeSpan{EDTF: fmt.Sprintf("%04d~", year), Label: fmt.Sprintf("around %d", year)}, true
}

func yearsBetween(from, to int, decade Decade) (TimeSpan, bool) {
	if from > to || !decade.Contains(from) || !decade.Contains(to) {
		return TimeSpan{}, false
	}
	if from == to {
		return TimeSpan{EDTF: fmt.Sprintf("%04d", from), Label: strconv.Itoa(from)}, true
	}
	return TimeSpan{
		EDTF:  fmt.Sprintf("%04d/%04d", from, to),
		Label: fmt.Sprintf("between %d and %d", from, to),
	}, true
}

// singleYearOf accepts a text only when it names exactly one year.
func singleYearOf(text string, decade Decade) (TimeSpan, bool) {
	years := singleYear.FindAllString(text, -1)
	if len(years) == 0 {
		return TimeSpan{}, false
	}
	for _, year := range years[1:] {
		if year != years[0] {
			return TimeSpan{}, false
		}
	}
	return yearsBetween(atoi(years[0]), atoi(years[0]), decade)
}

// monthNumber returns 1 to 12 for an English month name, 0 otherwise.
func monthNumber(name string) int {
	for i, candidate := range strings.Split(monthNames, "|") {
		if strings.EqualFold(candidate, name) {
			return i + 1
		}
	}
	return 0
}

func atoi(digits string) int {
	value, _ := strconv.Atoi(digits)
	return value
}
