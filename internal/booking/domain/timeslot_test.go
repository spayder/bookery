package domain_test

import (
	"testing"
	"time"

	"github.com/spayder/bookery/internal/booking/domain"
)

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("invalid fixture time %q: %v", s, err)
	}
	return parsed
}

func TestNewTimeSlot_RejectsEndBeforeOrEqualStart(t *testing.T) {
	start := mustTime(t, "2026-01-01T10:00:00Z")

	cases := map[string]time.Time{
		"end equals start": start,
		"end before start": start.Add(-time.Hour),
	}

	for name, end := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := domain.NewTimeSlot(start, end)
			if err != domain.ErrInvalidTimeSlot {
				t.Fatalf("expected ErrInvalidTimeSlot, got %v", err)
			}
		})
	}
}

func TestTimeSlot_Overlaps(t *testing.T) {
	base := mustTime(t, "2026-01-01T10:00:00Z")

	slot := func(startOffset, endOffset time.Duration) domain.TimeSlot {
		ts, err := domain.NewTimeSlot(base.Add(startOffset), base.Add(endOffset))
		if err != nil {
			t.Fatalf("unexpected error building fixture slot: %v", err)
		}
		return ts
	}

	a := slot(0, 2*time.Hour) // 10:00-12:00

	cases := []struct {
		name    string
		other   domain.TimeSlot
		overlap bool
	}{
		{"identical slot", slot(0, 2*time.Hour), true},
		{"partial overlap start", slot(-time.Hour, time.Hour), true},
		{"partial overlap end", slot(time.Hour, 3*time.Hour), true},
		{"contained inside", slot(30*time.Minute, 90*time.Minute), true},
		{"touches at end, no overlap", slot(2*time.Hour, 4*time.Hour), false},
		{"touches at start, no overlap", slot(-2*time.Hour, 0), false},
		{"fully separate", slot(3*time.Hour, 4*time.Hour), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := a.Overlaps(c.other); got != c.overlap {
				t.Errorf("a.Overlaps(%s) = %v, want %v", c.name, got, c.overlap)
			}
		})
	}
}
