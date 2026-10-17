package domain

import (
	"errors"
	"time"
)

var ErrSlotInPast = errors.New("cannot book a slot that has already started")
var ErrCannotConfirm = errors.New("only pending bookings can be confirmed")
var ErrCannotCancel = errors.New("only active bookings can be cancelled")
var ErrMissingCancellationPolicy = errors.New("missing cancellation policy")

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

type Booking struct {
	id         BookingID
	resourceID ResourceID
	customerID CustomerID
	slot       TimeSlot
	policy     CancellationPolicy
	status     Status
}

func NewBooking(resourceID ResourceID, customerID CustomerID, slot TimeSlot, policy CancellationPolicy, now time.Time) (*Booking, error) {
	if slot.Start().Before(now) {
		return nil, ErrSlotInPast
	}
	if policy == nil {
		return nil, ErrMissingCancellationPolicy
	}

	return &Booking{
		id:         NewBookingID(),
		resourceID: resourceID,
		customerID: customerID,
		slot:       slot,
		policy:     policy,
		status:     StatusPending,
	}, nil
}

func (b *Booking) Status() Status {
	return b.status
}

func (b *Booking) ID() BookingID {
	return b.id
}

func (b *Booking) ResourceID() ResourceID {
	return b.resourceID
}

func (b *Booking) CustomerID() CustomerID {
	return b.customerID
}

func (b *Booking) Slot() TimeSlot {
	return b.slot
}

func (b *Booking) Confirm() error {
	if b.status != StatusPending {
		return ErrCannotConfirm
	}

	b.status = StatusConfirmed
	return nil
}

func (b *Booking) Cancel(now time.Time) error {
	if !b.IsActive() {
		return ErrCannotCancel
	}

	if err := b.policy.CheckCancellation(b.slot, now); err != nil {
		return err
	}

	b.status = StatusCancelled
	return nil
}

func (b *Booking) IsActive() bool {
	return b.status == StatusPending || b.status == StatusConfirmed
}
