package domain_test

import (
	"testing"
	"time"

	"github.com/spayder/bookery/internal/booking/domain"
)

func mustSlot(t *testing.T, start, end string) domain.TimeSlot {
	t.Helper()
	slot, err := domain.NewTimeSlot(mustTime(t, start), mustTime(t, end))
	if err != nil {
		t.Fatalf("invalid fixture slot %s-%s: %v", start, end, err)
	}
	return slot
}

func TestNewBooking_StartsPendingWithGivenDetails(t *testing.T) {
	resourceID := domain.NewResourceID()
	customerID := domain.NewCustomerID()
	slot := mustSlot(t, "2026-01-01T10:00:00Z", "2026-01-01T12:00:00Z")
	now := mustTime(t, "2026-01-01T09:00:00Z")

	booking, err := domain.NewBooking(resourceID, customerID, slot, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if booking.ID().String() == "" {
		t.Error("expected a generated, non-empty booking ID")
	}
	if booking.Status() != domain.StatusPending {
		t.Errorf("Status() = %q, want %q", booking.Status(), domain.StatusPending)
	}
	if booking.ResourceID() != resourceID {
		t.Errorf("ResourceID() = %q, want %q", booking.ResourceID(), resourceID)
	}
	if booking.CustomerID() != customerID {
		t.Errorf("CustomerID() = %q, want %q", booking.CustomerID(), customerID)
	}
	if !booking.Slot().Start().Equal(slot.Start()) || !booking.Slot().End().Equal(slot.End()) {
		t.Errorf("Slot() = %v-%v, want %v-%v",
			booking.Slot().Start(), booking.Slot().End(), slot.Start(), slot.End())
	}
}

func TestNewBooking_GeneratesUniqueIDs(t *testing.T) {
	slot := mustSlot(t, "2026-01-01T10:00:00Z", "2026-01-01T12:00:00Z")
	now := mustTime(t, "2026-01-01T09:00:00Z")

	a, err := domain.NewBooking(domain.NewResourceID(), domain.NewCustomerID(), slot, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := domain.NewBooking(domain.NewResourceID(), domain.NewCustomerID(), slot, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.ID() == b.ID() {
		t.Errorf("expected different IDs, both were %q", a.ID())
	}
}

func TestNewBooking_StartTimeRelativeToNow(t *testing.T) {
	slot := mustSlot(t, "2026-01-01T10:00:00Z", "2026-01-01T12:00:00Z")

	cases := []struct {
		name    string
		now     time.Time
		wantErr error
	}{
		{"slot starts in the future", slot.Start().Add(-time.Hour), nil},
		{"slot starts exactly now (walk-in)", slot.Start(), nil},
		{"slot already started", slot.Start().Add(time.Minute), domain.ErrSlotInPast},
		{"slot already finished", slot.End().Add(time.Hour), domain.ErrSlotInPast},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			booking, err := domain.NewBooking(domain.NewResourceID(), domain.NewCustomerID(), slot, c.now)

			if err != c.wantErr {
				t.Fatalf("error = %v, want %v", err, c.wantErr)
			}
			if c.wantErr != nil && booking != nil {
				t.Errorf("expected nil booking on error, got %+v", booking)
			}
		})
	}
}
