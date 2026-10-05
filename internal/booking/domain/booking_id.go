package domain

import (
	"errors"

	"uuid"
)

var ErrEmptyBookingID = errors.New("booking Id can not be empty")

type BookingID struct {
	value string
}

func NewBookingID() BookingID {
	return BookingID{value: uuid.New().String()}
}

func BookingIDFromString(s string) (BookingID, error) {
	if s == "" {
		return BookingID{}, ErrEmptyBookingID
	}
	return BookingID{value: s}, nil
}

func (id BookingID) String() string {
	return id.value
}
