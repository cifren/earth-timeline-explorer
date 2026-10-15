package timespan

import (
	"errors"
	"testing"
	"time"
)

func TestNewDecade(t *testing.T) {
	now := time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		year    int
		wantErr bool
	}{
		{name: "first decade", year: 1840},
		{name: "middle of the timeline", year: 1950},
		{name: "current decade", year: 2020},
		{name: "not a multiple of 10", year: 1955, wantErr: true},
		{name: "before the timeline", year: 1830, wantErr: true},
		{name: "future decade", year: 2030, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decade, err := NewDecade(test.year, now)
			if test.wantErr {
				if !errors.Is(err, ErrInvalidDecade) {
					t.Fatalf("NewDecade(%d) error = %v, want ErrInvalidDecade", test.year, err)
				}
				return
			}
			if err != nil || int(decade) != test.year {
				t.Fatalf("NewDecade(%d) = %d, %v", test.year, decade, err)
			}
		})
	}
}

func TestDecadeContains(t *testing.T) {
	decade := Decade(1950)
	for year, want := range map[int]bool{1949: false, 1950: true, 1959: true, 1960: false} {
		if got := decade.Contains(year); got != want {
			t.Errorf("Decade(1950).Contains(%d) = %v, want %v", year, got, want)
		}
	}
}

func TestDecadeSpan(t *testing.T) {
	got := Decade(1950).Span()
	want := TimeSpan{EDTF: "195X", Label: "1950s"}
	if got != want {
		t.Fatalf("Decade(1950).Span() = %+v, want %+v", got, want)
	}
}
