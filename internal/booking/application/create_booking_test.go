package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spayder/bookery/internal/booking/application"
	"github.com/spayder/bookery/internal/booking/domain"
	"github.com/spayder/bookery/internal/booking/infrastructure/memory"
)

// fixedClock is a test Clock that always returns the same moment.
type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

// fixedPolicyProvider is a test PolicyProvider that returns the same policy
// for every resource.
type fixedPolicyProvider struct {
	policy domain.CancellationPolicy
}

func (p fixedPolicyProvider) PolicyFor(_ context.Context, _ domain.ResourceID) (domain.CancellationPolicy, error) {
	return p.policy, nil
}

// Compile-time checks: the build fails if the fakes stop satisfying the ports.
var (
	_ application.Clock          = fixedClock{}
	_ application.PolicyProvider = fixedPolicyProvider{}
)

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

func TestCreateBooking_SavesPendingBooking(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	createBooking := application.NewCreateBooking(
		repo,
		fixedPolicyProvider{policy: domain.NonCancellablePolicy{}},
		fixedClock{now: testNow},
	)
	cmd := application.CreateBookingCommand{
		ResourceID: domain.NewResourceID(),
		CustomerID: domain.NewCustomerID(),
		Slot:       mustSlot(t, 10, 12),
	}

	id, err := createBooking.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute: unexpected error: %v", err)
	}

	saved, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID: booking was not saved: %v", err)
	}
	if saved.ResourceID() != cmd.ResourceID {
		t.Errorf("ResourceID = %v, want %v", saved.ResourceID(), cmd.ResourceID)
	}
	if saved.CustomerID() != cmd.CustomerID {
		t.Errorf("CustomerID = %v, want %v", saved.CustomerID(), cmd.CustomerID)
	}
	if saved.Slot() != cmd.Slot {
		t.Errorf("Slot = %v, want %v", saved.Slot(), cmd.Slot)
	}
	if saved.Status() != domain.StatusPending {
		t.Errorf("Status = %v, want %v", saved.Status(), domain.StatusPending)
	}
}

func TestCreateBooking_RefusesDoubleBooking(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	createBooking := application.NewCreateBooking(
		repo,
		fixedPolicyProvider{policy: domain.NonCancellablePolicy{}},
		fixedClock{now: testNow},
	)
	resourceID := domain.NewResourceID()
	cmd := application.CreateBookingCommand{
		ResourceID: resourceID,
		CustomerID: domain.NewCustomerID(),
		Slot:       mustSlot(t, 10, 12),
	}

	_, err := createBooking.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute: unexpected error: %v", err)
	}

	secondCmd := cmd
	secondCmd.CustomerID = domain.NewCustomerID()

	_, err = createBooking.Execute(ctx, secondCmd)
	if !errors.Is(err, application.ErrResourceAlreadyBooked) {
		t.Fatalf("Execute: err = %v, want %v", err, application.ErrResourceAlreadyBooked)
	}

	bookings, err := repo.FindOverlapping(ctx, resourceID, mustSlot(t, 10, 12))
	if err != nil {
		t.Fatalf("FindOverlapping: unexpected error: %v", err)
	}
	if len(bookings) != 1 {
		t.Fatalf("FindOverlapping: expected 1 booking, got %d", len(bookings))
	}
}

func TestCreateBooking_CancelledBookingDoesNotBlockNewBooking(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	policy, err := domain.NewNoticePeriodPolicy(0)
	if err != nil {
		t.Fatalf("NewNoticePeriodPolicy: unexpected error: %v", err)
	}

	createBooking := application.NewCreateBooking(
		repo,
		fixedPolicyProvider{policy: policy},
		fixedClock{now: testNow},
	)
	resourceID := domain.NewResourceID()
	cmd := application.CreateBookingCommand{
		ResourceID: resourceID,
		CustomerID: domain.NewCustomerID(),
		Slot:       mustSlot(t, 10, 12),
	}

	bookingID, err := createBooking.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute: unexpected error: %v", err)
	}

	saved, err := repo.FindByID(ctx, bookingID)
	if err != nil {
		t.Fatalf("FindByID: unexpected error: %v", err)
	}

	if err := saved.Cancel(testNow); err != nil {
		t.Fatalf("Cancel: unexpected error: %v", err)
	}

	if err := repo.Save(ctx, saved); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	_, err = createBooking.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute: unexpected error: %v", err)
	}
}
