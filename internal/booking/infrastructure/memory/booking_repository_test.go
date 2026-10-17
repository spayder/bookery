package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spayder/bookery/internal/booking/domain"
	"github.com/spayder/bookery/internal/booking/infrastructure/memory"
)

// Compile-time check: the build fails if memory.BookingRepository
// stops satisfying the port defined in the domain.
var _ domain.BookingRepository = (*memory.BookingRepository)(nil)

// newTestBooking builds a valid pending booking. The memory_test package
// can't see the helpers in domain_test, so it has its own.
func newTestBooking(t *testing.T) *domain.Booking {
	t.Helper()
	return newBookingFor(t, domain.NewResourceID(), mustSlot(t, 10, 12))
}

// newBookingFor builds a valid pending booking for the given resource and slot.
func newBookingFor(t *testing.T, resourceID domain.ResourceID, slot domain.TimeSlot) *domain.Booking {
	t.Helper()
	policy, err := domain.NewNoticePeriodPolicy(0)
	if err != nil {
		t.Fatalf("invalid fixture policy: %v", err)
	}
	b, err := domain.NewBooking(resourceID, domain.NewCustomerID(), slot, policy, testNow)
	if err != nil {
		t.Fatalf("invalid fixture booking: %v", err)
	}
	return b
}

// testNow is before every fixture slot, so NewBooking accepts them all.
var testNow = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// mustSlot builds a slot on 2026-01-01 from startHour:00 to endHour:00 UTC.
func mustSlot(t *testing.T, startHour, endHour int) domain.TimeSlot {
	t.Helper()
	slot, err := domain.NewTimeSlot(
		time.Date(2026, 1, 1, startHour, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 1, endHour, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("invalid fixture slot %d-%d: %v", startHour, endHour, err)
	}
	return slot
}

func TestBookingRepository_SaveThenFindByID(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	want := newTestBooking(t)

	if err := repo.Save(ctx, want); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	got, err := repo.FindByID(ctx, want.ID())
	if err != nil {
		t.Fatalf("FindByID: unexpected error: %v", err)
	}

	if got.ID() != want.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), want.ID())
	}
	if got.ResourceID() != want.ResourceID() {
		t.Errorf("ResourceID = %v, want %v", got.ResourceID(), want.ResourceID())
	}
	if got.CustomerID() != want.CustomerID() {
		t.Errorf("CustomerID = %v, want %v", got.CustomerID(), want.CustomerID())
	}
	if got.Slot() != want.Slot() {
		t.Errorf("Slot = %v, want %v", got.Slot(), want.Slot())
	}
	if got.Status() != want.Status() {
		t.Errorf("Status = %v, want %v", got.Status(), want.Status())
	}
}

// A loaded booking is the caller's own copy: changing it must not touch the
// stored one until Save is called, just as with a real database.
func TestBookingRepository_ChangesToLoadedBookingAreNotStoredWithoutSave(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	original := newTestBooking(t)
	if err := repo.Save(ctx, original); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	loaded, err := repo.FindByID(ctx, original.ID())
	if err != nil {
		t.Fatalf("FindByID: unexpected error: %v", err)
	}
	if err := loaded.Confirm(); err != nil {
		t.Fatalf("Confirm: unexpected error: %v", err)
	}

	reloaded, err := repo.FindByID(ctx, original.ID())
	if err != nil {
		t.Fatalf("FindByID: unexpected error: %v", err)
	}
	if reloaded.Status() != domain.StatusPending {
		t.Errorf("stored Status = %v, want %v (the change was never saved)", reloaded.Status(), domain.StatusPending)
	}
}

func TestBookingRepository_ChangesToOriginalAfterSaveAreNotStored(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	original := newTestBooking(t)
	if err := repo.Save(ctx, original); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	if err := original.Confirm(); err != nil {
		t.Fatalf("Confirm: unexpected error: %v", err)
	}

	reloaded, err := repo.FindByID(ctx, original.ID())
	if err != nil {
		t.Fatalf("FindByID: unexpected error: %v", err)
	}
	if reloaded.Status() != domain.StatusPending {
		t.Errorf("stored Status = %v, want %v (the change was never saved)", reloaded.Status(), domain.StatusPending)
	}
}

func TestBookingRepository_NonExistingIDReturnsErrBookingNotFound(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()

	got, err := repo.FindByID(ctx, domain.NewBookingID())

	if got != nil {
		t.Errorf("FindByID: expected to receive no booking, got %v", got)
	}

	if !errors.Is(err, domain.ErrBookingNotFound) {
		t.Errorf("FindByID: unexpected error: %v, want %v", err, domain.ErrBookingNotFound)
	}
}

func TestBookingRepository_SaveOverwritesExistingBooking(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	original := newTestBooking(t)
	if err := repo.Save(ctx, original); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	if err := original.Confirm(); err != nil {
		t.Fatalf("Confirm: unexpected error: %v", err)
	}

	if err := repo.Save(ctx, original); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	reloaded, err := repo.FindByID(ctx, original.ID())
	if err != nil {
		t.Fatalf("FindByID: unexpected error: %v", err)
	}
	if reloaded.Status() != domain.StatusConfirmed {
		t.Errorf("stored Status = %v, want %v", reloaded.Status(), domain.StatusConfirmed)
	}
}

func TestBookingRepository_FindOverlapping(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name         string
		bookedSlot   domain.TimeSlot
		querySlot    domain.TimeSlot
		sameResource bool
		wantCount    int
		cancelled    bool
	}{
		{"overlapping slots same resource", mustSlot(t, 10, 12), mustSlot(t, 11, 13), true, 1, false},
		{"same slot different resource", mustSlot(t, 10, 12), mustSlot(t, 10, 12), false, 0, false},
		{"no overlap", mustSlot(t, 10, 12), mustSlot(t, 13, 14), true, 0, false},
		{"touching: query starts when booking ends", mustSlot(t, 10, 12), mustSlot(t, 12, 14), true, 0, false},
		{"touching: query ends when booking starts", mustSlot(t, 10, 12), mustSlot(t, 8, 10), true, 0, false},
		{"booking inside query", mustSlot(t, 11, 12), mustSlot(t, 10, 13), true, 1, false},
		{"query inside booking", mustSlot(t, 10, 13), mustSlot(t, 11, 12), true, 1, false},
		{"booking cancelled", mustSlot(t, 10, 12), mustSlot(t, 11, 13), true, 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := memory.NewBookingRepository()
			resourceID := domain.NewResourceID()
			var rID domain.ResourceID
			if tc.sameResource {
				rID = resourceID
			} else {
				rID = domain.NewResourceID()
			}
			existing := newBookingFor(t, rID, tc.bookedSlot)
			if err := repo.Save(ctx, existing); err != nil {
				t.Fatalf("Save: unexpected error: %v", err)
			}

			if tc.cancelled {
				if err := existing.Cancel(testNow); err != nil {
					t.Fatalf("Cancel: unexpected error: %v", err)
				}
				if err := repo.Save(ctx, existing); err != nil {
					t.Fatalf("Save: unexpected error: %v", err)
				}
			}

			got, err := repo.FindOverlapping(ctx, resourceID, tc.querySlot)
			if err != nil {
				t.Fatalf("FindOverlapping: unexpected error: %v", err)
			}

			if len(got) != tc.wantCount {
				t.Fatalf("got %d bookings, want %d", len(got), tc.wantCount)
			}
			if len(got) > 0 && got[0].ID() != existing.ID() {
				t.Errorf("ID = %v, want %v", got[0].ID(), existing.ID())
			}
		})
	}
}
