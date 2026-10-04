package domain

import (
	"errors"
	"time"
)

var ErrSlotInPast = errors.New("cannot book a slot that has already started")
var ErrCannotConfirm = errors.New("only pending bookings can be confirmed")
var ErrCannotCancel = errors.New("only active bookings can be cancelled")
var ErrAlreadyStarted = errors.New("can not cancel a booking that has already started")

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
	status     Status
}

func NewBooking(resourceID ResourceID, customerID CustomerID, slot TimeSlot, now time.Time) (*Booking, error) {
	if slot.Start().Before(now) {
		return nil, ErrSlotInPast
	}

	return &Booking{
		id:         NewBookingID(),
		resourceID: resourceID,
		customerID: customerID,
		slot:       slot,
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
	if b.status != StatusPending && b.status != StatusConfirmed {
		return ErrCannotCancel
	}

	if !now.Before(b.slot.Start()) {
		return ErrAlreadyStarted
	}

	b.status = StatusCancelled
	return nil
}
