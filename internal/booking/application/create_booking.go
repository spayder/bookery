package application

import (
	"context"

	"github.com/spayder/bookery/internal/booking/domain"
)

type createBookingRepository interface {
	domain.OverlappingBookingFinder
	domain.BookingSaver
}

type CreateBookingCommand struct {
	ResourceID domain.ResourceID
	CustomerID domain.CustomerID
	Slot       domain.TimeSlot
}

type CreateBooking struct {
	repository     createBookingRepository
	policyProvider PolicyProvider
	clock          Clock
}

func NewCreateBooking(
	repository createBookingRepository,
	policyProvider PolicyProvider,
	clock Clock,
) *CreateBooking {
	return &CreateBooking{
		repository:     repository,
		policyProvider: policyProvider,
		clock:          clock,
	}
}

func (c *CreateBooking) Execute(ctx context.Context, cmd CreateBookingCommand) (domain.BookingID, error) {
	now := c.clock.Now()
	policy, err := c.policyProvider.PolicyFor(ctx, cmd.ResourceID)
	if err != nil {
		return domain.BookingID{}, err
	}
	booking, err := domain.NewBooking(cmd.ResourceID, cmd.CustomerID, cmd.Slot, policy, now)
	if err != nil {
		return domain.BookingID{}, err
	}
	if err := c.repository.Save(ctx, booking); err != nil {
		return domain.BookingID{}, err
	}
	return booking.ID(), nil
}
