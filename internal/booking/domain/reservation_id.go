package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrEmptyReservationID = errors.New("reservation Id can not be empty")

type ReservationID struct {
	value string
}

func NewReservationID() ReservationID {
	return ReservationID{value: uuid.NewString()}
}

func ReservationIDFromString(s string) (ReservationID, error) {
	if s == "" {
		return ReservationID{}, ErrEmptyReservationID
	}
	return ReservationID{value: s}, nil
}

func (r ReservationID) String() string {
	return r.value
}
