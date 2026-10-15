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
	start := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	slot, err := domain.NewTimeSlot(start, start.Add(2*time.Hour))
	if err != nil {
		t.Fatalf("invalid fixture slot: %v", err)
	}
	policy, err := domain.NewNoticePeriodPolicy(0)
	if err != nil {
		t.Fatalf("invalid fixture policy: %v", err)
	}
	b, err := domain.NewBooking(domain.NewResourceID(), domain.NewCustomerID(), slot, policy, start.Add(-time.Hour))
	if err != nil {
		t.Fatalf("invalid fixture booking: %v", err)
	}
	return b
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
