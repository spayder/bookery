package domain

import (
	"context"
	"errors"
)

var ErrBookingNotFound = errors.New("booking not found")

type BookingFinder interface {
	FindByID(ctx context.Context, id BookingID) (*Booking, error)
}

type BookingSaver interface {
	Save(ctx context.Context, booking *Booking) error
}

type BookingRepository interface {
	BookingFinder
	BookingSaver
}
